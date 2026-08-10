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
	txRoot := filepath.Join(root, ".spectacular", "transactions")
	if err := os.MkdirAll(txRoot, 0o755); err != nil {
		return transactionRefusal("create transaction directory", err)
	}
	journalPath := filepath.Join(txRoot, txID+".json")
	journal := transactionJournal{Schema: "spectacular.transaction.v1", Key: key}
	seen := map[string]bool{}
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
			Target:    target,
			Temporary: filepath.Join(txRoot, fmt.Sprintf("%s-%d.new", txID, i)),
			Backup:    filepath.Join(txRoot, fmt.Sprintf("%s-%d.old", txID, i)),
			Mode:      uint32(change.Mode.Perm()),
		}
		if item.Mode == 0 {
			item.Mode = uint32(0o644)
		}
		if info, statErr := os.Stat(target); statErr == nil {
			item.HadOriginal = true
			item.Mode = uint32(info.Mode().Perm())
			data, readErr := os.ReadFile(target)
			if readErr != nil {
				return transactionRefusal("read transaction original", readErr)
			}
			if writeErr := writeSynced(item.Backup, data, info.Mode().Perm()); writeErr != nil {
				return transactionRefusal("write transaction backup", writeErr)
			}
		} else if !os.IsNotExist(statErr) {
			return transactionRefusal("inspect transaction target", statErr)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return transactionRefusal("create transaction target directory", err)
		}
		if err := writeSynced(item.Temporary, change.Data, os.FileMode(item.Mode)); err != nil {
			return transactionRefusal("write transaction candidate", err)
		}
		journal.Files = append(journal.Files, item)
	}
	journalData, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return transactionRefusal("encode transaction journal", err)
	}
	journalData = append(journalData, '\n')
	if err := writeSynced(journalPath, journalData, 0o600); err != nil {
		return transactionRefusal("write transaction journal", err)
	}
	installed := 0
	for _, item := range journal.Files {
		if failAfter >= 0 && installed == failAfter {
			if err := rollback(journal); err != nil {
				return err
			}
			_ = os.Remove(journalPath)
			return transactionRefusal("injected transaction interruption", fmt.Errorf("after %d installs", installed))
		}
		if err := os.Rename(item.Temporary, item.Target); err != nil {
			if rollbackErr := rollback(journal); rollbackErr != nil {
				return rollbackErr
			}
			_ = os.Remove(journalPath)
			return transactionRefusal("install transaction candidate", err)
		}
		installed++
	}
	for _, item := range journal.Files {
		_ = os.Remove(item.Backup)
		_ = os.Remove(item.Temporary)
	}
	if err := os.Remove(journalPath); err != nil && !os.IsNotExist(err) {
		return transactionRefusal("remove completed transaction journal", err)
	}
	return nil
}

// RecoverTransactions rolls back every incomplete journal. Recovery is
// deterministic and happens before any subsequent governed mutation.
func RecoverTransactions(root string) error {
	txRoot := filepath.Join(root, ".spectacular", "transactions")
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
		path := filepath.Join(txRoot, entry.Name())
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return transactionRefusal("read recovery journal", readErr)
		}
		var journal transactionJournal
		if jsonErr := json.Unmarshal(data, &journal); jsonErr != nil || journal.Schema != "spectacular.transaction.v1" {
			return domain.NewRefusal(domain.RefusalTransactionRecovery, path, "transaction journal is invalid", jsonErr)
		}
		if rollbackErr := rollback(journal); rollbackErr != nil {
			return rollbackErr
		}
		if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
			return transactionRefusal("remove recovered journal", removeErr)
		}
	}
	return nil
}

func rollback(journal transactionJournal) error {
	for i := len(journal.Files) - 1; i >= 0; i-- {
		item := journal.Files[i]
		if item.HadOriginal {
			data, err := os.ReadFile(item.Backup)
			if err != nil {
				return transactionRefusal("read rollback backup", err)
			}
			if err := writeSynced(item.Target, data, os.FileMode(item.Mode)); err != nil {
				return transactionRefusal("restore rollback target", err)
			}
		} else if err := os.Remove(item.Target); err != nil && !os.IsNotExist(err) {
			return transactionRefusal("remove rollback target", err)
		}
		_ = os.Remove(item.Temporary)
		_ = os.Remove(item.Backup)
	}
	return nil
}

func safeTarget(root, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) || filepath.Clean(relative) != relative || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalInvalidWorkspacePath, relative, "transaction target must be canonical and workspace-relative", nil)
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target := filepath.Join(rootAbs, relative)
	rel, err := filepath.Rel(rootAbs, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", domain.NewRefusal(domain.RefusalPathEscape, relative, "transaction target escapes workspace", err)
	}
	return target, nil
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
