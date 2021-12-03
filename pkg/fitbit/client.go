package fitbit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

var clientID = "229L7W"
var redirectURL = "https://localhost:7319"
var scopes = []string{"weight"}
var tokenFile = "fitibit-weight.token"

type Client struct {
	id          string
	secret      string
	redirectURL string
	token       *oauth2.Token
	client      *http.Client
}

func New(secret string) (*Client, error) {
	cli := &Client{id: clientID, secret: secret, redirectURL: redirectURL}
	err := cli.init()

	return cli, err
}

func (cli *Client) init() (err error) {
	conf := &oauth2.Config{
		ClientID:     cli.id,
		ClientSecret: cli.secret,
		Scopes:       scopes,
		Endpoint:     fitbit.Endpoint,
		RedirectURL:  cli.redirectURL,
	}

	err = cli.retrieveToken(conf)

	if err != nil {
		return fmt.Errorf("could not authenticate and get token. Error: %v", err)
	}

	// Like oauth2.Config.Client(), but using cacherTransport to persist tokens
	cli.client = &http.Client{
		Transport: &cacherTransport{
			Path: tokenFile,
			Base: &oauth2.Transport{
				Source: conf.TokenSource(oauth2.NoContext, cli.token),
			},
		},
	}
	return nil
}

func (cli *Client) retrieveToken(conf *oauth2.Config) (err error) {
	cli.token, err = tokenFromFile(tokenFile)
	if err != nil && os.IsNotExist(err) {
		cli.token, err = cli.authenticate(conf)
	}

	return
}

func (cli *Client) authenticate(conf *oauth2.Config) (*oauth2.Token, error) {
	tokens := make(chan *oauth2.Token)
	errors := make(chan error)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, "Missing 'code' parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Got code = [%s]", code)
		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			errors <- fmt.Errorf("could not exchange auth code for a token: %v", err)
			return
		}
		tokens <- tok
	})
	go func() {
		// Unfortunately, we need to hard-code this port — when registering
		// with fitbit, full RedirectURLs need to be whitelisted (incl. port).
		errors <- http.ListenAndServeTLS(":7319", "localhost.crt", "localhost.key", nil)
	}()

	authUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Println("Please visit the following URL to authorize:")
	fmt.Println(authUrl)
	select {
	case err := <-errors:
		return nil, err
	case token := <-tokens:
		return token, nil
	}
}
