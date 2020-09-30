package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/urfave/cli/v2"
)

func init() {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"v"},
		Usage: "print dawn version",
	}

	cli.VersionPrinter = func(c *cli.Context) {
		var (
			cur, latest string
			err         error
		)

		if cur, err = currentVersion(); err != nil {
			cur = err.Error()
		}

		if latest, err = latestVersion(); err != nil {
			fmt.Printf("dawn version: %v\n", err)
			return
		}

		fmt.Printf("dawn version: %s(latest %s)\n", cur, latest)
	}
}

var currentVersionRegexp = regexp.MustCompile(`github\.com/kiyonlin/dawn*?\s+(.*)\n`)
var currentVersionFile = "go.mod"

func currentVersion() (string, error) {
	b, err := ioutil.ReadFile(currentVersionFile)
	if err != nil {
		return "", err
	}

	if submatch := currentVersionRegexp.FindSubmatch(b); len(submatch) == 2 {
		return string(submatch[1]), nil
	}

	return "", errors.New("github.com/kiyonlin/dawn was not found in go.mod")
}

var latestVersionRegexp = regexp.MustCompile(`"name":"(v.*?)"`)

func latestVersion() (v string, err error) {
	var (
		res *http.Response
		b   []byte
	)

	if res, err = http.Get("https://api.github.com/repos/kiyonlin/dawn/releases/latest"); err != nil {
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if b, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}

	if submatch := latestVersionRegexp.FindSubmatch(b); len(submatch) == 2 {
		return string(submatch[1]), nil
	}

	return "", errors.New("no version found in github response body")
}
