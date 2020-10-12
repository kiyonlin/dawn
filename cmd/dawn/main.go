package main

import (
	"fmt"
	"os"

	"github.com/kiyonlin/dawn/cmd/dawn/internal"
	"github.com/spf13/cobra"
)

const version = "v0.0.7"

func init() {
	rootCmd.AddCommand(
		internal.VersionCmd, internal.NewCmd, internal.GenerateCmd, internal.DevCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "dawn",
	Short: "Dawn is an opinionated lightweight framework cli",
	Long:  longDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	if err = rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	return
}

const longDescription = `       __
   ___/ /__ __    _____    dawn-cli ` + version + `
 ~/ _  / _ '/ |/|/ / _ \~  For the opinionated lightweight framework dawn
~~\_,_/\_,_/|__,__/_//_/~~ Visit https://github.com/kiyonlin/dawn for detail
 ~~~  ~~ ~~~~~~~~~ ~~~~~~  (c) since 2020 by kiyon@gofiber.io`
