package markwhen

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	keyValueSeparator = ":"

	prefixComment = "//"
	prefixTitle   = "title:"
)

type Color = string

type Tags map[string]Color

type Header struct {
	Title       string
	Description string
	DateFormat  string
	Tags        Tags
}

func NewHeader() *Header {
	return &Header{
		Tags: make(Tags),
	}
}

type MarkWhen struct {
	Header *Header
}

func Parse(reader io.Reader) (*MarkWhen, error) {
	var err error
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	var line string
	header := NewHeader()
	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()
		if strings.HasPrefix(line, prefixComment) {
			continue
		}
		if strings.HasPrefix(line, prefixTitle) {
			if header.Title, err = getTitle(line); err != nil {
				return nil, err
			}
		}
	}
	return &MarkWhen{header}, nil
}

func getMetaValue(line string) (string, error) {
	index := strings.Index(line, keyValueSeparator)
	if index == -1 {
		return "", fmt.Errorf("invalid meta value")
	}
	return strings.TrimSpace(line[index+1:]), nil
}

func getTitle(line string) (string, error) {
	return getMetaValue(line)
}
