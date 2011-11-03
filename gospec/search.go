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

func (sc *SpecCollector) Walk(path string, f *os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if f.IsDirectory() {
        return nil
    }
	if strings.HasSuffix(f.Name, SpecSuffix) {
		*sc = append(*sc, path)
	}
    return nil
}

func (sc *SpecCollector) WalkFunc() (fn filepath.WalkFunc) {
    fn = func(path string, f *os.FileInfo, err error) error {
        return sc.Walk(path, f, err)
    }
    return
}

//  Finds all files <root>/**/*.spec
func SpecGoFiles(root string) (files []string, err error) {
	if root, err = filepath.Abs(root); err != nil {
		return
	}
	var sc SpecCollector
	err = filepath.Walk(root, sc.WalkFunc())
	files = sc
	return
}
