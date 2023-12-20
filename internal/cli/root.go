package cli

import (
	"github.com/spf13/cobra"
)

var ConfigDir string
var SrcDir string
var DstDir string

var rootCmd = &cobra.Command{
	Use:   "cisco-config-mapper",
	Short: "cisco-config-mapper is a CLI tool to convert production to lab configs",
	Long:  "cisco-config-mapper is a CLI tool to convert production to lab configs",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigDir, "config", "c", "./config/", "Directory with mapping files")
	rootCmd.PersistentFlags().StringVarP(&SrcDir, "src", "s", "./src/", "Directory with production config files")
	rootCmd.PersistentFlags().StringVarP(&DstDir, "dst", "d", "./dst/", "Directory with lab config files")
}
