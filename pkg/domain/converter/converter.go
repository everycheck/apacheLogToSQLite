package converter

import (
	"app/pkg/abstract"
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

const (
	apacheAccessLogDateLayout = "02/Jan/2006:15:04:05 -0700"

	regExpIp                          = `^(\S+)\s`
	regExpRemoteLogName               = `\S+\s+`
	regExpRemoteUser                  = `(?:\S+\s+)+`
	regExpDate                        = `\[([^]]+)\]\s`
	regExpMethod                      = `"(\S*)\s?`
	regExpUrl                         = `(?:((?:[^"]*(?:\\")?)*)\s`
	regExpProtocol                    = `([^"]*)"\s|`
	regExpOptionnalUrlWithoutProtocol = `((?:[^"]*(?:\\")?)*)"\s)`
	regExpStatusCode                  = `(\S+)\s`
	regExpBytes                       = `(\S+)\s`
	regExpReferres                    = `"((?:[^"]*(?:\\")?)*)"\s`
	regExpUserAgent                   = `"(.*)"$`

	lineRegExp = regExpIp + regExpRemoteLogName + regExpRemoteUser +
		regExpDate + regExpMethod + regExpUrl + regExpProtocol +
		regExpOptionnalUrlWithoutProtocol + regExpStatusCode + regExpBytes +
		regExpReferres + regExpUserAgent
)

func ConvertFile(f io.Reader, db abstract.DBLineBulkInserter, batch int) error {

	re, err := compileLineRegExp()
	if err != nil {
		return fmt.Errorf("Cannot compile reg exp : %w", err)
	}

	lines := make([]abstract.Line, 0, batch)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, err := parseLine(re, scanner.Text())
		if err != nil {
			return fmt.Errorf("Error while parsing line : %w", err)
		}
		lines = append(lines, line)
		if len(lines) == batch {
			err = db.BulkInsert(context.TODO(), lines)
			if err != nil {
				return fmt.Errorf("Cannot insert in db : %w", err)
			}
			lines = make([]abstract.Line, 0, batch)
		}
	}

	if len(lines) > 0 {
		err = db.BulkInsert(context.TODO(), lines)
		if err != nil {
			return fmt.Errorf("Cannot insert in db : %w", err)
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("Error while reading file : %w", scanner.Err())
	}
	return nil
}

func compileLineRegExp() (*regexp.Regexp, error) {
	return regexp.Compile(lineRegExp)
}

func parseLine(re *regexp.Regexp, l string) (abstract.Line, error) {

	var err error

	result := re.FindStringSubmatch(l)

	lineItem := abstract.Line{}
	lineItem.RemoteHost = result[1]
	lineItem.Time, err = time.Parse(apacheAccessLogDateLayout, result[2])
	if err != nil {
		return lineItem, fmt.Errorf("Cannont parse date %s : %w", result[2], err)
	}
	lineItem.Request = result[3] + " " + result[4] + " " + result[5]
	lineItem.Status, err = strconv.Atoi(result[7])
	if err != nil {
		lineItem.Status = 0
	}
	lineItem.Bytes, err = strconv.Atoi(result[8])
	if err != nil {
		lineItem.Bytes = 0
	}
	lineItem.Referer = result[9]
	lineItem.UserAgent = result[10]
	lineItem.URL = result[4]
	if len(lineItem.URL) == 0 {
		lineItem.URL = result[6]
	}

	return lineItem, nil
}
