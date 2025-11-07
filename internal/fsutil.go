package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Templated Path
type TPath_t string

// Perform template substitutions from vars
// and return the resolved path string
func (p TPath_t) Resolve(vars map[string]string) (string, error) {
	t := template.New("Vars").Option("missingkey=error")
	if _, err := t.Parse(string(p)); err != nil {
		return "", err
	}
	var buf strings.Builder
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Copy a file or dir (src)
// to destination dir (destDir),
// creating destDir & parent dirs leading up to it
func Copy(src, dstDir string) error {
	fInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if fInfo.IsDir() {
		err = os.CopyFS(dstDir,
			os.DirFS(src))
	} else {
		err = copyFile(src, dstDir)
	}

	return err
}

// Copy file src to dstDir.
// dstDir is created if it doesn't exist
func copyFile(src, dstDir string) error {
	err := os.MkdirAll(dstDir, 0700)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	dstFname := filepath.Join(dstDir, filepath.Base(src))
	out, err := os.OpenFile(dstFname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// Returns absolute path relsolving '~'
func AbsPath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		path = home
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func RealPath(link string) (string, error) {
	// Get the file info using Lstat, which doesn't follow symlinks
	info, err := os.Lstat(link)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %v", err)
	}

	// If it's not a symlink, return the original path
	if info.Mode()&os.ModeSymlink == 0 {
		return link, nil
	}

	// If it is a symlink, read the target
	realPath, err := os.Readlink(link)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink: %v", err)
	}

	// Return the real target of the symlink
	return realPath, nil
}

// ln -sf target link
func CreateSymlinkForce(target, link string) error {
	// Check if the symlink already exists
	if _, err := os.Lstat(link); err == nil {
		// If it exists, remove it
		err = os.Remove(link)
		if err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	// Create the symlink
	err := os.Symlink(target, link)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	return nil
}
