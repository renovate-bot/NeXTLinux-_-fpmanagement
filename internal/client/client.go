package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/nextlinux/enterprise-client-go/pkg/external"
)

type Configuration struct {
	Hostname       string
	Username       string
	Password       string
	Scheme         string
	Insecure       bool
	TimeoutSeconds int
	NextlinuxAccount string
}

type EnterpriseClient struct {
	config Configuration
	Client *external.APIClient
}

func NewEnterpriseClient(cfg Configuration) (*EnterpriseClient, error) {

	if cfg.Scheme == "" {
		cfg.Scheme = "https"
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Insecure}, //nolint:gosec
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(cfg.TimeoutSeconds) * time.Second,
	}
	return &EnterpriseClient{
		config: cfg,
		Client: external.NewAPIClient(&external.Configuration{
			BasePath:   "v1/enterprise",
			Host:       cfg.Hostname,
			Scheme:     cfg.Scheme,
			HTTPClient: client,
			DefaultHeader: map[string]string{
				"x-nextlinux-account": cfg.NextlinuxAccount,
			},
		}),
	}, nil
}

func (c *EnterpriseClient) NewRequestContext(parentContext context.Context) context.Context {
	if parentContext == nil {
		parentContext = context.Background()
	}
	return context.WithValue(
		parentContext,
		external.ContextBasicAuth,
		external.BasicAuth{
			UserName: c.config.Username,
			Password: c.config.Password,
		},
	)
}
