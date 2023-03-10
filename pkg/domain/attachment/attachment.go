package attachment

import (
	"errors"
	"time"
)

type Attachment struct {
	ID          int64
	URI         string
	FileName    string
	ContentType string
	Size        int
	Raw         []byte

	Created time.Time
	Changed time.Time
	Deleted time.Time
}

func (a Attachment) Validate() error {
	if a.URI == "" {
		return errors.New("missing 'URI' in attachment")
	}
	if a.FileName == "" {
		return errors.New("missing 'filename' in attachment")
	}
	if a.ContentType == "" {
		return errors.New("missing 'contentType' in attachment")
	}
	if a.Size == 0 {
		return errors.New("missing 'size' in attachment")
	}
	return nil
}
