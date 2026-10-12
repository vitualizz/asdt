package installer

import (
	"errors"
	"os"
	"path/filepath"
)

// atomicWriteFile writes data to a temp file beside path, sets its mode to
// perm, fsyncs it, then renames it over path, so a concurrent reader never
// sees a half-written file and a failed write leaves the target untouched.
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	f, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".asdt-tmp-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	fail := func(err error) error {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	// Chmod, not the create mode: the umask must not alter the mode we keep.
	if err := f.Chmod(perm); err != nil {
		return fail(err)
	}
	if _, err := f.Write(data); err != nil {
		return fail(err)
	}
	if err := f.Sync(); err != nil {
		return fail(err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// atomicReplaceFile atomically writes data to the file path names, keeping
// what the user set up around it: when path is a symlink, the file it points
// to is rewritten and the link itself stays in place; an existing file keeps
// its mode; a new file gets 0o644 and any missing parent directories.
func atomicReplaceFile(path string, data []byte) error {
	target := path
	perm := os.FileMode(0o644)
	resolved, err := filepath.EvalSymlinks(path)
	switch {
	case err == nil:
		target = resolved
		info, err := os.Stat(target)
		if err != nil {
			return err
		}
		perm = info.Mode().Perm()
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
	default:
		return err
	}
	return atomicWriteFile(target, data, perm)
}
