/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/wlbwlbwlb/tmpl/util"

	"github.com/spf13/cobra"
)

// layoutCmd represents the layout command
var layoutCmd = &cobra.Command{
	Use:   "layout",
	Short: "Create a new layout",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("layout called")

		tmplText, e := os.ReadFile(filepath.Join(os.Getenv("TMPL"), "templates/layout.tmpl"))
		if e != nil {
			log.Fatal(e)
		}

		tmpl := template.Must(template.New("layout").Funcs(util.FuncMap).Parse(string(tmplText)))

		if e = tmpl.Execute(io.Discard, map[string]interface{}{
			"moduleName": moduleName,
		}); e != nil {
			log.Fatal(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(layoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// layoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// layoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	layoutCmd.Flags().StringVarP(&moduleName, "moduleName", "m", "", "")
	layoutCmd.MarkFlagRequired("moduleName")
}
