package ipfd

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SiteURL             string
	SitePath            string
	SiteName            string
	BindAddress         string
	DatabaseURL         string
	TemplatesPath       string
	StaticPath          string
	StaticURL           string
	ServeStatic         bool
	IpfsAPI             []string
	IpfsGateway         string
	IpfsPin             bool
	MaxFileSize         int64
	AllowedContentTypes []string
	Premoderation       bool
	EnableComments      bool
	EnableVotes         bool
}

func LoadConfig(configPath string) (*Config, error) {
	var err error
	config := &Config{}

	if err = godotenv.Load(configPath); err != nil {
		return config, errors.New("Error loading .env file")
	}

	config.SiteURL = os.Getenv("SITE_URL")
	if string(config.SiteURL[len(config.SiteURL)-1]) == "/" {
		return nil, errors.New("invalid SITE_URL, remove ending slash")
	}

	config.SitePath = "/"
	var u *url.URL
	if u, err = url.Parse(config.SiteURL); err != nil {
		return nil, err
	} else if len(u.Path) != 0 {
		config.SitePath = u.Path + "/"
	}

	config.SiteName = os.Getenv("SITE_NAME")
	config.BindAddress = os.Getenv("BIND_ADDRESS")
	config.DatabaseURL = os.Getenv("DATABASE_URL")
	config.TemplatesPath = os.Getenv("TEMPLATES_PATH")
	config.StaticPath = os.Getenv("STATIC_PATH")
	config.StaticURL = os.Getenv("STATIC_URL")
	if config.ServeStatic, err = strconv.ParseBool(os.Getenv("SERVE_STATIC")); err != nil {
		return nil, err
	}

	config.IpfsAPI = strings.Split(os.Getenv("IPFS_API"), ",")
	config.IpfsGateway = os.Getenv("IPFS_GATEWAY")
	if config.IpfsPin, err = strconv.ParseBool(os.Getenv("IPFS_PIN")); err != nil {
		return nil, err
	}

	if config.MaxFileSize, err = strconv.ParseInt(os.Getenv("MAX_FILESIZE"), 10, 64); err != nil {
		return nil, err
	}
	config.AllowedContentTypes = strings.Split(os.Getenv("ALLOWED_CONTENT_TYPES"), ",")
	if config.Premoderation, err = strconv.ParseBool(os.Getenv("PREMODERATION")); err != nil {
		return nil, err
	}
	if config.EnableComments, err = strconv.ParseBool(os.Getenv("ENABLE_COMMENTS")); err != nil {
		return nil, err
	}
	if config.EnableVotes, err = strconv.ParseBool(os.Getenv("ENABLE_VOTES")); err != nil {
		return nil, err
	}

	return config, nil
}
