package repositories

import (
	"errors"
	"goticka/pkg/domain/attachment"
	"io/ioutil"
	"log"
	"os"
)

type BinaryRepositoryInterface interface {
	StoreBinary(attachment.Attachment) (attachment.Attachment, error)
	GetBinary(string) ([]byte, error)
}

// ----------------------------------
// FS implementation
// ----------------------------------

type BinaryRepositoryFS struct {
	basePath string
}

func NewBinaryRepositoryFS(basePath string) BinaryRepositoryFS {
	return BinaryRepositoryFS{
		basePath: basePath,
	}
}

func (br BinaryRepositoryFS) StoreBinary(a attachment.Attachment) (attachment.Attachment, error) {
	storedAttachment := a
	storedAttachment.URI = br.basePath + a.FileName

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

func (br BinaryRepositoryFS) GetBinary(URI string) ([]byte, error) {
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

type BinaryRepositoryMemory struct {
	basePath     string
	inMemoryData map[string][]byte
}

func NewBinaryRepositoryMemory(basePath string) BinaryRepositoryMemory {
	return BinaryRepositoryMemory{
		basePath:     basePath,
		inMemoryData: make(map[string][]byte),
	}
}

func (br BinaryRepositoryMemory) StoreBinary(a attachment.Attachment) (attachment.Attachment, error) {
	storedAttachment := a
	storedAttachment.URI = br.basePath + a.FileName

	if _, exists := br.inMemoryData[storedAttachment.URI]; exists {
		return attachment.Attachment{}, errors.New("can not owerwrite " + storedAttachment.URI)
	}

	br.inMemoryData[storedAttachment.URI] = a.Raw

	log.Print("Stored attachment " + storedAttachment.URI)
	return storedAttachment, nil
}

func (br BinaryRepositoryMemory) GetBinary(URI string) ([]byte, error) {
	if raw, exists := br.inMemoryData[URI]; !exists {
		return []byte{}, errors.New("uri " + URI + " not found")
	} else {
		return raw, nil
	}
}
