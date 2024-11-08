package util

import (
	"testing"

	"github.com/iancoleman/strcase"
)

func TestTouch(t *testing.T) {
	a := "TmplTmpl"
	t.Log(strcase.ToSnake(a))      //tmpl_tmpl
	t.Log(strcase.ToCamel(a))      //TmplTmpl
	t.Log(strcase.ToLowerCamel(a)) //tmplTmpl
	t.Log(strcase.ToKebab(a))      //tmpl-tmpl
}
