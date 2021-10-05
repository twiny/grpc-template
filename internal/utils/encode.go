package utils

import (
	"bytes"
	"encoding/gob"
	"phonebook/internal/types"
	"regexp"
)

//
var offsetRegexp = regexp.MustCompile(`^UTC([+-][\d]{1,2})$`)

// encode
func Encode(contact *types.Contact) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(contact); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decode
func Decode(data []byte) (*types.Contact, error) {
	var contact types.Contact
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&contact); err != nil {
		return nil, err
	}

	return &contact, nil
}
