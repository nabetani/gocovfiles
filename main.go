package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type options struct {
	Src     string   `json:"src"`
	Root    string   `json:"root"`
	Ignores []string `json:"ignores"`
}

var errHeadLine = errors.New("head line")

const defaultConfigPath = "gocovfiles.json"

func optionsFromJSON(jsonpath string) (*options, error) {
	jsonBytes, err := func() ([]byte, error) {
		file, err := os.Open(jsonpath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		buffer := make([]byte, fi.Size())
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		return buffer, nil
	}()
	if err != nil {
		return nil, err
	}
	o := options{}
	err = json.Unmarshal(jsonBytes, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func getOptions(args []string, exitOnError bool) (*options, error) {
	errorHandling := func() flag.ErrorHandling {
		if exitOnError {
			return flag.ExitOnError
		}
		return flag.ContinueOnError
	}()
	flagSet := flag.NewFlagSet(args[0], errorHandling)
	configPath := defaultConfigPath
	flagSet.StringVar(&configPath, "c", defaultConfigPath, "config JSON path")
	if err := flagSet.Parse(args[1:]); err != nil {
		return nil, err
	}
	o, err := optionsFromJSON(configPath)
	if err != nil {
		return nil, err
	}
	return o, nil
}

type covItem struct {
	fn        string
	sentences int
	covered   bool
}

type covInfo struct {
	items []covItem
}

func convItemFromLine(line string, reg *regexp.Regexp, opts *options) (*covItem, error) {
	m := reg.FindStringSubmatch(line)
	if len(m) < 1 {
		return nil, errHeadLine
	}
	fn, err := filepath.Rel(opts.Root, m[1])
	if err != nil {
		return nil, err
	}
	sentences, err := strconv.Atoi(m[2])
	if err != nil {
		return nil, err
	}
	covered, err := strconv.Atoi(m[3])
	if err != nil {
		return nil, err
	}
	return &covItem{
		fn:        fn,
		sentences: sentences,
		covered:   covered != 0,
	}, nil
}

func stringArrayContains(a []string, s string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func covInfoFromReader(src io.Reader, opts *options) (*covInfo, error) {
	scanner := bufio.NewScanner(src)
	//github.com/nabetani/gocovfiles/samplesrc/hoge.go:15.18,17.2 1 1
	reg := regexp.MustCompile(`^(.*)\:\d+\.\d+\,\d+\.\d+\s+(\d+)\s+(\d+)\s*$`)
	c := covInfo{}
	for scanner.Scan() {
		item, err := convItemFromLine(scanner.Text(), reg, opts)
		if err == errHeadLine {
			continue
		}
		if err != nil {
			return nil, err
		}
		if stringArrayContains(opts.Ignores, item.fn) {
			continue
		}
		c.items = append(c.items, *item)
	}

	if err := scanner.Err(); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func covInfoFromFilename(opts *options) (*covInfo, error) {
	file, err := os.Open(opts.Src)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return covInfoFromReader(file, opts)
}

type sumItem struct {
	covered    int
	notCovered int
}

func (s *sumItem) percent() float64 {
	total := s.covered + s.notCovered
	if total == 0 {
		return 0
	}
	return float64(s.covered) * 100 / float64(total)
}

type summary struct {
	sumItems map[string]*sumItem
	total    sumItem
}

func (s *summary) toString() string {
	fns := []string{}
	fncol := "filename"
	maxlen := len(fncol)
	for fn := range s.sumItems {
		fns = append(fns, fn)
		if lenFn := len(fn); maxlen < lenFn {
			maxlen = lenFn
		}
	}
	sort.Strings(fns)
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%*s  covered  not covered   ratio\n", maxlen, fncol))
	horzLine := strings.Repeat("-", maxlen) + "  -------  -----------  ------\n"
	b.WriteString(horzLine)
	format := "%*s  %7d  %11d  %5.1f%%\n"
	for _, fn := range fns {
		i := s.sumItems[fn]
		b.WriteString(fmt.Sprintf(format, maxlen, fn, i.covered, i.notCovered, i.percent()))
	}
	b.WriteString(horzLine)
	b.WriteString(fmt.Sprintf(format, maxlen, "total", s.total.covered, s.total.notCovered, s.total.percent()))
	return b.String()
}

func summarize(c *covInfo, opts *options) (*summary, error) {
	s := summary{sumItems: map[string]*sumItem{}}
	for _, ci := range c.items {
		if _, ok := s.sumItems[ci.fn]; !ok {
			s.sumItems[ci.fn] = &sumItem{}
		}
		if ci.covered {
			s.sumItems[ci.fn].covered += ci.sentences
			s.total.covered += ci.sentences
		} else {
			s.sumItems[ci.fn].notCovered += ci.sentences
			s.total.notCovered += ci.sentences
		}
	}
	return &s, nil
}

func main() {
	opts, err := getOptions(os.Args, true)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	c, err := covInfoFromFilename(opts)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	sum, err := summarize(c, opts)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	fmt.Println(sum.toString())
}
