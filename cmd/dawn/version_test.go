package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_Version_Printer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, latestVersionUrl, httpmock.NewBytesResponder(200, fakeVersionResponse))

		cli.VersionPrinter(nil)
	})

	t.Run("latest err", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, latestVersionUrl, httpmock.NewBytesResponder(200, []byte("no version")))

		cli.VersionPrinter(nil)
	})
}

func Test_Version_Current(t *testing.T) {
	at := assert.New(t)

	t.Run("file not found", func(t *testing.T) {
		setupCurrentVersionFile()
		defer teardownCurrentVersionFile()

		_, err := currentVersion()
		at.NotNil(err)
	})

	t.Run("match version", func(t *testing.T) {
		content := `module github.com/kiyonlin/dawn-demo
go 1.14
require (
	github.com/gofiber/fiber/v2 v2.0.4
	github.com/kiyonlin/dawn v0.3.0
	github.com/jarcoal/httpmock v1.0.6
)`

		setupCurrentVersionFile(content)
		defer teardownCurrentVersionFile()

		v, err := currentVersion()
		at.Nil(err)
		at.Equal("v0.3.0", v)
	})

	t.Run("match master", func(t *testing.T) {
		content := `module github.com/kiyonlin/dawn-demo
go 1.14
require (
	github.com/gofiber/fiber/v2 v2.0.4
	github.com/kiyonlin/dawn v0.0.0-20200926082917-55763e7e6ee3
	github.com/jarcoal/httpmock v1.0.6
)`

		setupCurrentVersionFile(content)
		defer teardownCurrentVersionFile()

		v, err := currentVersion()
		at.Nil(err)
		at.Equal("v0.0.0-20200926082917-55763e7e6ee3", v)
	})

	t.Run("package not found", func(t *testing.T) {
		content := `module github.com/kiyonlin/dawn-demo
go 1.14
require (
	github.com/gofiber/fiber/v2 v2.0.4
)`

		setupCurrentVersionFile(content)
		defer teardownCurrentVersionFile()

		_, err := currentVersion()
		at.NotNil(err)
	})
}

func setupCurrentVersionFile(content ...string) {
	currentVersionFile = "current-version"
	if len(content) > 0 {
		_ = ioutil.WriteFile(currentVersionFile, []byte(content[0]), 0600)
	}
}

func teardownCurrentVersionFile() {
	_ = os.Remove(currentVersionFile)
}

func Test_Version_Latest(t *testing.T) {
	at := assert.New(t)
	t.Run("http get error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, latestVersionUrl, httpmock.NewErrorResponder(errors.New("network error")))

		_, err := latestVersion()
		at.NotNil(err)
	})

	t.Run("version matched", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, latestVersionUrl, httpmock.NewBytesResponder(200, fakeVersionResponse))

		v, err := latestVersion()
		at.Nil(err)
		at.Equal("v0.3.0", v)
	})

	t.Run("no version matched", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(http.MethodGet, latestVersionUrl, httpmock.NewBytesResponder(200, []byte("no version")))

		_, err := latestVersion()
		at.NotNil(err)
	})
}

var fakeVersionResponse = []byte(`
{"url":"https://api.github.com/repos/kiyonlin/dawn/releases/31943095","assets_url":"https://api.github.com/repos/kiyonlin/dawn/releases/31943095/assets","upload_url":"https://uploads.github.com/repos/kiyonlin/dawn/releases/31943095/assets{?name,label}","html_url":"https://github.com/kiyonlin/dawn/releases/tag/v0.3.0","id":31943095,"node_id":"MDc6UmVsZWFzZTMxOTQzMDk1","tag_name":"v0.3.0","target_commitish":"4f0d4d8254630b68e17a7b305cd9d72c601ac7e1","name":"v0.3.0","draft":false,"author":{"login":"kiyonlin","id":1214670,"node_id":"MDQ6VXNlcjEyMTQ2NzA=","avatar_url":"https://avatars1.githubusercontent.com/u/1214670?v=4","gravatar_id":"","url":"https://api.github.com/users/kiyonlin","html_url":"https://github.com/kiyonlin","followers_url":"https://api.github.com/users/kiyonlin/followers","following_url":"https://api.github.com/users/kiyonlin/following{/other_user}","gists_url":"https://api.github.com/users/kiyonlin/gists{/gist_id}","starred_url":"https://api.github.com/users/kiyonlin/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/kiyonlin/subscriptions","organizations_url":"https://api.github.com/users/kiyonlin/orgs","repos_url":"https://api.github.com/users/kiyonlin/repos","events_url":"https://api.github.com/users/kiyonlin/events{/privacy}","received_events_url":"https://api.github.com/users/kiyonlin/received_events","type":"User","site_admin":false},"prerelease":false,"created_at":"2020-09-26T04:43:02Z","published_at":"2020-09-29T15:48:14Z","assets":[],"tarball_url":"https://api.github.com/repos/kiyonlin/dawn/tarball/v0.3.0","zipball_url":"https://api.github.com/repos/kiyonlin/dawn/zipball/v0.3.0"}
`)
