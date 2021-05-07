package main

import "fmt"

type Pile struct {
	content *Download
	next    *Pile
}

type Stack struct {
	top   *Pile
	count int
}

func (self *Stack) init() {
	self.top = nil
	self.count = 0
}

func (self *Stack) createPile(d *Download, next *Pile) *Pile {
	var new_pile *Pile = new(Pile)
	new_pile.content = d
	new_pile.next = next
	return new_pile
}

func (self *Stack) push(d *Download) {
	var new_pile *Pile = self.createPile(d, self.top)
	self.top = new_pile
	self.count++
}

func (self *Stack) pop() (*Download, error) {
	var err error
	if self.top == nil {
		err = fmt.Errorf("Pile is empty")
		return nil, err
	}

	var current_top *Pile = self.top
	self.top = self.top.next
	self.count--

	return current_top.content, err
}

func (self *Stack) peak() *Download {
	var top_download *Download
	if self.top != nil {
		top_download = self.top.content
	}
	return top_download
}
