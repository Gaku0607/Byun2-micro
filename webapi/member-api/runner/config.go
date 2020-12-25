package runner

import (
	"os"
	"strconv"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/joho/godotenv"
)

var (
	name     string
	regaddrs string
	port     string
	ttl      int
	interval int
)

func Init() error {

	var err error

	if err = godotenv.Load("member_api.env"); err != nil {
		return err
	}

	if name = os.Getenv("api_name"); name == "" {
		return models.ErrAPINameIsEmpty
	}

	if regaddrs = os.Getenv("registry_address"); regaddrs == "" {
		return models.ErrRegAddressIsEmpty
	}

	if port = os.Getenv("api_address"); port == "" {
		return models.ErrAPIAddressIsEmpty
	}

	if val := os.Getenv("api_ttl"); val == "" {
		return models.ErrAPITTLIsEmpty
	} else {
		if ttl, err = strconv.Atoi(val); err != nil {
			return err
		}
	}

	if val := os.Getenv("api_interval"); val == "" {
		return models.ErrAPIIntervalIsEmpty
	} else {
		if interval, err = strconv.Atoi(val); err != nil {
			return err
		}
	}

	return nil
}
