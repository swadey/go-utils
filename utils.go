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
	"strings"
)

// ----------------------------------------------------------------------------------------------------------------
// Kontstants
// ----------------------------------------------------------------------------------------------------------------
var info_color			= chalk.White.NewStyle()
var debug_color			= chalk.Cyan.NewStyle()
var warn_color			= chalk.Magenta.NewStyle()
var error_color			= chalk.Red.NewStyle()
var spin_color			= chalk.Yellow.NewStyle()
var complete_color  = chalk.Green.NewStyle()

var Spinners = spinner.CharSets
const complete_symbol = "✔"
const clear = "\033[2K\033[1A\033[J"

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
// Spinner
// ----------------------------------------------------------------------------------------------------------------
type Spinner struct {
	s *spinner.Spinner
	N int
	start_time time.Time
	interval_time time.Time
	update_interval int
}

func (spin *Spinner) Update(n int) {
	spin.N += n
	if spin.N % spin.update_interval == 0 {
		t_rate := float64(spin.N) / time.Now().Sub(spin.start_time).Seconds()
		rate   := float64(spin.update_interval) / time.Now().Sub(spin.interval_time).Seconds()
		spin.s.Suffix = fmt.Sprintf(" %12d complete (%.3f items/sec [total], %.3f items/sec [current interval])", spin.N, t_rate, rate)
		spin.interval_time = time.Now()
	}
}

func (spin *Spinner) Stop() {
	total := time.Now().Sub(spin.start_time)
	rate  := float64(spin.N) / time.Now().Sub(spin.start_time).Seconds()
	spin.s.FinalMSG = prefix(complete_color.Style("[✔]")) + 
		fmt.Sprintf("%d %s (total time: %s, %.3f items/sec [total])\n", spin.N, complete_color.Style("[complete]"), total.String(), rate)
	spin.s.Stop()
}

func StartSpinner(spin []string, update_interval int) Spinner {
	s := spinner.New(spin, 100 * time.Millisecond)
	s.Color("yellow")
	s.Prefix = spin_color.Style(prefix("[RUNNING]"))
	s.Start()

	return Spinner{s, 0, time.Now(), time.Now(), update_interval}
}

// ----------------------------------------------------------------------------------------------------------------
// Gauge
// ----------------------------------------------------------------------------------------------------------------
type Gauge struct {
	N int
	total int
	start_time time.Time
	interval_time time.Time
	update_interval int
}

func (g *Gauge) bar(width int) string {
	frac := float64(g.N) / float64(g.total)
	n_c  := int(frac * float64(width) + 0.5)
	r    := width - n_c
	
	return "[" + strings.Repeat("#", n_c) + strings.Repeat(" ", r) + "]"
}

func (g *Gauge) Update(n int) {
	g.N += n
	if g.N % g.update_interval == 0 {
		t_rate := float64(g.N) / time.Now().Sub(g.start_time).Seconds()
		rate   := float64(g.update_interval) / time.Now().Sub(g.interval_time).Seconds()
		str    := spin_color.Style(prefix("[RUNNING]") + g.bar(20)) + fmt.Sprintf(" (%.3f items/sec [total], %.3f items/sec [current interval])", t_rate, rate)
		fmt.Fprintln(os.Stdout, clear + str)
		g.interval_time = time.Now()
	}
}

func (g *Gauge) Stop() {
	total := time.Now().Sub(g.start_time)
	rate  := float64(g.N) / time.Now().Sub(g.start_time).Seconds()
	str   := prefix(complete_color.Style("[✔]")) +
		fmt.Sprintf("%d %s (total time: %s, %.3f items/sec [total])", g.N, complete_color.Style("[complete]"), total.String(), rate)
	fmt.Fprintln(os.Stdout, clear + str)
}

func StartGauge(total int, update_interval int) Gauge {
	fmt.Fprint(os.Stdout, spin_color.Style(prefix("[RUNNING]") + "\n"))
	return Gauge{0, total, time.Now(), time.Now(), update_interval}
}

// ----------------------------------------------------------------------------------------------------------------
// Logging
// ----------------------------------------------------------------------------------------------------------------
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
