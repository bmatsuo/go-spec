package main
/*
 *  Filename:    test.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 18:42:30 PDT 2011
 *  Description: <no value>
 */
import (
    "exec"
    "fmt"
    "os"
)

type GoTest []string

func GoTestCommand() GoTest { return nil }

func (cmd1 GoTest) SpecGoFiles(specfiles []string) (cmd2 GoTest) {
    cmd2 = cmd1
    for i := range specfiles {
        cmd2 = append(cmd2, "-file", specfiles[i])
    }
    return
}

func (cmd1 GoTest) Verbose() (cmd2 GoTest) {
    cmd2 = cmd1
    cmd2 = append(cmd2, "-v")
    return
}

func (cmd1 GoTest) TestPattern(pattern string) (cmd2 GoTest) {
    cmd2 = cmd1
    cmd2 = append(cmd2, fmt.Sprintf("-run=%s", pattern))
    return
}

func (cmd GoTest) Run(specpattern string) os.Error {
    excmd := exec.Command("gotest", cmd...)

    excmd.Env = os.Environ()
    excmd.Env = append(excmd.Env, fmt.Sprintf("GOSPECPATTERN=%s", specpattern))

    excmd.Stdout = os.Stdout
    excmd.Stderr = os.Stderr
    excmd.Stdin = os.Stdin
    return excmd.Run()
}
