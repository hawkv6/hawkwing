package bpf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

var (
	BpffsRoot = "/sys/fs/bpf"
)

func MapToPath(name string) string {
	return filepath.Join(BpffsRoot, name)
}

func Mount() error {
	mounted, err := isMounted(BpffsRoot)
	if err != nil {
		return err
	}
	if mounted {
		return nil
	}
	return mountFS()
}

func isMounted(path string) (bool, error) {
	var stat, pstat unix.Stat_t

	err := unix.Lstat(path, &stat)
	if err != nil {
		if errors.Is(err, unix.ENOENT) {
			// path does not exist -> no mount point
			return false, nil
		}
		return false, fmt.Errorf("could not stat %q: %s", path, err)
	}

	parentDir := filepath.Dir(path)
	err = unix.Lstat(parentDir, &pstat)
	if err != nil {
		return false, fmt.Errorf("could not stat %q: %s", parentDir, err)
	}
	if stat.Dev == pstat.Dev {
		// parent has same device -> no mount point
		return false, nil
	}

	filesytemType := unix.Statfs_t{}
	err = unix.Statfs(path, &filesytemType)
	if err != nil {
		return false, fmt.Errorf("could not statfs %q: %s", path, err)
	}

	return true, nil
}

func mountFS() error {
	bpffsRootStat, err := os.Stat(BpffsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(BpffsRoot, 0755); err != nil {
				return fmt.Errorf("could not create bpf mount directory: %q: %s", BpffsRoot, err)
			}
		} else {
			return fmt.Errorf("could not stat bpf mount directory: %q: %s", BpffsRoot, err)
		}
	} else if !bpffsRootStat.IsDir() {
		return fmt.Errorf("bpf mount directory is not a directory: %q", BpffsRoot)
	}

	if err := unix.Mount(BpffsRoot, BpffsRoot, "bpf", 0, ""); err != nil {
		return fmt.Errorf("could not mount bpf filesystem: %s: %s", BpffsRoot, err)
	}

	return nil
}
