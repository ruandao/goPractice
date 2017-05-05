package tmpl

import (
	"io"
	"path/filepath"
	"os"
	"html/template"
	"strings"
)

var tmpl *template.Template

func AddDir(path string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".tmpl") {
			if tmpl == nil {
				tmpl,err = template.ParseFiles(path)
				if err != nil {
					panic(err)
				}
			} else {
				tmpl, err = tmpl.ParseFiles(path)
				if err != nil {
					panic(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Execute(wr io.Writer, blockName string, data interface{}) error {
	return tmpl.ExecuteTemplate(wr, blockName, data)
}