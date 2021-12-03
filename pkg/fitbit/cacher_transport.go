package fitbit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var errExpiredToken = errors.New("expired token")

type cacherTransport struct {
	Base *oauth2.Transport
	Path string
}

func (c *cacherTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	cachedToken, err := tokenFromFile(c.Path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if _, err := c.Base.Source.Token(); err != nil {
		return nil, errExpiredToken
	}
	resp, err = c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	newTok, err := c.Base.Source.Token()
	if err != nil {
		// While we’re unable to obtain a new token, the request was still
		// successful, so let’s gracefully handle this error by not caching a
		// new token. In either case, the user will need to re-authenticate.
		return resp, nil
	}
	if cachedToken == nil ||
		cachedToken.AccessToken != newTok.AccessToken ||
		cachedToken.RefreshToken != newTok.RefreshToken {
		bytes, err := json.Marshal(&newTok)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(c.Path, bytes, 0600); err != nil {
			return nil, err
		}
	}
	return resp, nil
}
