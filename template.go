package suti

/*
Copyright (C) 2021 gearsix <gearsix@tuta.io>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"bytes"
	"fmt"
	mst "github.com/cbroglie/mustache"
	hmpl "html/template"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	tmpl "text/template"
)

// SupportedTemplateLangs provides a list of supported languages for template files (lower-case)
var SupportedTemplateLangs = []string{"tmpl", "hmpl", "mst"}

// IsSupportedTemplateLang provides the index of `SupportedTemplateLangs` that `lang` is at.
// If `lang` is not in `SupportedTemplateLangs`, `-1` will be returned.
// File extensions can be passed in `lang`, the prefixed `.` will be trimmed.
func IsSupportedTemplateLang(lang string) int {
	lang = strings.ToLower(lang)
	if lang[0] == '.' {
		lang = lang[1:]
	}
	for i, l := range SupportedTemplateLangs {
		if lang == l {
			return i
		}
	}
	return -1
}

func getTemplateType(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

// Template is a generic interface to any template parsed from LoadTemplateFile
type Template struct {
	Source   string
	Template interface{}
}

// Execute executes `t` against `d`. Reflection is used to determine
// the template type and call it's execution fuction.
func (t *Template) Execute(d interface{}) (result bytes.Buffer, err error) {
	var funcName string
	var params []reflect.Value
	switch (reflect.TypeOf(t.Template).String()) {
	case "*template.Template": // golang templates
		funcName = "Execute"
		params = []reflect.Value{reflect.ValueOf(&result), reflect.ValueOf(d)}
	case "*mustache.Template":
		funcName = "FRender"
		params = []reflect.Value{reflect.ValueOf(&result), reflect.ValueOf(d)}
	default:
		err = fmt.Errorf("unable to infer template type '%s'", reflect.TypeOf(t.Template).String())
	}

	if err == nil {
		rval := reflect.ValueOf(t.Template).MethodByName(funcName).Call(params)
		if !rval[0].IsNil() { // err != nil
			err = rval[0].Interface().(error)
		}
	}

	return
}

func loadTemplateFileTmpl(root string, partials ...string) (*tmpl.Template, error) {
	var stat os.FileInfo
	t, e := tmpl.ParseFiles(root)

	for i := 0; i < len(partials) && e == nil; i++ {
		p := partials[i]
		ptype := getTemplateType(p)

		stat, e = os.Stat(p)
		if e == nil {
			if ptype == "tmpl" || ptype == "gotmpl" {
				t, e = t.ParseFiles(p)
			} else if strings.Contains(p, "*") {
				t, e = t.ParseGlob(p)
			} else if stat.IsDir() {
				e = filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
					if err == nil && !info.IsDir() {
						ptype = getTemplateType(path)
						if ptype == "tmpl" || ptype == "gotmpl" {
							t, err = t.ParseFiles(path)
						}
					}
					return err
				})
			} else {
				return nil, fmt.Errorf("non-matching filetype")
			}
		}
	}

	return t, e
}

func loadTemplateFileHmpl(root string, partials ...string) (*hmpl.Template, error) {
	var stat os.FileInfo
	t, e := hmpl.ParseFiles(root)

	for i := 0; i < len(partials) && e == nil; i++ {
		p := partials[i]
		ptype := getTemplateType(p)

		stat, e = os.Stat(p)
		if e == nil {
			if ptype == "hmpl" || ptype == "gohmpl" {
				t, e = t.ParseFiles(p)
			} else if strings.Contains(p, "*") {
				t, e = t.ParseGlob(p)
			} else if stat.IsDir() {
				e = filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
					if err == nil && !info.IsDir() {
						ptype = getTemplateType(path)
						if ptype == "hmpl" || ptype == "gohmpl" {
							t, err = t.ParseFiles(path)
						}
					}
					return err
				})
			} else {
				return nil, fmt.Errorf("non-matching filetype")
			}
		}
	}

	return t, e
}

func loadTemplateFileMst(root string, partials ...string) (*mst.Template, error) {
	var err error
	for p, partial := range partials {
		if err != nil {
			break
		}

		if stat, e := os.Stat(partial); e != nil {
			partials = append(partials[:p], partials[p+1:]...)
			err = e
		} else if stat.IsDir() == false {
			partials[p] = filepath.Dir(partial)
		} else if strings.Contains(partial, "*") {
			if paths, e := filepath.Glob(partial); e != nil {
				partials = append(partials[:p], partials[p+1:]...)
				err = e
			} else {
				partials = append(partials[:p], partials[p+1:]...)
				partials = append(partials, paths...)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	mstfp := &mst.FileProvider{
		Paths:      partials,
		Extensions: []string{".mst", ".mustache"},
	}
	return mst.ParseFilePartials(root, mstfp)
}

// LoadTemplateFile loads a Template from file `root`. All files in `partials`
// that have the same template type (identified by file extension) are also
// parsed and associated with the parsed root template.
func LoadTemplateFile(root string, partials ...string) (t Template, e error) {
	if len(root) == 0 {
		e = fmt.Errorf("no root template specified")
	}

	if stat, err := os.Stat(root); err != nil {
		e = err
	} else if stat.IsDir() {
		e = fmt.Errorf("root path must be a file, not a directory: %s", root)
	}

	if e == nil {
		t = Template{Source: root}
		ttype := getTemplateType(root)
		if ttype == "tmpl" || ttype == "gotmpl" {
			t.Template, e = loadTemplateFileTmpl(root, partials...)
		} else if ttype == "hmpl" || ttype == "gohmpl" {
			t.Template, e = loadTemplateFileHmpl(root, partials...)
		} else if ttype == "mst" || ttype == "mustache" {
			t.Template, e = loadTemplateFileMst(root, partials...)
		} else {
			e = fmt.Errorf("'%s' is not a supported template language", ttype)
		}
	}

	return
}

