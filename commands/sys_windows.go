package commands

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

func workdir() string {
	if programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT); err != nil {
		return `C:\uhppoted`
	} else {
		return filepath.Join(programData, "uhppoted")
	}
}
