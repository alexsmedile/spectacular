package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
)

type FileChange struct {
	Path string
	Data []byte
	Mode os.FileMode
}

type transactionFile struct {
	Target      string `json:"target"`
	Temporary   string `json:"temporary"`
	Backup      string `json:"backup"`
	HadOriginal bool   `json:"had_original"`
	Mode        uint32 `json:"mode"`
}

type transactionJournal struct {
	Schema string            `json:"schema"`
	Key    string            `json:"key"`
	Files  []transactionFile `json:"files"`
}

// ApplyTransaction installs a set of canonical files as one recoverable
// logical transaction. A process interruption may leave installed files, but
// RecoverTransactions restores every original before a later mutation.
func ApplyTransaction(root, key string, changes []FileChange) error {
	return applyTransaction(root, key, changes, -1)
}

func applyTransaction(root, key string, changes []FileChange, failAfter int) error {
	if len(changes) == 0 {
		return nil
	}
	if strings.TrimSpace(key) == "" {
		return domain.NewRefusal(domain.RefusalMissingRequiredField, "idempotency_key", "transaction key is required", nil)
	}
	if err := RecoverTransactions(root); err != nil {
		return err
	}
	digest := sha256.Sum256([]byte(key))
	txID := hex.EncodeToString(digest[:16])
	txRoot, err := effectPath(root, filepath.Join(".spectacular", "transactions"), ".spectacular")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(txRoot, 0o755); err != nil {
		return transactionRefusal("create transaction directory", err)
	}
	if _, err := effectPath(root, filepath.Join(".spectacular", "transactions"), ".spectacular"); err != nil {
		return err
	}
	journalRelative := filepath.Join(".spectacular", "transactions", txID+".json")
	journalPath, err := effectPath(root, journalRelative, filepath.Join(".spectacular", "transactions"))
	if err != nil {
		return err
	}
	journal := transactionJournal{Schema: "spectacular.transaction.v1", Key: key}
	seen := map[string]bool{}
	journalWritten := false
	defer func() {
		if !journalWritten {
			cleanupPrepared(root, journal.Files)
		}
	}()
	for i, change := range changes {
		target, err := safeTarget(root, change.Path)
		if err != nil {
			return err
		}
		if seen[target] {
			return domain.NewRefusal(domain.RefusalCollision, change.Path, "transaction targets the same path twice", nil)
		}
		seen[target] = true
		item := transactionFile{
			Target:    filepath.ToSlash(change.Path),
			Temporary: filepath.ToSlash(filepath.Join(".spectacular", "transactions", fmt.Sprintf("%s-%d.new", txID, i))),
			Backup:    filepath.ToSlash(filepath.Join(".spectacular", "transactions", fmt.Sprintf("%s-%d.old", txID, i))),
			Mode:      uint32(change.Mode.Perm()),
		}
		temporary, pathErr := transactionArtifact(root, item.Temporary, txID, i, ".new")
		if pathErr != nil {
			return pathErr
		}
		backup, pathErr := transactionArtifact(root, item.Backup, txID, i, ".old")
		if pathErr != nil {
			return pathErr
		}
		if item.Mode == 0 {
			item.Mode = uint32(0o644)
		}
		journal.Files = append(journal.Files, item)
		tracked := &journal.Files[len(journal.Files)-1]
		if info, statErr := os.Stat(target); statErr == nil {
			tracked.HadOriginal = true
			tracked.Mode = uint32(info.Mode().Perm())
			data, readErr := os.ReadFile(target)
			if readErr != nil {
				return transactionRefusal("read transaction original", readErr)
			}
			if writeErr := writeSynced(backup, data, info.Mode().Perm()); writeErr != nil {
				return transactionRefusal("write transaction backup", writeErr)
			}
		} else if !os.IsNotExist(statErr) {
			return transactionRefusal("inspect transaction target", statErr)
		}
		if _, err := safeTarget(root, change.Path); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return transactionRefusal("create transaction target directory", err)
		}
		if _, err := safeTarget(root, change.Path); err != nil {
			return err
		}
		if err := writeSynced(temporary, change.Data, os.FileMode(tracked.Mode)); err != nil {
			return transactionRefusal("write transaction candidate", err)
		}
	}
	journalData, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return transactionRefusal("encode transaction journal", err)
	}
	journalData = append(journalData, '\n')
	if err := writeSynced(journalPath, journalData, 0o600); err != nil {
		return transactionRefusal("write transaction journal", err)
	}
	journalWritten = true
	installed := 0
	for i, item := range journal.Files {
		if failAfter >= 0 && installed == failAfter {
			if err := rollback(root, journal); err != nil {
				return err
			}
			if err := removeEffect(root, journalRelative, filepath.Join(".spectacular", "transactions")); err != nil {
				return err
			}
			return transactionRefusal("injected transaction interruption", fmt.Errorf("after %d installs", installed))
		}
		temporary, pathErr := transactionArtifact(root, item.Temporary, txID, i, ".new")
		if pathErr != nil {
			return pathErr
		}
		target, pathErr := safeTarget(root, item.Target)
		if pathErr != nil {
			return pathErr
		}
		if err := os.Rename(temporary, target); err != nil {
			if rollbackErr := rollback(root, journal); rollbackErr != nil {
				return rollbackErr
			}
			if removeErr := removeEffect(root, journalRelative, filepath.Join(".spectacular", "transactions")); removeErr != nil {
				return removeErr
			}
			return transactionRefusal("install transaction candidate", err)
		}
		installed++
	}
	for i, item := range journal.Files {
		if _, err := transactionArtifact(root, item.Backup, txID, i, ".old"); err != nil {
			return err
		}
		if err := removeEffect(root, item.Backup, filepath.Join(".spectacular", "transactions")); err != nil {
			return err
		}
		if _, err := transactionArtifact(root, item.Temporary, txID, i, ".new"); err != nil {
			return err
		}
		if err := removeEffect(root, item.Temporary, filepath.Join(".spectacular", "transactions")); err != nil {
			return err
		}
	}
	if err := removeEffect(root, journalRelative, filepath.Join(".spectacular", "transactions")); err != nil {
		return err
	}
	return nil
}

