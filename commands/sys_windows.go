package commands

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

const RemoveLockfile = true
const RemoveClientLockfile = true

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}
