package ver

import (
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func TryReadFromGit() (*VersionInformation, error) {
	// get revision
	rev, err := exec.Command("git", "log", "-n", "1", "--pretty=format:%H").Output()

	if err != nil {
		// not a git repository?
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 128 /* not found */ {
					log.Debugf("Not a Git repository; no revision retrieved")

					return nil, nil
				}
			}
		}

		return nil, err
	}

	// done
	vi := MakeVersionInformation()

	vi.Revision = string(rev)

	log.Debugf("Found Git revision: %v", vi.Revision)

	return vi, nil
}
