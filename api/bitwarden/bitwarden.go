package bitwarden

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	my_http "github.com/jpradass/bitwarden-backups/http"
	"github.com/jpradass/bitwarden-backups/logging"
	"go.uber.org/zap"
)

type BitWardenAPI struct {
	authURL string
	token   string
}

var logger *zap.SugaredLogger

func init() {
	logger = logging.New()
}

// We need to check until when it's validated
func New() *BitWardenAPI {
	return &BitWardenAPI{
		authURL: "https://identity.bitwarden.com/connect/token",
		token:   "",
	}
}

func (bwAPI *BitWardenAPI) auth() error {
	logger.Debug("Obtaining authorization token...")

	response, err := my_http.MakeRequest(
		"POST",
		bwAPI.authURL,
		[]byte(fmt.Sprintf("deviceName=firefox&deviceIdentifier=0&deviceType=0&grant_type=client_credentials&scope=api&client_id=%s&client_secret=%s", os.Getenv("BITWARDEN_CLIENT_ID"), os.Getenv("BITWARDEN_CLIENT_SECRET"))),
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

func (bwAPI *BitWardenAPI) ListCollections() error {
	logger.Debug("Listing collections...")

	bwAPI.checkAuth()

	response, err := my_http.MakeRequest(
		"GET",
		"https://api.bitwarden.com/public/collections",
		nil,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", bwAPI.token),
			"Content-Type":  "application/json",
		},
	)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	bbody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(bbody))
	return nil
}

func (bwAPI *BitWardenAPI) checkAuth() {
	if bwAPI.token == "" {
		if err := bwAPI.auth(); err != nil {
			logger.Error("There was an error trying to authenticate on Bitwarden. Error: ", err.Error())
		}
	}
}
