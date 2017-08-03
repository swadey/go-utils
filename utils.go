package utils

import (
	"os"
	"io"
	"fmt"
	"strconv"
	docopt "github.com/docopt/docopt-go"
	xz "github.com/remyoudompheng/go-liblzma"
	"compress/gzip"
	"compress/bzip2"
	"regexp"
)

// ----------------------------------------------------------------------------------------------------------------
// Utilities
// ----------------------------------------------------------------------------------------------------------------
var filenamePattern = regexp.MustCompile(`^.*\.(gz|bz2|xz)$`)

func Zopen(fn string) (io.Reader, *os.File) {
	compress := filenamePattern.FindStringSubmatch(fn)
	if len(compress) == 0 {
		in, _ := os.Open(fn)
		return in, in
	} else {
		in, _ := os.Open(fn)
		var reader io.Reader

		switch compress[1] {
		case "gz":
			gz, _ := gzip.NewReader(in)
			reader = gz
		case "bz2":
			bz := bzip2.NewReader(in)
			reader = bz
		case "xz":
			xz, _ := xz.NewReader(in)
			reader = xz
		}

		return reader, in
	}
}

func Zcreate(fn string) (io.Writer, *os.File) {
	compress := filenamePattern.FindStringSubmatch(fn)
	if len(compress) == 0 {
		outf, _	:= os.Create(fn)
		return outf, outf
	} else {
		outf, _	:= os.Create(fn)
		var writer io.Writer

		switch compress[1] {
		case "gz":
			gz := gzip.NewWriter(outf)
			writer = gz
		case "bz2":
			fmt.Fprintf(os.Stderr, "error: bzip2 is not supported")
			writer = outf
		case "xz":
			xz, _ := xz.NewWriter(outf, xz.Level2)
			writer = xz
		}

		return writer, outf
	}
}

// ----------------------------------------------------------------------------------------------------------------
// DocOpt wrapper
// ----------------------------------------------------------------------------------------------------------------
type Args struct {
	raw map[string] interface{}
}

func (args *Args) Bool(key string) bool {
	return args.raw[key].(bool)
}

func (args *Args) String(key string) string {
	return args.raw[key].(string)
}

func (args *Args) Int(key string) int {
	i, _ := strconv.ParseInt(args.raw[key].(string), 10, 0)
	return int(i)
}

func (args *Args) Hex(key string) int {
	i, _ := strconv.ParseInt(args.raw[key].(string), 16, 0)
	return int(i)
}

func (args *Args) Float(key string) float64 {
	f, _ := strconv.ParseFloat(args.raw[key].(string), 64)
	return f
}

func Docopt(usage string, version string) *Args {
	args, _	:= docopt.Parse(usage, nil, true, version, false)
	return &Args{args}
}
