package converter

import (
	"app/pkg/abstract"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRegExpBuild(t *testing.T) {
	_, err := compileLineRegExp()
	if err != nil {
		t.Fatalf("regexp should compile %v", err)
	}
}

func TestParseLine(t *testing.T) {
	re, _ := compileLineRegExp()
	tests := []struct {
		input    string
		expected abstract.Line
	}{
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: abstract.Line{
				RemoteHost: "51.255.43.108",
				Time:       time.Date(2020, 7, 5, 6, 26, 0, 0, time.Local),
				Request:    "POST /tokens/jwt HTTP/1.1",
				Status:     201,
				Bytes:      3472,
				Referer:    "-",
				UserAgent:  "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28",
				URL:        "/tokens/jwt",
			},
		},
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 20z1 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: abstract.Line{
				RemoteHost: "51.255.43.108",
				Time:       time.Date(2020, 7, 5, 6, 26, 0, 0, time.Local),
				Request:    "POST /tokens/jwt HTTP/1.1",
				Status:     0,
				Bytes:      3472,
				Referer:    "-",
				UserAgent:  "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28",
				URL:        "/tokens/jwt",
			},
		},
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3z472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: abstract.Line{
				RemoteHost: "51.255.43.108",
				Time:       time.Date(2020, 7, 5, 6, 26, 0, 0, time.Local),
				Request:    "POST /tokens/jwt HTTP/1.1",
				Status:     201,
				Bytes:      0,
				Referer:    "-",
				UserAgent:  "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28",
				URL:        "/tokens/jwt",
			},
		},
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: abstract.Line{
				RemoteHost: "51.255.43.108",
				Time:       time.Date(2020, 7, 5, 6, 26, 0, 0, time.Local),
				Request:    "POST  ",
				Status:     201,
				Bytes:      3472,
				Referer:    "-",
				UserAgent:  "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28",
				URL:        "/tokens/jwt",
			},
		},
	}

	for i, tt := range tests {
		result, err := parseLine(re, tt.input)
		if err != nil {
			t.Fatalf("[%d] error while parsing %s : %v", i, tt.input, err)
		}
		if result != tt.expected {
			t.Fatalf("[%d] error result is not equal to expected\n exp %#v\n got %#v\n, ", i, tt.expected, result)
		}
	}
}

func TestParseLineErrors(t *testing.T) {
	re, _ := compileLineRegExp()
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `51.255.43.108 - - [05/Jul/20:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: "Cannont parse date 05/Jul/20:06:26:00 +0200 : parsing time \"05/Jul/20:06:26:00 +0200\" as \"02/Jan/2006:15:04:05 -0700\": cannot parse \"6:26:00 +0200\" as \"2006\"",
		},
	}

	for i, tt := range tests {
		_, err := parseLine(re, tt.input)
		if err == nil {
			t.Fatalf("[%d] should have failed while parsing %s", i, tt.input)
		}
		if err.Error() != tt.expected {
			t.Fatalf("[%d] error result is not equal to expected\n exp %#v\n got %#v\n, ", i, tt.expected, err.Error())
		}
	}
}

type fakeInserter struct {
	called       int
	shouldFailed error
}

func (fi *fakeInserter) Insert(ctx context.Context, line abstract.Line) error {
	fi.called++
	if fi.shouldFailed != nil {
		return fi.shouldFailed
	}
	return nil
}

func TestConvertFileSuccess(t *testing.T) {

	tests := []struct {
		input    string
		expected int
	}{
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"
51.255.43.108 - - [05/Jul/2020:06:26:01 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"
51.255.43.108 - - [05/Jul/2020:06:26:02 +0200] "POST /emails/process/take HTTP/1.1" 404 240 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"
51.255.43.108 - - [05/Jul/2020:06:27:01 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"
51.255.43.108 - - [05/Jul/2020:06:27:01 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: 5,
		},

		{
			input:    `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
			expected: 1,
		},
	}

	for i, tt := range tests {
		fi := &fakeInserter{}
		err := ConvertFile(strings.NewReader(tt.input), fi)
		if err != nil {
			t.Fatalf("[%d] error while parsing %s : %v", i, tt.input, err)
		}

		if fi.called != tt.expected {
			t.Fatalf("[%d] error Insert should have been called %d got %d, ", i, tt.expected, fi.called)
		}

	}
}

func TestConvertFileCannotInsert(t *testing.T) {

	tests := []struct {
		input string
	}{
		{
			input: `51.255.43.108 - - [05/Jul/2020:06:26:00 +0200] "POST /tokens/jwt HTTP/1.1" 201 3472 "-" "GuzzleHttp/6.5.1 curl/7.52.1 PHP/7.1.33-8+0~20200202.31+debian9~1.gbp266c28"`,
		},
	}

	for i, tt := range tests {
		errStr := "Cannont insert in test"
		packerErrStr := "Cannot insert in db : " + errStr
		fi := &fakeInserter{shouldFailed: fmt.Errorf(errStr)}
		err := ConvertFile(strings.NewReader(tt.input), fi)
		if err == nil {
			t.Fatalf("[%d] should have failed while parsing %s", i, tt.input)
		}
		if err.Error() != packerErrStr {
			t.Fatalf("[%d] error result is not equal to expected\n exp %#v\n got %#v\n, ", i, packerErrStr, err.Error())
		}
	}
}
