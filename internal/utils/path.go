package utils

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetFilePath(elem ...string) (string, error) {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	rootPath := strings.TrimSpace(string(cmdOut))
	pathList := []string{rootPath}
	pathList = append(pathList, elem...)
	return filepath.Join(pathList...), nil
}
