package email

import (
	"io/ioutil"
	"path/filepath"
)

// Attachment represents an email attachment.
type Attachment struct {
	Name   string
	Data   []byte
	Inline bool
}

func NewAttachment(name string, inline bool) (*Attachment, error) {
	var err error
	attach := &Attachment{}
	attach.Data, err = ioutil.ReadFile(name)
	_, attach.Name = filepath.Split(name)
	return attach, err
}
