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
	"github.com/ttacon/chalk"
	"time"
	"github.com/briandowns/spinner"
	//"github.com/gosuri/uilive"
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

// ----------------------------------------------------------------------------------------------------------------
// Loggers and Spinners
// ----------------------------------------------------------------------------------------------------------------
var info_color  = chalk.White.NewStyle()
var debug_color = chalk.Cyan.NewStyle()
var warn_color  = chalk.Magenta.NewStyle()
var error_color = chalk.Red.NewStyle()

type Progress interface {
	Write(p []byte) (n int, err error)
	Stop()
}

type ExtSpin struct {
	spinner.Spinner
}
// patch spinner
func (spin *ExtSpin) Write(p []byte) (n int, err error) {
	spin.Suffix = string(p)
	return len(p), nil
}
var Spinners = spinner.CharSets

func logf(format string, args... interface{}) {
	fmt.Printf(format, args...)
}

func prefix(tag string) string {
	return fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05.000") + " %-10s ", tag)
}

func logger(tag string, color chalk.Style, format string, args... interface{}) {
	adjfmt := prefix(tag) + format + "\n"
	logf(color.Style(adjfmt), args...)
}

func Info(format string, args... interface{}) {
	logger("[INFO]", info_color, format, args...)
}

func Debug(format string, args... interface{}) {
	logger("[DEBUG]", debug_color, format, args...)
}

func Warn(format string, args... interface{}) {
	logger("[WARN]", warn_color, format, args...)
}

func Error(format string, args... interface{}) {
	logger("[ERROR]", error_color, format, args...)
}

func SpinCustom(spin []string, interval time.Duration, finalizer string) Progress {
	s := spinner.New(spin, interval)
	s.Color("yellow")
	s.Prefix = prefix("[RUNNING]")
	s.FinalMSG = finalizer
	s.Start()
	return &ExtSpin{*s}
}

func Spin(spin []string) Progress {
	return SpinCustom(spin, 100*time.Millisecond, " completed.\n")
}

