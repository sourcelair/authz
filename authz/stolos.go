package authz

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/authorization"
	"github.com/twistlock/authz/core"
)

const (
	// StolosUrl indicates which URL to use as base when connecting to the API
	StolosUrl = "https://api.stolos.io"

	// StolosToken indicates the token to use when connecting to the API
	StolosToken = ""
)

type stolosAuthorizer struct {
	settings *StolosAuthorizerSettings
}

// BasicAuthorizerSettings provides settings for the basic authoerizer flow
type StolosAuthorizerSettings struct {
	StolosUrl   string // StolosUrl indicates which URL to use as base when connecting to the API
	StolosToken string // StolosToken indicates the token to use when connecting to the API
}

// NewBasicAuthZAuthorizer creates a new basic authorizer
func NewStolosAuthZAuthorizer(settings *StolosAuthorizerSettings) core.Authorizer {
	return &stolosAuthorizer{settings: settings}
}

func (f *stolosAuthorizer) Init() error {
	return nil
}

func (f *stolosAuthorizer) AuthZReq(authZReq *authorization.Request) *authorization.Response {
	logrus.Debugf("Received AuthZ request, method: '%s', url: '%s'", authZReq.RequestMethod, authZReq.RequestURI)
	if authZReq.User == "" || authZReq.User == "client" {
		return &authorization.Response{Allow: true}
	}
	url := fmt.Sprintf("%s/api/a0.1/certs/%s/", f.settings.StolosUrl, authZReq.User)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", f.settings.StolosToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Request failed: %q", err.Error())
		return &authorization.Response{Allow: false}
	}
	defer resp.Body.Close()

	logrus.Debugf("Request success, status code: %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("Response failed with status code: %d, body: %q", resp.StatusCode, body)
		return &authorization.Response{Allow: false}
	}
	return &authorization.Response{Allow: true}
}

// AuthZRes always allow responses from server
func (f *stolosAuthorizer) AuthZRes(authZReq *authorization.Request) *authorization.Response {
	return &authorization.Response{Allow: true}
}
