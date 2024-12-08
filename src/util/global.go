package util

import (
	"os"
	"path"
)

func GetWdFile(name string) string {
	dir, _ := os.Getwd()
	dir = path.Join(dir, name)
	return dir
}
