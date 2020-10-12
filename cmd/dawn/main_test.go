package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

func Test_Main_Run(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	t.Run("success", func(t *testing.T) {
		assert.Nil(t, run())
	})

	t.Run("error", func(t *testing.T) {
		rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
			return errors.New("")
		}
		assert.NotNil(t, run())
	})
}
