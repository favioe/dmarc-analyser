package dmarc

import (
	"encoding/xml"
	"io"
)

// Decode reads a single DMARC feedback report from r (plain XML).
func Decode(r io.Reader) (*Feedback, error) {
	var f Feedback
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}
