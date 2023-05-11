package word

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type updateItem struct {
	count int
	file  *os.File
}

type updateData struct {
	items [MaxLength]updateItem
}

func (d *updateData) init() error {
	for i := 1; i < MaxLength; i++ {
		file, err := os.CreateTemp(".", "passman")
		if err != nil {
			d.clear()
			return err
		}
		d.items[i].file = file
	}
	return nil
}

func (d *updateData) clear() {
	for i := 1; i < MaxLength; i++ {
		if d.items[i].file != nil {
			_ = d.items[i].file.Close()
			_ = os.Remove(d.items[i].file.Name())
		}
	}
}

func (d *updateData) readAll(source Source) error {
	return ReadSource(source, d.consume)
}

func (d *updateData) consume(word string) error {
	if word == "" || len(word) >= MaxLength {
		panic(fmt.Errorf("incorrect word '%s'", word))
	}
	index := len(word)
	if _, err := d.items[index].file.WriteString(fmt.Sprintf("%s\n", word)); err != nil {
		return err
	}
	d.items[index].count++
	return nil
}

func (d *updateData) countFor(length int) int {
	return d.items[length].count
}

func (d *updateData) fill(metadata *storageMetadata) {
	offset := pageSize
	for i := 1; i < MaxLength; i++ {
		count := d.countFor(i)
		aligned := alignedSize(i)
		metadata.items[i].count = count
		metadata.items[i].alignedSize = aligned
		if count == 0 {
			metadata.items[i].startOffset = 0
		} else {
			metadata.items[i].startOffset = offset
			offset += (int64(count*aligned)/pageSize + 1) * pageSize
		}
	}
}

func (d *updateData) writeWords(file *os.File, length int, item *metadataItem) error {
	if item.count == 0 {
		return nil
	}
	if _, err := d.items[length].file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	count := 0
	offset := item.startOffset
	scanner := bufio.NewScanner(d.items[length].file)
	for ; scanner.Scan(); count++ {
		word := scanner.Text()
		buffer := make([]byte, item.alignedSize)
		copy(buffer, word)
		if _, err := file.WriteAt(buffer, offset); err != nil {
			return err
		}
		offset += int64(item.alignedSize)
	}
	if count != item.count {
		return fmt.Errorf("%w: expected %d words, got %d", errIncorrectWrite, item.count, count)
	}
	return nil
}

func prepareLine(line string) string {
	if index := strings.IndexRune(line, '#'); index != -1 {
		line = line[:index]
	}
	return strings.TrimSpace(line)
}
