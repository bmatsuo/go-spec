package main
/*
 *  Filename:    search.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 18:13:46 PDT 2011
 *  Description: <no value>
 */
import (
    "path/filepath"
    "strings"
    "os"
)

const SpecSuffix = "_spec.go"

type SpecCollector []string

func (sc *SpecCollector) VisitDir(path string, f *os.FileInfo) bool { return true }
func (sc *SpecCollector) VisitFile(path string, f *os.FileInfo) {
    if strings.HasSuffix(f.Name, SpecSuffix) {
        *sc = append(*sc, path)
    }
}

//  Finds all files <root>/**/*.spec
func SpecGoFiles(root string) (files []string, err os.Error) {
    if root, err = filepath.Abs(root); err != nil {
        return
    }
    var sc SpecCollector
    filepath.Walk(root, &sc, nil)
    files = sc
    return
}
