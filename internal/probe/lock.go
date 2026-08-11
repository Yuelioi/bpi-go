package probe

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func BuildContractLock(root, sourceCommit string) (ContractLock, error) {
	if sourceCommit == "" {
		return ContractLock{}, fmt.Errorf("probe: source commit cannot be empty")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return ContractLock{}, fmt.Errorf("probe: resolve contract root: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return ContractLock{}, fmt.Errorf("probe: stat contract root: %w", err)
	}
	if !info.IsDir() {
		return ContractLock{}, fmt.Errorf("probe: contract root is not a directory")
	}
	lock := ContractLock{SchemaVersion: SchemaVersion, SourceCommit: sourceCommit}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink is not allowed in contract snapshot: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		digest := sha256.Sum256(data)
		lock.Files = append(lock.Files, LockedFile{
			Path:   filepath.ToSlash(relative),
			SHA256: hex.EncodeToString(digest[:]),
			Size:   int64(len(data)),
		})
		return nil
	})
	if err != nil {
		return ContractLock{}, fmt.Errorf("probe: hash contract snapshot: %w", err)
	}
	sort.Slice(lock.Files, func(i, j int) bool { return lock.Files[i].Path < lock.Files[j].Path })
	return lock, nil
}

func VerifyContractLock(root string, want ContractLock) error {
	if want.SchemaVersion != SchemaVersion {
		return fmt.Errorf("probe: contract lock schema is %d, want %d", want.SchemaVersion, SchemaVersion)
	}
	got, err := BuildContractLock(root, want.SourceCommit)
	if err != nil {
		return err
	}
	if len(got.Files) != len(want.Files) {
		return fmt.Errorf("probe: contract snapshot has %d files, lock has %d", len(got.Files), len(want.Files))
	}
	for index := range want.Files {
		if got.Files[index] != want.Files[index] {
			return fmt.Errorf("probe: contract snapshot drift at %q", want.Files[index].Path)
		}
	}
	return nil
}

func ReadJSON(path string, destination any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("probe: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("probe: decode %s: %w", path, err)
	}
	return nil
}

func WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("probe: encode %s: %w", path, err)
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("probe: create output directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("probe: write %s: %w", path, err)
	}
	return nil
}