// RecoverTransactions rolls back every incomplete journal. Recovery is
// deterministic and happens before any subsequent governed mutation.
func RecoverTransactions(root string) error {
	txRoot, pathErr := effectPath(root, filepath.Join(".spectacular", "transactions"), ".spectacular")
	if pathErr != nil {
		return pathErr
	}
	entries, err := os.ReadDir(txRoot)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return transactionRefusal("read transaction directory", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return domain.NewRefusal(domain.RefusalPathEscape, entry.Name(), "transaction journal must not be a symlink", nil)
		}
		path, effectErr := effectPath(root, filepath.Join(".spectacular", "transactions", entry.Name()), filepath.Join(".spectacular", "transactions"))
		if effectErr != nil {
			return effectErr
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return transactionRefusal("read recovery journal", readErr)
		}
		var journal transactionJournal
		if jsonErr := json.Unmarshal(data, &journal); jsonErr != nil || journal.Schema != "spectacular.transaction.v1" {
			return domain.NewRefusal(domain.RefusalTransactionRecovery, path, "transaction journal is invalid", jsonErr)
		}
		digest := sha256.Sum256([]byte(journal.Key))
		txID := hex.EncodeToString(digest[:16])
		if entry.Name() != txID+".json" {
			return domain.NewRefusal(domain.RefusalPathEscape, entry.Name(), "transaction journal identity does not match its key", nil)
		}
		if rollbackErr := rollback(root, journal); rollbackErr != nil {
			return rollbackErr
		}
		if removeErr := removeEffect(root, filepath.Join(".spectacular", "transactions", entry.Name()), filepath.Join(".spectacular", "transactions")); removeErr != nil {
			return removeErr
		}
	}
	return nil
}

