package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

type payload struct {
	User    string          `json:"username"`
	Pass    string          `json:"password"`
	Options map[string]bool `json:"options"`
}

type okta struct {
	User        string
	Pass        string
	AuthAPI     string
	IDPLoginURL string
	C           *http.Client
}

func NewOkta(user string, pass string, api string, idpLoginURL string, c *http.Client) *okta {
	return &okta{
		User:        user,
		Pass:        pass,
		AuthAPI:     api,
		IDPLoginURL: idpLoginURL,
		C:           c,
	}
}

func (o *okta) GetSessionToken() (string, error) {
	p := &payload{
		User: o.User,
		Pass: o.Pass,
		Options: map[string]bool{
			"multiOptionalFactorEnroll": false,
			"warnBeforePasswordExpired": false,
		},
	}
	pJSON, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	r, err := o.C.Post(o.AuthAPI, "application/json", bytes.NewBuffer(pJSON))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != 200 {
		err = fmt.Errorf("Error: %s", body)
		return "", err
	}

	status := gjson.GetBytes(body, "status").String()
	sessionToken := ""
	if status == "SUCCESS" {
		sessionToken = gjson.GetBytes(body, "sessionToken").String()
	} else {
		err = fmt.Errorf("Error: %s", status)
		return "", err
	}

	return sessionToken, nil
}

func (o *okta) GetSAMLAssertion(sessionToken string) (string, error) {
	formData := url.Values{
		"onetimetoken": {sessionToken},
	}
	formDataBytes := bytes.NewBufferString(formData.Encode())
	r, err := o.C.Post(o.IDPLoginURL, "application/x-www-form-urlencoded", formDataBytes)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	if r.StatusCode != 200 {
		err = fmt.Errorf("Error: %s", body)
		return "", err
	}
	// can not generate doc directly from r.Body
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	samlAssertion, exits := doc.Find("input[name=SAMLResponse]").Attr("value")
	if !exits {
		err = fmt.Errorf("Error: SAMLAssertion token is not found in the response")
		return "", err
	}

	return samlAssertion, nil
}
