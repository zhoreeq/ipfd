package ipfd

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"net/url"

	"github.com/joho/godotenv"
)

type Config struct {
	siteURL             string
	sitePath            string
	siteName            string
	bindAddress         string
	databaseURL         string
	templatesPath       string
	staticPath          string
	staticURL           string
	serveStatic         bool
	ipfsAPI             string
	ipfsGateway         string
	ipfsPin             bool
	maxFileSize         int64
	allowedContentTypes []string
	premoderation       bool
	enableComments      bool
	enableVotes         bool
}

func LoadConfig(configPath string) (*Config, error) {
	var err error
	config := &Config{}

	if err = godotenv.Load(configPath); err != nil {
		return config, errors.New("Error loading .env file")
	}

	config.siteURL = os.Getenv("SITE_URL")
	if string(config.siteURL[len(config.siteURL)-1]) == "/" {
		return nil, errors.New("invalid SITE_URL, remove ending slash")
	}

	config.sitePath = "/"
	var u *url.URL
	if u, err = url.Parse(config.siteURL); err != nil {
		return nil, err
	} else if len(u.Path) != 0 {
		config.sitePath = u.Path + "/"
	}

	config.siteName = os.Getenv("SITE_NAME")
	config.bindAddress = os.Getenv("BIND_ADDRESS")
	config.databaseURL = os.Getenv("DATABASE_URL")
	config.templatesPath = os.Getenv("TEMPLATES_PATH")
	config.staticPath = os.Getenv("STATIC_PATH")
	config.staticURL = os.Getenv("STATIC_URL")
	if config.serveStatic, err = strconv.ParseBool(os.Getenv("SERVE_STATIC")); err != nil {
		return nil, err
	}

	config.ipfsAPI = os.Getenv("IPFS_API")
	config.ipfsGateway = os.Getenv("IPFS_GATEWAY")
	if config.ipfsPin, err = strconv.ParseBool(os.Getenv("IPFS_PIN")); err != nil {
		return nil, err
	}

	if config.maxFileSize, err = strconv.ParseInt(os.Getenv("MAX_FILESIZE"), 10, 64); err != nil {
		return nil, err
	}
	config.allowedContentTypes = strings.Split(os.Getenv("ALLOWED_CONTENT_TYPES"), ",")
	if config.premoderation, err = strconv.ParseBool(os.Getenv("PREMODERATION")); err != nil {
		return nil, err
	}
	if config.enableComments, err = strconv.ParseBool(os.Getenv("ENABLE_COMMENTS")); err != nil {
		return nil, err
	}
	if config.enableVotes, err = strconv.ParseBool(os.Getenv("ENABLE_VOTES")); err != nil {
		return nil, err
	}

	return config, nil
}
