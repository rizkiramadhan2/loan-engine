package errs

import (
	"fmt"
	"go/build"
	"os"
	"runtime"
	"strings"
)

var (
	fileDirWd = ""
	fileDir   = ""
	funcDir   = ""
)

func init() {
	wd, _ := os.Getwd()
	srcDirs := build.Default.SrcDirs()

	if len(srcDirs) != 2 {
		return
	}

	fileDirWd = wd
	fileDir = srcDirs[0]
	funcDir = strings.TrimPrefix(fileDirWd, srcDirs[1]+"/")
}

func makeTrace(skip int) []string {
	trace := []string{}
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, next := frames.Next()
		t := formatFrame(frame)
		trace = append(trace, t)
		if !next {
			break
		}
	}
	return trace
}

func formatFrame(frame runtime.Frame) string {
	file := strings.TrimPrefix(frame.File, fileDirWd)
	file = strings.TrimPrefix(file, fileDir)
	file = fmt.Sprintf("%s:%d", file, frame.Line)

	function := strings.TrimPrefix(frame.Function, funcDir)
	roundUp := 12
	roundUp = (len(function) + roundUp - 1) / roundUp * roundUp
	function = fmt.Sprintf("%-*s", roundUp, function)
	t := fmt.Sprintf("%s %s", function, file)
	return t
}
