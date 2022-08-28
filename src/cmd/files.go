/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"flag"

	"github.com/jotego/jtframe/jtfiles"

	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files <sim|syn> <core-name>",
	Short: "Generates the project compilation and simulation files",
	Long: `The project files are defined in cores/corename/game.yaml.
jtframe files command will also add the required files for the
selected compilation or simulation target.
The first argument selects simulation (sim) or synthesis (output). The
synthesis output consists of .qip files compatible with Intel Quartus.`,
	Run:  run_files,
	Args: cobra.ExactArgs(2),
}

var files_args jtfiles.Args

func init() {
	rootCmd.AddCommand(filesCmd)

	flag.StringVar(&files_args.Output, "output", "", "Output file name with no extension. Default is 'game'")
	flag.StringVar(&files_args.Target, "target", "", "Target platform: mist, mister, pocket, etc.")
	flag.BoolVar(&files_args.Rel, "rel", false, "Output relative paths")
	flag.BoolVar(&files_args.SkipVHDL, "novhdl", false, "Skip VHDL files")
}

func run_files(cmd *cobra.Command, args []string) {
	files_args.Corename = args[1]
	files_args.Format = args[0]

	jtfiles.Run(files_args)
}
