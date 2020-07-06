package converter

import (
	"app/pkg/abstract"
	"bufio"
	"io"
	"regexp"
)

const (
	apacheAccessLogDateLayout = "02/Jan/2006:15:04:05 -0700"
)

func ConvertFile(f io.Reader, db abstract.DBLineInserter) error {

	re, err := lineRegExp()
	if err != nil {
		return fmt.Errorf("Cannot compile reg exp : %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, err := parseLine(scanner.Text(), re)
		db.Insert(context.TODO(), line)
	}

	if scanner.Err() != nil {
		return fmt.Errorf("Error while reading file : %w", scanner.Err())
	}
	return nil
}

func lineRegExp() (regexp, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`^(\S+)\s`)                  // 1) IP
	buffer.WriteString(`\S+\s+`)                    // remote logname
	buffer.WriteString(`(?:\S+\s+)+`)               // remote user
	buffer.WriteString(`\[([^]]+)\]\s`)             // 2) date
	buffer.WriteString(`"(\S*)\s?`)                 // 3) method
	buffer.WriteString(`(?:((?:[^"]*(?:\\")?)*)\s`) // 4) URL
	buffer.WriteString(`([^"]*)"\s|`)               // 5) protocol
	buffer.WriteString(`((?:[^"]*(?:\\")?)*)"\s)`)  // 6) or, possibly URL with no protocol
	buffer.WriteString(`(\S+)\s`)                   // 7) status code
	buffer.WriteString(`(\S+)\s`)                   // 8) bytes
	buffer.WriteString(`"((?:[^"]*(?:\\")?)*)"\s`)  // 9) referrer
	buffer.WriteString(`"(.*)"$`)                   // 10) user agent

	return regexp.Compile(buffer.String())
}

func parseLine(re regex, l string) (abstract.Line, error) {

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

	return lineItem
}
