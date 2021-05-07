package main

import "os"

func pathExists(path_name string) bool {
	_, err := os.Stat(path_name)
	return !os.IsNotExist(err)
}
