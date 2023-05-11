package word

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"
)

type StorageConfig interface {
	WordStorageFile() string
}

type defaultStorage struct {
	file     *os.File
	metadata storageMetadata
}

type metadataItem struct {
	count       int
	alignedSize int
	startOffset int64
}

type storageMetadata struct {
	code    [8]byte
	version int
	items   [MaxLength]metadataItem
}

const pageSize = int64(4096)
const metadataSize = int(unsafe.Sizeof(storageMetadata{}))

var code = [8]byte{0x01, 0x22, 0x3, 0x44, 0x5, 0x66, 0x7, 0x88}

func (s *defaultStorage) Words(length int) (int, error) {
	if length <= 0 || length >= MaxLength {
		return 0, ErrIncorrectLength
	}
	return s.metadata.items[length].count, nil
}

func (s *defaultStorage) WordAt(length, index int) (string, error) {
	if length <= 0 || length >= MaxLength {
		return "", ErrIncorrectLength
	}
	if index < 0 || index >= s.metadata.items[length].count {
		return "", ErrNoSuchWord
	}
	if s.metadata.items[length].startOffset == 0 {
		return "", ErrNoSuchWord
	}
	buffer := make([]byte, s.metadata.items[length].alignedSize)
	offset := s.metadata.items[length].startOffset + int64(s.metadata.items[length].alignedSize*index)
	n, err := s.file.ReadAt(buffer, offset)
	switch {
	case err != nil:
		return "", err
	case n != s.metadata.items[length].alignedSize:
		return "", errIncorrectRead
	}
	return string(buffer[:length]), nil
}

func (s *defaultStorage) Close() error {
	return s.file.Close()
}

func (s *defaultStorage) init() error {
	buffer := (*[metadataSize]byte)(unsafe.Pointer(&s.metadata))
	n, err := s.file.ReadAt(buffer[:], 0)
	switch {
	case err != nil:
		return err
	case n != metadataSize:
		return errIncorrectRead
	case !bytes.Equal(code[:], s.metadata.code[:]):
		return fmt.Errorf("%w %s", ErrIncorrectFile, s.file.Name())
	}
	return nil
}

func (i *metadataItem) String() string {
	return fmt.Sprintf("metadataItem[count=%d,alignedSize=%d,startOffset=%d]",
		i.count, i.alignedSize, i.startOffset)
}

func NewStorage(config StorageConfig) (CloseableStorage, error) {
	return NewStorageNamed(config.WordStorageFile())
}

func NewStorageNamed(name string) (CloseableStorage, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	result := defaultStorage{file: file}
	if err := result.init(); err != nil {
		_ = file.Close()
		return nil, err
	}
	return &result, nil
}

func UpdateStorage(name string, source Source) (CloseableStorage, error) {
	data := updateData{}
	if err := data.init(); err != nil {
		return nil, err
	}
	defer data.clear()
	if err := data.readAll(source); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	result := defaultStorage{file: file}
	copy(result.metadata.code[:], code[:])
	result.metadata.version = 1
	data.fill(&result.metadata)
	buffer := (*[metadataSize]byte)(unsafe.Pointer(&result.metadata))
	if _, err := file.WriteAt(buffer[:], 0); err != nil {
		return nil, err
	}
	for i := 1; i < MaxLength; i++ {
		if err := data.writeWords(result.file, i, &result.metadata.items[i]); err != nil {
			return nil, err
		}
	}
	return &result, nil
}
