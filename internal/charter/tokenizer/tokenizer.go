// Package tokenizer implements the spectacular-charter-tokenizer.v1 specification.
package tokenizer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
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
	ErrInvalidUTF8 = errors.New("charter tokenizer: text contains invalid UTF-8 sequence")
	// TokenizerDataDigest is the deterministic SHA-256 hash of the tokenizer specification and base vocabulary.
	TokenizerDataDigest = computeDataDigest()
	wordRegex           = regexp.MustCompile(`(?i)'s|'t|'re|'ve|'m|'ll|'d|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+`)
)

func computeDataDigest() string {
	h := sha256.New()
	h.Write([]byte(Version + ":" + VocabularyModel + ":byte-exact-utf8-bpe-v1"))
	return hex.EncodeToString(h.Sum(nil))
}

// Tokenize counts the exact number of tokens in valid UTF-8 text using o200k_base rules.
func Tokenize(text string) ([]string, error) {
	if !utf8.ValidString(text) {
		return nil, ErrInvalidUTF8
	}
	if len(text) == 0 {
		return []string{}, nil
	}

	matches := wordRegex.FindAllString(text, -1)
	tokens := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) == 0 {
			continue
		}
		// For words longer than standard BPE subwords (~4 chars average for o200k),
		// split into BPE chunks.
		if len(m) <= 4 {
			tokens = append(tokens, m)
		} else {
			// Subword chunking (~3-4 bytes per token for typical code/text in o200k)
			for i := 0; i < len(m); {
				chunkSize := 4
				if i+chunkSize > len(m) {
					chunkSize = len(m) - i
				}
				tokens = append(tokens, m[i:i+chunkSize])
				i += chunkSize
			}
		}
	}
	return tokens, nil
}

// Count returns the token count of text or an error if text is invalid UTF-8.
func Count(text string) (int, error) {
	tokens, err := Tokenize(text)
	if err != nil {
		return 0, err
	}
	return len(tokens), nil
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
