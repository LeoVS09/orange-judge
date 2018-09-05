package log

import (
	"fmt"
	"github.com/fatih/color"
	baselog "log"
	"os"
	"runtime"
)

var (
	cyan   = color.New(color.FgCyan).SprintfFunc()
	blue   = color.New(color.FgBlue).SprintfFunc()
	red    = color.New(color.FgRed).SprintfFunc()
	yellow = color.New(color.FgHiYellow).SprintfFunc()
	green  = color.New(color.FgGreen).SprintfFunc()
)

var standard = logger{
	warning: yellow,
	err:     red,
	log:     green,
	info:    cyan,
	debug:   blue,
	core:    baselog.New(os.Stdout, "", baselog.LstdFlags),
}

func init() {
	if runtime.GOOS == "windows" {
		SetNoColoring()
	}

	VerboseAll()
}

type logger struct {
	warning     func(format string, a ...interface{}) string
	err         func(format string, a ...interface{}) string
	log         func(format string, a ...interface{}) string
	info        func(format string, a ...interface{}) string
	debug       func(format string, a ...interface{}) string
	showError   bool
	showWarning bool
	showLog     bool
	showInfo    bool
	showDebug   bool
	core        *baselog.Logger
}

func Check(id string, err error) {
	if err != nil {
		baselog.Fatal(notColored(id), err) // TODO: fix coloring, change to red
	}
}

func Panic(id string, err error) {
	if err != nil {
		baselog.Println(notColored(id)) // TODO: fix coloring, change to red
		panic(err)
	}
}

func (l *logger) SetNoColoring() {
	l.warning = notColored
	l.err = notColored
	l.log = notColored
	l.info = notColored
	l.debug = notColored
}

// TODO: add levels of log
func (l *logger) VerboseAll() {
	l.showError = true
	l.showWarning = true
	l.showLog = true
	l.showInfo = true
	l.showDebug = true
}

func (l *logger) VerboseOnlyErrors() {
	l.showError = true
	l.showWarning = true
	l.showLog = false
	l.showInfo = false
	l.showDebug = false
}

func (l *logger) VerboseProduction() {
	l.showError = true
	l.showWarning = false
	l.showLog = true
	l.showInfo = false
	l.showDebug = false
}

type AsyncStdout struct {
	stack chan []byte
}

func (a *AsyncStdout) Write(p []byte) (n int, err error) {
	// TODO: add writing in file
	a.stack <- p

	return 0, nil
}

func createAsyncStdout() *AsyncStdout {
	// TODO: create close system
	out := make(chan []byte, 100)

	go func() {
		for p := range out {
			os.Stdout.Write(p)
		}
	}()

	return &AsyncStdout{
		stack: out,
	}
}

func (l *logger) UseAsync() {
	l.core.SetOutput(createAsyncStdout())
}

func (l *logger) UseSync() {
	l.core.SetOutput(os.Stdout)
}

func (l *logger) Printf(s string, a ...interface{}) {
	if l.showLog {
		l.core.Printf(s, a...)
	}
}

func (l *logger) Println(s string) {
	if l.showLog {
		l.core.Println(s)
	}
}

func (l *logger) Error(s string, v ...interface{}) {
	if l.showError {
		baselog.Printf(l.err("Error: "+s+"\n%s"), v...)
	}
}

func (l *logger) Warning(s string) {
	if l.showWarning {
		l.core.Println(l.warning(s))
	}
}

func (l *logger) WarningFmt(s string, v ...interface{}) {
	if l.showWarning {
		l.core.Println(l.warning(fmt.Sprintf(s, v...)))
	}
}

func (l *logger) Log(s string) {
	if l.showLog {
		l.core.Println(l.log(s))
	}
}

func (l *logger) LogFmt(s string, a ...interface{}) {
	if l.showLog {
		l.core.Println(l.log(fmt.Sprintf(s, a...)))
	}
}

func (l *logger) LogWhite(s string) {
	if l.showLog {
		l.core.Println(s)
	}
}

func (l *logger) Info(s string) {
	if l.showInfo {
		l.core.Println(l.info(s))
	}
}

func (l *logger) InfoFmt(s string, a ...interface{}) {
	if l.showInfo {
		l.core.Println(l.info(fmt.Sprintf(s, a...)))
	}
}

func (l *logger) Debug(s string) {
	if l.showDebug {
		l.core.Println(l.debug(s))
	}
}

func (l *logger) DebugFmt(s string, a ...interface{}) {
	if l.showDebug {
		l.core.Println(l.debug(fmt.Sprintf(s, a...)))
	}
}

func notColored(format string, a ...interface{}) string {
	return format
}

// ------------------ STANDARD LOGGER ----------------------------

func VerboseAll() {
	standard.VerboseAll()
}

func SetNoColoring() {
	standard.SetNoColoring()
}

func VerboseOnlyErrors() {
	standard.VerboseOnlyErrors()
}

func VerboseProduction() {
	standard.VerboseProduction()
}

// This can work unstable, use on you at you own risk
// Some logs can be logged not on real order
func UseAsync() {
	fmt.Println(standard.warning("------------Logger------------"))
	fmt.Println(standard.warning("Now using async logging."))
	fmt.Println(standard.warning("This in beta! Use it at you own risk"))
	fmt.Println(standard.warning("Some logs can be logged not in real order"))
	standard.UseAsync()
}

func UseSync() {
	standard.UseSync()
}

func Printf(s string, a ...interface{}) {
	standard.Printf(s, a...)
}

func Println(s string) {
	standard.Println(s)
}

func Error(s string, v ...interface{}) {
	standard.Error(s, v...)
}

func Warning(s string) {
	standard.Warning(s)
}

func WarningFmt(s string, v ...interface{}) {
	standard.WarningFmt(s, v...)
}

func Log(s string) {
	standard.Log(s)
}

func LogFmt(s string, a ...interface{}) {
	standard.LogFmt(s, a...)
}

func LogWhite(s string) {
	standard.LogWhite(s)
}

func Info(s string) {
	standard.Info(s)
}

func InfoFmt(s string, a ...interface{}) {
	standard.InfoFmt(s, a...)
}

func Debug(s string) {
	standard.Debug(s)
}

func DebugFmt(s string, a ...interface{}) {
	standard.DebugFmt(s, a...)
}
