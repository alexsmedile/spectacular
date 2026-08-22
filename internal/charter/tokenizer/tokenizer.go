package tokenizer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"unicode/utf8"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

const (
	Version           = "spectacular-charter-tokenizer.v1"
	VocabularyModel   = "o200k_base"
	MaxTargetTokens   = 1200
	MaxWarningTokens  = 1400
	MaxCeilingTokens  = 1440
	HardCeilingTokens = 1440
)

type Disposition string

const (
	DispositionPass             Disposition = "pass"
	DispositionWarn             Disposition = "warn"
	DispositionSplitRecommended Disposition = "split_recommended"
	DispositionRefusal          Disposition = "refusal"
)

var (
	ErrInvalidUTF8      = errors.New("charter tokenizer: text contains invalid UTF-8 sequence")
	TokenizerDataDigest = computeDataDigest()

	bpeOnce    sync.Once
	bpeEncoder *tiktoken.Tiktoken
	bpeErr     error
)

func getEncoder() (*tiktoken.Tiktoken, error) {
	bpeOnce.Do(func() {
		bpeEncoder, bpeErr = tiktoken.GetEncoding(VocabularyModel)
	})
	return bpeEncoder, bpeErr
}

func computeDataDigest() string {
	h := sha256.New()
	h.Write([]byte(Version + ":" + VocabularyModel + ":official-o200k-base-v1"))
	return hex.EncodeToString(h.Sum(nil))
}

// Tokenize counts the exact number of tokens in valid UTF-8 text using official o200k_base BPE rules.
func Tokenize(text string) ([]string, error) {
	if !utf8.ValidString(text) {
		return nil, ErrInvalidUTF8
	}
	if len(text) == 0 {
		return []string{}, nil
	}
	enc, err := getEncoder()
	if err != nil {
		return nil, err
	}
	tokenIDs := enc.Encode(text, nil, nil)
	tokens := make([]string, len(tokenIDs))
	for i, id := range tokenIDs {
		tokens[i] = enc.Decode([]int{id})
	}
	return tokens, nil
}

// Count returns the exact token count of text using official o200k_base BPE.
func Count(text string) (int, error) {
	if !utf8.ValidString(text) {
		return 0, ErrInvalidUTF8
	}
	if len(text) == 0 {
		return 0, nil
	}
	enc, err := getEncoder()
	if err != nil {
		return 0, err
	}
	tokenIDs := enc.Encode(text, nil, nil)
	return len(tokenIDs), nil
}

// EvaluateDisposition checks a token count against the frozen threshold gates.
func EvaluateDisposition(tokenCount int) (Disposition, string) {
	switch {
	case tokenCount <= MaxTargetTokens:
		return DispositionPass, "Token count is within the 1,200 token budget envelope."
	case tokenCount <= MaxWarningTokens:
		return DispositionWarn, "Token count exceeds 1,200 (within 1,400 warning envelope). Safe compaction applied."
	case tokenCount <= MaxCeilingTokens:
		return DispositionSplitRecommended, "Token count exceeds 1,400 (within 1,440 ceiling). Strong recommendation to split Objective into serial Runs."
	default:
		return DispositionRefusal, "Token count exceeds the hard 1,440 ceiling. Compilation refused."
	}
}

// FormatReceiptSummary produces a deterministic string receipt for the charter.
func FormatReceiptSummary(tokenCount int, disposition Disposition) string {
	var b strings.Builder
	b.WriteString("Tokenizer: " + Version + " (" + VocabularyModel + ")\n")
	b.WriteString("Digest: sha256:" + TokenizerDataDigest + "\n")
	b.WriteString("Disposition: " + string(disposition) + "\n")
	return b.String()
}
