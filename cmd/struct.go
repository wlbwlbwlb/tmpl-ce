/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wlbwlbwlb/tmpl/util"

	"github.com/iancoleman/strcase"
	"github.com/smallnest/gen/dbmeta"
	"github.com/spf13/cobra"
)

// structCmd represents the struct command
var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Structure model from a biz table",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("struct called")

		tmplText, e := os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/column.tmpl"))
		if e != nil {
			log.Fatal(e)
		}

		db, e := sql.Open("mysql", config.DSN)
		if e != nil {
			log.Fatal(fmt.Sprintf("Error in open database: %v\n\n", e.Error()))
		}
		e = db.Ping()
		if e != nil {
			log.Fatal(fmt.Sprintf("Error pinging database: %v\n\n", e.Error()))
		}

		cols := make([]string, 0)

		meta, e := dbmeta.LoadMeta("mysql", db, config.DbName, tableName)
		if e != nil {
			log.Fatal(fmt.Sprintf("Warning - LoadMeta skipping table info for %s error: %v\n", tableName, e))
		}

		for _, col := range meta.Columns() {
			fmt.Printf("%s\n", col.String())

			var buf bytes.Buffer

			sqlType := strings.ToLower(col.DatabaseTypeName())

			goType, _ := typeMapping[sqlType]

			//代码或表维护一份default标签即可，防止有变更，两边都要改
			gormTag := fmt.Sprintf("column:%s;comment:%s;default:", col.Name(), col.Comment())

			tmpl := template.Must(template.New("struct").Funcs(util.FuncMap).Parse(string(tmplText)))

			if e = tmpl.Execute(&buf, map[string]interface{}{
				"fieldName": strcase.ToCamel(col.Name()),
				"fieldType": goType,
				"gormTag":   gormTag,
				"jsonTag":   col.Name(),
				"formTag":   col.Name(),
			}); e != nil {
				log.Fatal(e)
			}
			cols = append(cols, buf.String())
		}

		buf := bytes.Buffer{}

		//if "v2" == ver {
		//	tmplText, e = os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/struct.v2.tmpl"))
		//}
		//if "v1" == ver {
		//	tmplText, e = os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/struct.tmpl"))
		//}
		tmplText, e = os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/struct.tmpl"))
		if e != nil {
			log.Fatal(e)
		}

		tmpl := template.Must(template.New("struct").Funcs(util.FuncMap).Parse(string(tmplText)))

		//1，moduleName lowerCamel风格
		//2，structName lowerCamel风格
		if e = tmpl.Execute(&buf, map[string]interface{}{
			"projectName":          config.Project,
			"moduleName":           moduleName,
			"structName":           structName,
			"structNameLowerCamel": strcase.ToLowerCamel(structName),
			"dbName":               meta.SQLDatabase(),
			"tableName":            meta.TableName(),
			"fields":               strings.Join(cols, "\n"),
		}); e != nil {
			log.Fatal(e)
		}

		dst := fmt.Sprintf("./%s/%sModel/%s.go", moduleName, moduleName, strcase.ToKebab(structName))

		if e = os.WriteFile(dst, buf.Bytes(), os.ModePerm); e != nil {
			log.Fatal(e)
		}

		//格式化
		gofmt := exec.Command("gofmt", "-w", dst)
		if e = gofmt.Run(); e != nil {
			log.Fatal(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(structCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// structCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// structCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	structCmd.Flags().StringVarP(&moduleName, "moduleName", "m", "", "")
	structCmd.MarkFlagRequired("moduleName")

	structCmd.Flags().StringVarP(&tableName, "tableName", "t", "", "")
	structCmd.MarkFlagRequired("tableName")

	structCmd.Flags().StringVarP(&structName, "structName", "s", "", "")
	structCmd.MarkFlagRequired("structName")

	//structCmd.Flags().StringVar(&ver, "ver", "v1", "")
}

var typeMapping = map[string]string{
	"bigint":       "int64",
	"int":          "int",
	"tinyint":      "int",
	"smallint":     "int",
	"unsigned int": "uint",
	"decimal":      "float64",

	"varchar": "string",
	"char":    "string",
	"text":    "string",

	"timestamp": "time.Time",
	"datetime":  "time.Time",
	"date":      "time.Time",

	"json": "[]byte",
}
