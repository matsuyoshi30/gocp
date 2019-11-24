package util

import (
	"testing"
)

func TestScrape(t *testing.T) {
	testcase := []struct {
		source  string
		tagtype string
		output  []string
		err     error
	}{
		{
			"",
			"",
			[]string{},
			nil,
		},
		{
			"<html><head></head><body></body></html>",
			"body",
			[]string{},
			nil,
		},
		{
			"<html><head></head><body><pre>1</pre><pre>2</pre></body></html>",
			"pre",
			[]string{"1", "2"},
			nil,
		},
		{
			"<html><head></head><body><tbody><th></th><td>AC</td></tbody></body></html>",
			"tbody",
			[]string{"AC"},
			nil,
		},
		{
			"<html><head><title>ログイン - AtCoder</title></head><body></body></html>",
			"title",
			nil,
			NotLoginError,
		},
		{
			"<html><head><title>Another</title></head><body></body></html>",
			"title",
			nil,
			nil,
		},
	}

	for _, tt := range testcase {
		actual, err := Scrape(tt.source, tt.tagtype)
		if tt.err != err {
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
	LogWrite(-1, "show default log")
}
