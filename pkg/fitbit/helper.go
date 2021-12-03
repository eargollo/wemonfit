package fitbit

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
)

func tokenFromFile(path string) (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t oauth2.Token
	if err := json.Unmarshal(content, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
