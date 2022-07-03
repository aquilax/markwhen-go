package markwhen

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func marshal(mw *MarkWhen) string {
	b, _ := json.Marshal(mw)
	return string(b)
}

func mustParseTime(t string) time.Time {
	result, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return result
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *MarkWhen
		wantErr bool
	}{
		{
			"parses title",
			"title: This is a title",
			&MarkWhen{
				[]*Page{
					{
						&Header{Title: "This is a title", DateFormat: DefaultDateFormat, Tags: make(Tags)},
						[]*Collection{},
					},
				},
			},
			false,
		},
		{
			"parses description",
			"description: This is a description",
			&MarkWhen{
				[]*Page{
					{
						&Header{Description: "This is a description", DateFormat: DefaultDateFormat, Tags: make(Tags)},
						[]*Collection{},
					},
				},
			},
			false,
		},
		{
			"EU Date example",
			`// To indicate we are using European date formatting
dateFormat: d/M/y

// 2 weeks
01/01/2023 - 14/01/2023: Phase 1 #Exploratory

// Another 2 weeks
15/01/2023 - 31/01/2023: Phase 2 #Implementation

// 3 days, after a one week buffer
07/03/2023 - 10/03/2023: Phase 4 - kickoff! #Launch
`,
			&MarkWhen{
				[]*Page{
					{
						&Header{
							DateFormat: euDateFormat,
							Tags:       make(Tags),
						},
						[]*Collection{
							{
								Type: CollectionFree,
								Events: []*Event{
									{From: mustParseTime("2023-01-01T00:00:00Z"), To: mustParseTime("2023-01-14T00:00:00Z"), Body: "Phase 1 #Exploratory"},
									{From: mustParseTime("2023-01-15T00:00:00Z"), To: mustParseTime("2023-01-31T00:00:00Z"), Body: "Phase 2 #Implementation"},
									{From: mustParseTime("2023-03-07T00:00:00Z"), To: mustParseTime("2023-03-10T00:00:00Z"), Body: "Phase 4 - kickoff! #Launch"},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"works with multiple pages",
			`title: This is a title for page 1
_-_-_break_-_-_
title: This is a title for page 2`,
			&MarkWhen{
				[]*Page{
					{
						&Header{Title: "This is a title for page 1", DateFormat: DefaultDateFormat, Tags: make(Tags)},
						[]*Collection{},
					},
					{
						&Header{Title: "This is a title for page 2", DateFormat: DefaultDateFormat, Tags: make(Tags)},
						[]*Collection{},
					},
				},
			},
			false,
		},
		{
			"handles groups",
			`group potato
endGroup`,
			&MarkWhen{
				[]*Page{
					{
						&Header{DateFormat: DefaultDateFormat, Tags: make(Tags)},
						[]*Collection{NewCollection(CollectionGroup)},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := Parse(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = \n%+v\n, want \n%+v\n", marshal(got), marshal(tt.want))
			}
		})
	}
}
