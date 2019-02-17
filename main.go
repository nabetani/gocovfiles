package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type options struct {
	src  string
	root string
}

func parseCommandLine(args []string, exitOnError bool) (*options, error) {
	errorHandling := func() flag.ErrorHandling {
		if exitOnError {
			return flag.ExitOnError
		}
		return flag.ContinueOnError
	}()
	flagSet := flag.NewFlagSet(args[0], errorHandling)
	o := options{}
	flagSet.StringVar(&o.src, "src", "cover.out", "input file name")
	flagSet.StringVar(&o.root, "root", "", "root directory of files")
	err := flagSet.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	return &o, nil
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
		return &covItem{}, nil
	}
	fn, err := filepath.Rel(opts.root, m[1])
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

func covInfoFromReader(src io.Reader, opts *options) (*covInfo, error) {
	scanner := bufio.NewScanner(src)
	//github.com/nabetani/gocovfiles/samplesrc/hoge.go:15.18,17.2 1 1
	reg := regexp.MustCompile(`^(.*)\:\d+\.\d+\,\d+\.\d+\s+(\d+)\s+(\d+)\s*$`)
	c := covInfo{}
	for scanner.Scan() {
		item, err := convItemFromLine(scanner.Text(), reg, opts)
		if err != nil {
			return nil, err
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
	file, err := os.Open(opts.src)
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

type fileSumItem struct {
	sumItem
	fn string
}

type summary struct {
	sumItems []fileSumItem
	total    sumItem
}

func summarize(c *covInfo, opts *options) (*summary, error) {
	files := map[string]bool{}
	s := summary{}
	for _, ci := range c.items {
		if !files[ci.fn] {
			s.sumItems = append(s.sumItems, fileSumItem{fn: ci.fn})
		}
	}
	return &s, nil
}

func main() {
	opts, err := parseCommandLine(os.Args, true)
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
	log.Println(sum)
}
