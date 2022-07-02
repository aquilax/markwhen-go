package markwhen

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func marshal(mw *MarkWhen) string {
	b, _ := json.Marshal(mw)
	return string(b)
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
			&MarkWhen{&Header{Title: "This is a title", Tags: make(Tags)}},
			false,
		},
		{
			"parses description",
			"description: This is a description",
			&MarkWhen{&Header{Description: "This is a description", Tags: make(Tags)}},
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
