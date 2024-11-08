/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/wlbwlbwlb/tmpl/util"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// crudapiCmd represents the crudapi command
var crudapiCmd = &cobra.Command{
	Use:   "crudapi",
	Short: "Generate crud http api",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("crudapi called")

		tmplText, e := os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/crudapi.tmpl"))
		if e != nil {
			log.Fatal(e)
		}
		buf := bytes.Buffer{}

		tmpl := template.Must(template.New("create").Funcs(util.FuncMap).Parse(string(tmplText)))

		//1，structName 区分大小写
		if e = tmpl.Execute(&buf, map[string]interface{}{
			"projectName": config.Project,
			"moduleName":  moduleName,
			"structName":  structName,
			"pluralStyle": pluralStyle,
		}); e != nil {
			log.Fatal(e)
		}

		dst := fmt.Sprintf("./%s/%sControl/%s.go", moduleName, moduleName, strcase.ToKebab(structName))

		if e = os.WriteFile(dst, buf.Bytes(), os.ModePerm); e != nil {
			log.Fatal(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(crudapiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crudapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crudapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	crudapiCmd.Flags().StringVarP(&moduleName, "moduleName", "m", "", "")
	crudapiCmd.MarkFlagRequired("moduleName")

	crudapiCmd.Flags().StringVarP(&structName, "structName", "s", "", "")
	crudapiCmd.MarkFlagRequired("structName")

	crudapiCmd.Flags().StringVarP(&pluralStyle, "pluralStyle", "p", "", "plural style of struct")
	crudapiCmd.MarkFlagRequired("pluralStyle")
}
