package main
import (
	"fmt"
	"strings"
)
var verbose = false

func SetVerbose(b bool) {
	verbose=b
}

func VerbosePrintf(format string, a ...interface{}) (n int, err error) {
	if verbose {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

var verboseSpaceNum = 0
func VerbosePrintfStart(format string, a ...interface{}) (n int, err error) {
	n,err = VerbosePrintfIn(format, a...)
	verboseSpaceNum++
	return n, err
}
func VerbosePrintfIn(format string, a ...interface{}) (n int, err error) {
	n,err = VerbosePrintf(strings.Repeat(" ", verboseSpaceNum) + format, a...)
	return n, err
}
func VerbosePrintfEnd(format string, a ...interface{}) (n int, err error) {
	verboseSpaceNum--
	n,err = VerbosePrintfIn("/" + format, a...)
	return n, err
}