func rollback(root string, journal transactionJournal) error {
	digest := sha256.Sum256([]byte(journal.Key))
	txID := hex.EncodeToString(digest[:16])
	for i := len(journal.Files) - 1; i >= 0; i-- {
		item := journal.Files[i]
		target, err := safeTarget(root, item.Target)
		if err != nil {
			return err
		}
		backup, err := transactionArtifact(root, item.Backup, txID, i, ".old")
		if err != nil {
			return err
		}
		if _, err := transactionArtifact(root, item.Temporary, txID, i, ".new"); err != nil {
			return err
		}
		if item.HadOriginal {
			data, err := os.ReadFile(backup)
			if err != nil {
				return transactionRefusal("read rollback backup", err)
			}
			target, err = safeTarget(root, item.Target)
			if err != nil {
				return err
			}
			if err := writeSynced(target, data, os.FileMode(item.Mode)); err != nil {
				return transactionRefusal("restore rollback target", err)
			}
		} else if err := removeEffect(root, item.Target, ""); err != nil {
			return err
		}
		if err := removeEffect(root, item.Temporary, filepath.Join(".spectacular", "transactions")); err != nil {
			return err
		}
		if err := removeEffect(root, item.Backup, filepath.Join(".spectacular", "transactions")); err != nil {
			return err
		}
	}
	return nil
}

func safeTarget(root, relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "transaction target must remain workspace-relative", nil)
	}
	if relative == "" || filepath.Clean(relative) != relative || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, relative, "transaction target must be canonical and workspace-relative", nil)
	}
	return effectPath(root, relative, "")
}

func effectPath(root, relative, scope string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) || filepath.Clean(relative) != relative || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, relative, "effect path must be canonical and workspace-relative", nil)
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return "", domain.NewRefusal(domain.RefusalPathEscape, root, "resolve canonical workspace root", err)
	}
	if scope != "" {
		cleanScope := filepath.Clean(scope)
		if relative != cleanScope && !strings.HasPrefix(relative, cleanScope+string(filepath.Separator)) {
			return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "effect path escapes required transaction scope", nil)
		}
	}
	target := filepath.Join(rootReal, relative)
	ancestor := filepath.Dir(target)
	for {
		if _, statErr := os.Lstat(ancestor); statErr == nil {
			break
		} else if !os.IsNotExist(statErr) {
			return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "inspect effect parent", statErr)
		}
		next := filepath.Dir(ancestor)
		if next == ancestor {
			return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "effect parent has no contained ancestor", nil)
		}
		ancestor = next
	}
	ancestorReal, err := filepath.EvalSymlinks(ancestor)
	if err != nil || !pathWithin(rootReal, ancestorReal) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "effect parent escapes canonical workspace", err)
	}
	if info, statErr := os.Lstat(target); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
		resolved, resolveErr := filepath.EvalSymlinks(target)
		if resolveErr != nil || !pathWithin(rootReal, resolved) {
			return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "effect target symlink escapes canonical workspace", resolveErr)
		}
	} else if statErr != nil && !os.IsNotExist(statErr) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "inspect effect target", statErr)
	}
	return target, nil
}

func transactionArtifact(root, relative, txID string, index int, suffix string) (string, error) {
	want := filepath.ToSlash(filepath.Join(".spectacular", "transactions", fmt.Sprintf("%s-%d%s", txID, index, suffix)))
	if filepath.ToSlash(relative) != want {
		return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "transaction artifact identity is invalid", nil)
	}
	return effectPath(root, filepath.FromSlash(want), filepath.Join(".spectacular", "transactions"))
}

func cleanupPrepared(root string, files []transactionFile) {
	for _, item := range files {
		for _, path := range []string{item.Temporary, item.Backup} {
			if absolute, err := effectPath(root, filepath.FromSlash(path), filepath.Join(".spectacular", "transactions")); err == nil {
				_ = os.Remove(absolute)
			}
		}
	}
}

func removeEffect(root, relative, scope string) error {
	absolute, err := effectPath(root, filepath.FromSlash(relative), scope)
	if err != nil {
		return err
	}
	if err := os.Remove(absolute); err != nil && !os.IsNotExist(err) {
		return transactionRefusal("remove contained transaction artifact", err)
	}
	return nil
}

func pathWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func writeSynced(path string, data []byte, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func transactionRefusal(operation string, err error) error {
	return domain.NewRefusal(domain.RefusalPersistence, "transaction", operation, err)
}
