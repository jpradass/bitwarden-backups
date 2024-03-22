package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	my_http "github.com/jpradass/bitwarden-backups/http"
)

type BitWardenAPI struct {
	client *my_http.HttpClient
	token  string
}

func New() *BitWardenAPI {
	return &BitWardenAPI{
		client: my_http.New(),
		token:  "",
	}
}

func (bwAPI *BitWardenAPI) auth() error {
	response, err := bwAPI.client.MakeRequest(
		"POST",
		"https://identity.bitwarden.com/connect/token",
		[]byte(fmt.Sprintf("grant_type=client_credentials&scope=api&client_id=%s&client_secret=%s", os.Getenv("BITWARDEN_CLIENT_ID"), os.Getenv("BITWARDEN_CLIENT_SECRET"))),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	bbody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var body map[string]any
	err = json.Unmarshal(bbody, &body)
	if err != nil {
		return err
	}

	bwAPI.token = body["access_token"].(string)
	return nil
}
