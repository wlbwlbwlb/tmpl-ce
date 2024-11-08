/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/wlbwlbwlb/tmpl/db"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmpl",
	Short: "A cli tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var yaml string

func init() {
	cobra.OnInitialize(readConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tmpl.yaml)")
	rootCmd.PersistentFlags().StringVar(&yaml, "yaml", filepath.Join(os.Getenv("TMPL"), "configs/config.yaml"), "config path")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Config struct {
	DSN     string
	Project string //如github.com/wlbwlbwlb/tmpl
	DbName  string
}

var config Config

func readConfig() {
	viper.SetConfigFile(yaml)

	e := viper.ReadInConfig()
	if e != nil {
		panic(e.Error())
	}
	viper.Unmarshal(&config) //将配置文件绑定到config上

	conf, e := mysql.ParseDSN(config.DSN)
	if e != nil {
		panic(e)
	}
	config.DbName = conf.DBName

	db.Init(config.DSN)
}
