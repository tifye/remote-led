package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
)

type RGB struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

var (
	url string
)

func init() {
	url = os.Getenv("BACKEND")
}

func NewRGB(red, green, blue int) RGB {
	return RGB{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

func Fill(ctx context.Context, rgb RGB) error {
	url := fmt.Sprintf("%s/fill", url)
	b, err := json.Marshal(rgb)
	if err != nil {
		log.Errorf("failed to marshal request: %v", err)
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Errorf("failed to post request: %v", err)
		return err
	}

	return nil
}

func IsOn(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/isOn", url)
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to get request: %v", err)
		return false, err
	}
	defer res.Body.Close()

	isOn := false
	err = json.NewDecoder(res.Body).Decode(&isOn)
	return isOn, err
}
