package util

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

var FuncMap = template.FuncMap{
	"mkdir": Mkdir,
	"touch": Touch,
}

func Mkdir(dstDir string) string {
	e := os.MkdirAll(dstDir, os.ModePerm)
	if e != nil {
		return fmt.Sprintf("mkdir returned an error %v", e)
	}
	return fmt.Sprintf("mkdir %s", dstDir)
}

// Touch template command to touch a file under the output directory
func Touch(dstDir string) string {
	_, e := os.Stat(dstDir)
	if os.IsNotExist(e) {
		file, e2 := os.Create(dstDir)
		if e2 != nil {
			return fmt.Sprintf("touch returned an error %v", e2)
		}
		defer file.Close()
	} else {
		t := time.Now().Local()
		e = os.Chtimes(dstDir, t, t)
		if e != nil {
			return fmt.Sprintf("touch returned an error %v", e)
		}
	}
	return fmt.Sprintf("touch %s", dstDir)
}
