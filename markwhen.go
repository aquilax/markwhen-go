package markwhen

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	keyValueSeparator  = ":"
	dateRangeSeparator = "-"
	pageBreak          = "_-_-_break_-_-_"

	prefixComment     = "//"
	prefixTitle       = "title:"
	prefixDescription = "description:"
	prefixDateFormat  = "dateFormat:"

	usDateFormat = "01/02/2006" // MM/dd/yy
	euDateFormat = "02/01/2006" // MM/dd/yy
)

var knownDateFormats = map[string]string{
	"MM/dd/yy": usDateFormat,
	"d/M/y":    euDateFormat,
}

const DefaultDateFormat = usDateFormat

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
		DateFormat: DefaultDateFormat,
		Tags:       make(Tags),
	}
}

type Event struct {
	From time.Time
	To   time.Time
	Body string
}

type Page struct {
	Header *Header
	Events []*Event
}

func NewPage() *Page {
	return &Page{
		NewHeader(),
		make([]*Event, 0),
	}
}

type MarkWhen struct {
	Pages []*Page
}

func Parse(reader io.Reader) (*MarkWhen, error) {
	var err error
	var line string
	var trimmedLine string
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	pages := make([]*Page, 0)
	page := NewPage()
	inHeader := true
	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()
		trimmedLine = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefixComment) {
			// ignore line comments
			continue
		}
		if trimmedLine == "" {
			// ignore empty lines
			continue
		}
		if line == pageBreak {
			oldDateFormat := page.Header.DateFormat
			pages = append(pages, page)
			page = NewPage()
			page.Header.DateFormat = oldDateFormat
			continue
		}
		if inHeader {
			if strings.HasPrefix(line, prefixTitle) {
				if _, page.Header.Title, err = getKeyValue(line); err != nil {
					return nil, err
				}
				continue
			}
			if strings.HasPrefix(line, prefixDescription) {
				if _, page.Header.Description, err = getKeyValue(line); err != nil {
					return nil, err
				}
				continue
			}
			if strings.HasPrefix(line, prefixDateFormat) {
				if page.Header.DateFormat, err = getDateFormatValue(line); err != nil {
					return nil, err
				}
				continue
			}
		}
		if event, err := getEvent(line, page.Header); err != nil {
			return nil, err
		} else {
			page.Events = append(page.Events, event)
			continue
		}
	}
	pages = append(pages, page)
	return &MarkWhen{pages}, nil
}

func getKeyValue(line string) (string, string, error) {
	index := strings.Index(line, keyValueSeparator)
	if index == -1 {
		return "", "", fmt.Errorf("invalid key value")
	}
	return strings.TrimSpace(line[:index]), strings.TrimSpace(line[index+1:]), nil
}

func getDateFormatValue(line string) (string, error) {
	if _, value, err := getKeyValue(line); err != nil {
		return DefaultDateFormat, nil
	} else {
		trimmedValue := strings.TrimSpace(value)
		if format, found := knownDateFormats[trimmedValue]; found {
			return format, nil
		} else {
			return "", fmt.Errorf("unknown dateFormat: %s", trimmedValue)
		}
	}
}

func getEvent(line string, header *Header) (*Event, error) {
	key, value, err := getKeyValue(line)
	if err != nil {
		return nil, err
	}
	from, to, err := getRange(key, header)
	if err != nil {
		return nil, err
	}
	return &Event{from, to, value}, nil
}

func parseTime(dateFormat string, t string) (time.Time, error) {
	trimmedLine := strings.TrimSpace(t)
	if trimmedLine == "now" {
		return time.Time{}, nil
	}
	return time.Parse(dateFormat, trimmedLine)
}

func getRange(dateRange string, header *Header) (time.Time, time.Time, error) {
	var err error
	fromTime := time.Time{}
	toTime := time.Time{}
	index := strings.Index(dateRange, dateRangeSeparator)
	if index == -1 {
		// single date
	} else {
		// date range
		fromTime, err = parseTime(header.DateFormat, dateRange[:index])
		if err != nil {
			return fromTime, toTime, err
		}
		toTime, err = parseTime(header.DateFormat, dateRange[index+1:])
		if err != nil {
			return fromTime, toTime, err
		}
	}
	return fromTime, toTime, nil
}
