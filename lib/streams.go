package lib

import (
	"fmt"
	"os"
	"strings"
)

// StdOut is the standard output. Use this for informational output
var StdOut = os.Stdout

// StdErr is the standard error. Use this for error output
var StdErr = os.Stderr

// StdIn is the standard input. Use this for user input
var StdIn = os.Stdin

// PrintError prints a message to the standard output. The arguments work like Printf
func PrintError(template string, params ...any) {
	_, _ = StdErr.WriteString(fmt.Sprintf("ERR: "+template+"\n", params...))
}

// prompt will ask the user for a boolean input. If nonInteractive is set to true, then the prompt won't be displayed
// and true will be returned. The expected input is either "y" or "n", regardless of the case. The default is "N",
// therefore false
func prompt(message string, nonInteractive bool) bool {
	if nonInteractive {
		return true
	}
	_, _ = StdOut.WriteString(message + "[y/N]: ")
	b := make([]byte, 1)
	_, _ = StdIn.Read(b)
	val := strings.ToLower(string(b))
	if len(val) == 0 {
		return false
	}
	return val == "y"
}
