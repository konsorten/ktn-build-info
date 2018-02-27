package ver

import (
	"os/exec"
)

func TryReadFromGit() (*VersionInformation, error) {
	// get revision
	rev, err := exec.Command("git", "log", "-n", "1", "--pretty=format:\"%h\"").Output()

	if err != nil {
		return nil, err
	}

	// done
	vi := MakeVersionInformation()

	vi.Revision = string(rev)

	return vi, nil
}
