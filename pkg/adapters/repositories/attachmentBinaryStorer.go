package repositories

import (
	"errors"
	"goticka/pkg/domain/attachment"
	"io/ioutil"
	"log"
	"os"
)

type AttachmentBinaryStorerInterface interface {
	StoreBinary(attachment.Attachment) (attachment.Attachment, error)
	GetBinary(string) ([]byte, error)
}

// ----------------------------------
// FS implementation
// ----------------------------------

type AttachmentBinaryStorerFS struct {
	basePath string
}

func NewAttachmentBinaryStorerFS(basePath string) AttachmentBinaryStorerFS {
	return AttachmentBinaryStorerFS{
		basePath: basePath,
	}
}

func (bs AttachmentBinaryStorerFS) StoreBinary(a attachment.Attachment) (attachment.Attachment, error) {
	storedAttachment := a
	storedAttachment.URI = bs.basePath + a.FileName

	if _, checkFileErr := os.Stat(storedAttachment.URI); checkFileErr == nil {
		return attachment.Attachment{}, errors.New("can not owerwrite " + storedAttachment.URI)
	}

	f, err := os.Create(storedAttachment.URI)
	defer f.Close()

	if err != nil {
		return attachment.Attachment{}, err
	}

	_, err2 := f.Write(a.Raw)

	if err2 != nil {
		return attachment.Attachment{}, err2
	}

	log.Print("Stored attachment " + storedAttachment.URI)
	return storedAttachment, nil
}

func (bs AttachmentBinaryStorerFS) GetBinary(URI string) ([]byte, error) {
	if _, checkFileErr := os.Stat(URI); checkFileErr != nil {
		return []byte{}, errors.New("uri " + URI + " not found")
	}
	raw, err := ioutil.ReadFile(URI)
	if err != nil {
		return []byte{}, err
	}
	return raw, nil
}

// ----------------------------------
// In memory implementation
// ----------------------------------

type AttachmentBinaryStorerMemory struct {
	basePath     string
	inMemoryData map[string][]byte
}

func NewAttachmentBinaryStorerMemory(basePath string) AttachmentBinaryStorerMemory {
	return AttachmentBinaryStorerMemory{
		basePath:     basePath,
		inMemoryData: make(map[string][]byte),
	}
}

func (bs AttachmentBinaryStorerMemory) StoreBinary(a attachment.Attachment) (attachment.Attachment, error) {
	storedAttachment := a
	storedAttachment.URI = bs.basePath + a.FileName

	if _, exists := bs.inMemoryData[storedAttachment.URI]; exists {
		return attachment.Attachment{}, errors.New("can not owerwrite " + storedAttachment.URI)
	}

	bs.inMemoryData[storedAttachment.URI] = a.Raw

	log.Print("Stored attachment " + storedAttachment.URI)
	return storedAttachment, nil
}

func (bs AttachmentBinaryStorerMemory) GetBinary(URI string) ([]byte, error) {
	if raw, exists := bs.inMemoryData[URI]; !exists {
		return []byte{}, errors.New("uri " + URI + " not found")
	} else {
		return raw, nil
	}
}
