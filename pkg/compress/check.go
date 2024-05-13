package compress

import (
	"os/exec"
)

func _CommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func SevenZipAvailable() bool {
	return _CommandAvailable("7z")
}

func ZipAvailable() bool {
	return _CommandAvailable("zip") && _CommandAvailable("zipnote")
}

func TarAvailable() bool {
	return _CommandAvailable("tar")
}
