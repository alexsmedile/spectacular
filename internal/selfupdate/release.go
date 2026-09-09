package selfupdate

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// releaseAPI is the GitHub endpoint naming the newest published release. It is
// a variable so a test can point it at a local server.
var releaseAPI = "https://api.github.com/repos/alexsmedile/spectacular/releases/latest"

// downloadBase is the prefix release assets hang from, completed with the tag.
var downloadBase = "https://github.com/alexsmedile/spectacular/releases/download"

// maxArchiveBytes caps a download. A release archive is a few megabytes; a
// response far larger is a redirect to something unexpected rather than a
// build, and streaming it into memory unbounded is how an update turns into an
// out-of-memory kill.
const maxArchiveBytes = 128 << 20

// Release is the newest published release and the asset for this platform.
type Release struct {
	Version  string
	Tag      string
	Archive  string
	Platform string
}

// Platform names the release asset slice for the running machine, in the same
// os-arch form assemble-release uses.
func Platform() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

// LatestRelease asks GitHub for the newest published release. A release with no
// assets is reported as such rather than silently treated as absent: three tags
// published no archives when the release workflow broke, and an update that
// called that "up to date" would hide the outage.
func LatestRelease(ctx context.Context, client *http.Client) (Release, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, releaseAPI, nil)
	if err != nil {
		return Release{}, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	response, err := client.Do(request)
	if err != nil {
		return Release{}, fmt.Errorf("reach the release API: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf("release API returned %s", response.Status)
	}
	var payload struct {
		TagName string `json:"tag_name"`
		Draft   bool   `json:"draft"`
		Assets  []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4<<20)).Decode(&payload); err != nil {
		return Release{}, fmt.Errorf("parse the release API response: %w", err)
	}
	if payload.Draft || payload.TagName == "" {
		return Release{}, errors.New("no published release found")
	}
	release := Release{
		Tag:      payload.TagName,
		Version:  strings.TrimPrefix(payload.TagName, "v"),
		Platform: Platform(),
	}
	want := "spectacular-" + payload.TagName + "-" + release.Platform + ".tar.gz"
	for _, asset := range payload.Assets {
		if asset.Name == want {
			release.Archive = asset.Name
			return release, nil
		}
	}
	if len(payload.Assets) == 0 {
		return release, fmt.Errorf("release %s published no assets", payload.TagName)
	}
	return release, fmt.Errorf("release %s has no archive for %s", payload.TagName, release.Platform)
}

// DownloadBinary fetches the release archive, verifies it against the release
// SHA256SUMS, and returns the binary held inside it.
//
// The checksum is verified over the whole archive before a single entry is
// read, so nothing extracted has escaped verification.
func DownloadBinary(ctx context.Context, client *http.Client, release Release) ([]byte, error) {
	if release.Archive == "" {
		return nil, errors.New("release has no archive for this platform")
	}
	base := downloadBase + "/" + release.Tag + "/"
	sums, err := fetch(ctx, client, base+"SHA256SUMS", 4<<20)
	if err != nil {
		return nil, fmt.Errorf("download SHA256SUMS: %w", err)
	}
	expected, err := checksumFor(string(sums), release.Archive)
	if err != nil {
		return nil, err
	}
	archive, err := fetch(ctx, client, base+release.Archive, maxArchiveBytes)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", release.Archive, err)
	}
	digest := sha256.Sum256(archive)
	if actual := hex.EncodeToString(digest[:]); actual != expected {
		return nil, fmt.Errorf("checksum mismatch for %s: archive is %s, SHA256SUMS says %s", release.Archive, actual, expected)
	}
	return extractBinary(archive)
}

func fetch(ctx context.Context, client *http.Client, url string, limit int64) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned %s", response.Status)
	}
	return io.ReadAll(io.LimitReader(response.Body, limit))
}

// checksumFor finds one archive's line in a SHA256SUMS file. An absent entry is
// an error rather than a skipped check: an unverifiable download must refuse.
func checksumFor(sums, archive string) (string, error) {
	for _, line := range strings.Split(sums, "\n") {
		checksum, name, found := strings.Cut(strings.TrimSpace(line), "  ")
		if !found || name != archive {
			continue
		}
		if len(checksum) != 64 {
			return "", fmt.Errorf("SHA256SUMS entry for %s is malformed", archive)
		}
		return checksum, nil
	}
	return "", fmt.Errorf("SHA256SUMS has no entry for %s", archive)
}

// extractBinary pulls spectacular/bin/spectacular out of a verified archive.
// Every other entry is ignored: this replaces one binary and has no business
// writing anything else to disk.
func extractBinary(archive []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(strings.NewReader(string(archive)))
	if err != nil {
		return nil, fmt.Errorf("read release archive: %w", err)
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read release archive: %w", err)
		}
		if header.Name != "spectacular/bin/spectacular" {
			continue
		}
		if header.Typeflag != tar.TypeReg {
			return nil, errors.New("release binary is not a regular file")
		}
		data, err := io.ReadAll(io.LimitReader(reader, maxArchiveBytes))
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return nil, errors.New("release binary is empty")
		}
		return data, nil
	}
	return nil, errors.New("release archive holds no spectacular binary")
}

// HTTPClient is the client used for release downloads, with a timeout so a
// hung connection fails instead of blocking an interactive command forever.
func HTTPClient() *http.Client {
	return &http.Client{Timeout: 2 * time.Minute}
}
