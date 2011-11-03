package main
/*
 *  Filename:    options.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 17:58:46 PDT 2011
 *  Description: Parse arguments and options from the command line.
 */
import (
    "strings"
    "flag"
    "fmt"
    "os"
)
/*
 *  Constants, variables, and functions that users may actually want to call
 *  are capitalized.
 */

var (
    // Set this variable to customize the help message header.
    // For example, `gospec [options] action [arg2 ...]`.
    CommandLineHelpUsage = `gospec [-v] [-test=PATTERN] [ROOT [PATTERN ...]]`
    // Set this variable to print a message after the option specifications.
    // For example, "For more help:\n\tgospec help [action]"
    CommandLineHelpFooter = `Spec files must end with a suffix "_spec.go".`
)

//  A struct that holds parsed option values.
type options struct {
    Root        string
    TestPattern string
    SpecPattern string
    Verbose     bool
}

//  Create a flag.FlagSet to parse the command line options/arguments.
func setupFlags(opt *options) *flag.FlagSet {
    fs := flag.NewFlagSet("gospec", flag.ExitOnError)
    fs.BoolVar(&(opt.Verbose), "v", false, "Verbose program output.")
    fs.StringVar(&(opt.Root), "root", "./spec", "Directory containing spec files.")
    fs.StringVar(&(opt.TestPattern), "test", ".*", "Regexp matching tests to run.")
    fs.StringVar(&(opt.SpecPattern), "spec", ".*", "Regexp matching tests to run.")
    setupUsage(fs)
    return fs
}

//  Check the options for acceptable values. Panics or otherwise exits
//  with a non-zero exitcode when errors are encountered.
func verifyFlags(opt *options, fs *flag.FlagSet) {
    args := fs.Args()
    if len(args) > 0 {
        opt.Root = args[0]
        args = args[1:]
    }
    if len(args) > 0 {
        patterns := make([]string, len(args))
        for i := range args {
            patterns[i] = fmt.Sprintf("(%s)", args[i])
        }
        opt.SpecPattern = strings.Join(patterns, "|")
    }
}

//  Print a help message to standard error. See constants CommandLineHelpUsage
//  and CommandLineHelpFooter.
func PrintHelp() {
    fs := setupFlags(&options{})
    fs.Usage()
}

//  Hook up the commandLineHelpUsage and commandLineHelpFooter strings
//  to the standard Go flag.Usage function.
func setupUsage(fs *flag.FlagSet) {
    printNonEmpty := func(s string) {
        if s != "" {
            fmt.Fprintf(os.Stderr, "%s\n", s)
        }
    }
    tmpUsage := fs.Usage
    fs.Usage = func() {
        printNonEmpty(CommandLineHelpUsage)
        tmpUsage()
        printNonEmpty(CommandLineHelpFooter)
    }
}

//  Parse the command line options, validate them, and process them
//  further (e.g. Initialize more complex structs) if need be.
func parseFlags() options {
    var opt options
    fs := setupFlags(&opt)
    fs.Parse(os.Args[1:])
    verifyFlags(&opt, fs)
    // Process the verified options...
    return opt
}
