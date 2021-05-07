package main

import "fmt"

type Download struct {
	Name         string
	Directory    string
	Files        []string
	file_pointer int
}

func (self *Download) init(name string, directory string, files []string) {
	self.Name = name
	self.Directory = directory
	self.Files = files
	self.file_pointer = -1
}

func (self *Download) getNext() (string, error) {
	self.file_pointer++
	if self.file_pointer == len(self.Files) {
		return "", fmt.Errorf("Downloads stack ended")
	}
	return self.Files[self.file_pointer], nil
}

func (self *Download) resetFileStream() {
	self.file_pointer = -1
}
