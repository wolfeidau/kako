package kako

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/rpc/v2/json2"
)

type TimeseriesArgs struct {
	Events []*Event
}

type TimeseriesReply struct {
	Message string
}

type Config struct {
	Email      string
	Name       string
	URL        string
	SigningKey []byte
}

type Client struct {
	config *Config
	client *http.Client
}

func NewClient(config *Config) *Client {
	return &Client{config: config, client: &http.Client{}}
}

func (c *Client) SetHttpClient(client *http.Client) {
	c.client = client
}

func (c *Client) SaveEvents(events []*Event) (string, error) {

	args := &TimeseriesArgs{Events: events}

	var results TimeseriesReply

	err := c.rpcCall("Timeseries.SaveEvents", args, &results)

	log.Printf("reply %s", results.Message)

	return results.Message, err

}

func (c *Client) rpcCall(method string, args, reply interface{}) error {

	buf, _ := json2.EncodeClientRequest(method, args)

	body := bytes.NewBuffer(buf)
	token, err := buildClaim(c.config)

	if err != nil {
		return err
	}
	client := &http.Client{}

	req, err := buildClient(c.config, token, body)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	buf, err = ioutil.ReadAll(resp.Body)

	log.Printf("body %s", string(buf))

	body = bytes.NewBuffer(buf)

	return json2.DecodeClientResponse(body, reply)
}

func buildClient(config *Config, token string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest("POST", config.URL, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("jwt", token)

	return req, nil
}

func buildClaim(config *Config) (string, error) {

	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims["AccessToken"] = "level1"
	t.Claims["ServiceInfo"] = struct {
		Email string
		Name  string
	}{config.Email, config.Name}

	t.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := t.SignedString(config.SigningKey)

	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
