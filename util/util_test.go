package util

import (
	"testing"
)

func TestScrape(t *testing.T) {
	testcase := []struct {
		source  string
		tagtype string
		output  []string
	}{
		{
			"",
			"",
			[]string{},
		},
		{
			"<html><head></head><body></body></html>",
			"body",
			[]string{},
		},
		{
			"<html><head></head><body><pre>1</pre><pre>2</pre></body></html>",
			"pre",
			[]string{"1", "2"},
		},
		{
			"<html><head></head><body><tbody><th></th><td>AC</td></tbody></body></html>",
			"tbody",
			[]string{"AC"},
		},
	}

	for _, tt := range testcase {
		actual, err := Scrape(tt.source, tt.tagtype)
		if err != nil {
			t.Errorf("Error %v\n", err)
		}

		if len(tt.output) != len(actual) {
			t.Errorf("Length Error: expected %v, but got %v\n", len(tt.output), len(actual))
		}

		for idx, elem := range actual {
			if elem != tt.output[idx] {
				t.Errorf("Element Error: expected %v, but got %v\n", tt.output[idx], elem)
			}
		}
	}
}

func TestLogWrite(t *testing.T) {
	LogWrite(SUCCESS, "show success log")
	LogWrite(FAILED, "show failed log")
	LogWrite(INFO, "show info log")
	LogWrite(INFO, "show default log")
}
