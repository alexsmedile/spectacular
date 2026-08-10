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
	Path   string
	Data   []byte
	Mode   os.FileMode
	Delete bool
}

type transactionFile struct {
	Target      string `json:"target"`
	Temporary   string `json:"temporary"`
	Backup      string `json:"backup"`
	HadOriginal bool   `json:"had_original"`
	Mode        uint32 `json:"mode"`
	Delete      bool   `json:"delete,omitempty"`
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
	return applyTransactionWithInstallHook(root, key, changes, failAfter, nil)
}

func applyTransactionWithInstallHook(root, key string, changes []FileChange, failAfter int, beforeInstall func(int)) error {
	if len(changes) == 0 {
		return nil
	}
	if strings.TrimSpace(key) == "" {
		return domain.NewRefusal(domain.RefusalMissingRequiredField, "idempotency_key", "transaction key is required", nil)
	}
	if err := RecoverTransactions(root); err != nil {
		return err
	}
	effects, err := openRootedWorkspace(root)
	if err != nil {
		return err
	}
	defer effects.Close()
	digest := sha256.Sum256([]byte(key))
	txID := hex.EncodeToString(digest[:16])
	txRelative := filepath.Join(".spectacular", "transactions")
	if err := effects.mkdirAll(txRelative, 0o755, ".spectacular"); err != nil {
		return transactionRefusal("create transaction directory", err)
	}
	journalRelative := filepath.Join(".spectacular", "transactions", txID+".json")
	journal := transactionJournal{Schema: "spectacular.transaction.v1", Key: key}
	seen := map[string]bool{}
	journalWritten := false
	defer func() {
		if !journalWritten {
			cleanupPrepared(effects, journal.Files)
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
			Delete:    change.Delete,
		}
		if _, pathErr := transactionArtifact(root, item.Temporary, txID, i, ".new"); pathErr != nil {
			return pathErr
		}
		if _, pathErr := transactionArtifact(root, item.Backup, txID, i, ".old"); pathErr != nil {
			return pathErr
		}
		if item.Mode == 0 {
			item.Mode = uint32(0o644)
		}
		journal.Files = append(journal.Files, item)
		tracked := &journal.Files[len(journal.Files)-1]
		if info, statErr := effects.stat(change.Path, ""); statErr == nil {
			tracked.HadOriginal = true
			tracked.Mode = uint32(info.Mode().Perm())
			data, readErr := effects.readFile(change.Path, "")
			if readErr != nil {
				return transactionRefusal("read transaction original", readErr)
			}
			if _, pathErr := transactionArtifact(root, item.Backup, txID, i, ".old"); pathErr != nil {
				return pathErr
			}
			if writeErr := effects.writeSynced(item.Backup, data, info.Mode().Perm(), txRelative); writeErr != nil {
				return transactionRefusal("write transaction backup", writeErr)
			}
		} else if !os.IsNotExist(statErr) {
			return transactionRefusal("inspect transaction target", statErr)
		}
		if change.Delete {
			continue
		}
		if _, err := safeTarget(root, change.Path); err != nil {
			return err
		}
		if err := effects.mkdirAll(filepath.Dir(change.Path), 0o755, ""); err != nil {
			return transactionRefusal("create transaction target directory", err)
		}
		if _, err := safeTarget(root, change.Path); err != nil {
			return err
		}
		if _, pathErr := transactionArtifact(root, item.Temporary, txID, i, ".new"); pathErr != nil {
			return pathErr
		}
		if err := effects.writeSynced(item.Temporary, change.Data, os.FileMode(tracked.Mode), txRelative); err != nil {
			return transactionRefusal("write transaction candidate", err)
		}
	}
	journalData, err := json.MarshalIndent(journal, "", "  ")
	if err != nil {
		return transactionRefusal("encode transaction journal", err)
	}
	journalData = append(journalData, '\n')
	if err := effects.writeSynced(journalRelative, journalData, 0o600, txRelative); err != nil {
		return transactionRefusal("write transaction journal", err)
	}
	journalWritten = true
	installed := 0
	for i, item := range journal.Files {
		if failAfter >= 0 && installed == failAfter {
			if err := rollback(effects, journal); err != nil {
				return err
			}
			if err := removeEffect(effects, journalRelative, txRelative); err != nil {
				return err
			}
			return transactionRefusal("injected transaction interruption", fmt.Errorf("after %d installs", installed))
		}
		if _, pathErr := transactionArtifact(root, item.Temporary, txID, i, ".new"); pathErr != nil {
			return pathErr
		}
		if _, pathErr := safeTarget(root, item.Target); pathErr != nil {
			return pathErr
		}
		if beforeInstall != nil {
			beforeInstall(i)
		}
		if item.Delete {
			if err := removeEffect(effects, item.Target, ""); err != nil {
				if rollbackErr := rollback(effects, journal); rollbackErr != nil {
					return rollbackErr
				}
				return transactionRefusal("delete transaction target", err)
			}
			installed++
			continue
		}
		if err := effects.rename(item.Temporary, item.Target, txRelative, ""); err != nil {
			if rollbackErr := rollback(effects, journal); rollbackErr != nil {
				return rollbackErr
			}
			if removeErr := removeEffect(effects, journalRelative, txRelative); removeErr != nil {
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
		if err := removeEffect(effects, item.Backup, txRelative); err != nil {
			return err
		}
		if _, err := transactionArtifact(root, item.Temporary, txID, i, ".new"); err != nil {
			return err
		}
		if err := removeEffect(effects, item.Temporary, txRelative); err != nil {
			return err
		}
	}
	if err := removeEffect(effects, journalRelative, txRelative); err != nil {
		return err
	}
	return nil
}

// RecoverTransactions rolls back every incomplete journal. Recovery is
// deterministic and happens before any subsequent governed mutation.
func RecoverTransactions(root string) error {
	effects, err := openRootedWorkspace(root)
	if err != nil {
		return err
	}
	defer effects.Close()
	txRelative := filepath.Join(".spectacular", "transactions")
	entries, err := effects.readDir(txRelative, ".spectacular")
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
		journalRelative := filepath.Join(txRelative, entry.Name())
		data, readErr := effects.readFile(journalRelative, txRelative)
		if readErr != nil {
			return transactionRefusal("read recovery journal", readErr)
		}
		var journal transactionJournal
		if jsonErr := json.Unmarshal(data, &journal); jsonErr != nil || journal.Schema != "spectacular.transaction.v1" {
			return domain.NewRefusal(domain.RefusalTransactionRecovery, journalRelative, "transaction journal is invalid", jsonErr)
		}
		digest := sha256.Sum256([]byte(journal.Key))
		txID := hex.EncodeToString(digest[:16])
		if entry.Name() != txID+".json" {
			return domain.NewRefusal(domain.RefusalPathEscape, entry.Name(), "transaction journal identity does not match its key", nil)
		}
		if rollbackErr := rollback(effects, journal); rollbackErr != nil {
			return rollbackErr
		}
		if removeErr := removeEffect(effects, journalRelative, txRelative); removeErr != nil {
			return removeErr
		}
	}
	return nil
}

func rollback(effects *rootedWorkspace, journal transactionJournal) error {
	digest := sha256.Sum256([]byte(journal.Key))
	txID := hex.EncodeToString(digest[:16])
	txRelative := filepath.Join(".spectacular", "transactions")
	for i := len(journal.Files) - 1; i >= 0; i-- {
		item := journal.Files[i]
		if _, err := safeTarget(effects.path, item.Target); err != nil {
			return err
		}
		if _, err := transactionArtifact(effects.path, item.Backup, txID, i, ".old"); err != nil {
			return err
		}
		if _, err := transactionArtifact(effects.path, item.Temporary, txID, i, ".new"); err != nil {
			return err
		}
		if item.HadOriginal {
			data, err := effects.readFile(item.Backup, txRelative)
			if err != nil {
				return transactionRefusal("read rollback backup", err)
			}
			if _, err = safeTarget(effects.path, item.Target); err != nil {
				return err
			}
			if err := effects.writeSynced(item.Target, data, os.FileMode(item.Mode), ""); err != nil {
				return transactionRefusal("restore rollback target", err)
			}
		} else if err := removeEffect(effects, item.Target, ""); err != nil {
			return err
		}
		if err := removeEffect(effects, item.Temporary, txRelative); err != nil {
			return err
		}
		if err := removeEffect(effects, item.Backup, txRelative); err != nil {
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

type rootedWorkspace struct {
	path string
	root *os.Root
}

func openRootedWorkspace(path string) (*rootedWorkspace, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, domain.NewRefusal(domain.RefusalPathEscape, path, "resolve canonical workspace root", err)
	}
	root, err := os.OpenRoot(canonical)
	if err != nil {
		return nil, domain.NewRefusal(domain.RefusalPathEscape, path, "open canonical workspace root", err)
	}
	return &rootedWorkspace{path: canonical, root: root}, nil
}

func (r *rootedWorkspace) Close() {
	_ = r.root.Close()
}

func (r *rootedWorkspace) validate(relative, scope string) (string, error) {
	name := filepath.FromSlash(relative)
	if name == "." && scope == "" {
		return name, nil
	}
	if _, err := effectPath(r.path, name, scope); err != nil {
		return "", err
	}
	return name, nil
}

func (r *rootedWorkspace) operationError(relative, scope string, operationErr error) error {
	if operationErr == nil {
		return nil
	}
	if _, validationErr := r.validate(relative, scope); validationErr != nil {
		return validationErr
	}
	return operationErr
}

func (r *rootedWorkspace) stat(relative, scope string) (os.FileInfo, error) {
	name, err := r.validate(relative, scope)
	if err != nil {
		return nil, err
	}
	info, err := r.root.Stat(name)
	return info, r.operationError(relative, scope, err)
}

func (r *rootedWorkspace) readFile(relative, scope string) ([]byte, error) {
	name, err := r.validate(relative, scope)
	if err != nil {
		return nil, err
	}
	data, err := r.root.ReadFile(name)
	return data, r.operationError(relative, scope, err)
}

func (r *rootedWorkspace) readDir(relative, scope string) ([]os.DirEntry, error) {
	name, err := r.validate(relative, scope)
	if err != nil {
		return nil, err
	}
	directory, err := r.root.Open(name)
	if err != nil {
		return nil, r.operationError(relative, scope, err)
	}
	defer directory.Close()
	entries, err := directory.ReadDir(-1)
	return entries, err
}

func (r *rootedWorkspace) mkdirAll(relative string, mode os.FileMode, scope string) error {
	name, err := r.validate(relative, scope)
	if err != nil {
		return err
	}
	err = r.root.MkdirAll(name, mode.Perm())
	return r.operationError(relative, scope, err)
}

func (r *rootedWorkspace) writeSynced(relative string, data []byte, mode os.FileMode, scope string) error {
	name, err := r.validate(relative, scope)
	if err != nil {
		return err
	}
	file, err := r.root.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
	if err != nil {
		return r.operationError(relative, scope, err)
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

func (r *rootedWorkspace) rename(oldRelative, newRelative, oldScope, newScope string) error {
	oldName, err := r.validate(oldRelative, oldScope)
	if err != nil {
		return err
	}
	newName, err := r.validate(newRelative, newScope)
	if err != nil {
		return err
	}
	if err := r.root.Rename(oldName, newName); err != nil {
		if _, validationErr := r.validate(oldRelative, oldScope); validationErr != nil {
			return validationErr
		}
		if _, validationErr := r.validate(newRelative, newScope); validationErr != nil {
			return validationErr
		}
		return err
	}
	return nil
}

func (r *rootedWorkspace) remove(relative, scope string) error {
	name, err := r.validate(relative, scope)
	if err != nil {
		return err
	}
	err = r.root.Remove(name)
	if os.IsNotExist(err) {
		return nil
	}
	return r.operationError(relative, scope, err)
}

func cleanupPrepared(effects *rootedWorkspace, files []transactionFile) {
	for _, item := range files {
		for _, path := range []string{item.Temporary, item.Backup} {
			_ = effects.remove(path, filepath.Join(".spectacular", "transactions"))
		}
	}
}

func removeEffect(effects *rootedWorkspace, relative, scope string) error {
	if err := effects.remove(relative, scope); err != nil {
		return transactionRefusal("remove contained transaction artifact", err)
	}
	return nil
}

func pathWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func transactionRefusal(operation string, err error) error {
	return domain.NewRefusal(domain.RefusalPersistence, "transaction", operation, err)
}
