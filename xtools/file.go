package xtools

import "os"

//FileExists FileExists
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return false
		}
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
