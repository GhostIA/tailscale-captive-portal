package hcloud

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/GhostIA/tailscale-captive-portal/internal/pkg/logger"
	"github.com/GhostIA/tailscale-captive-portal/internal/pkg/utils"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"golang.org/x/oauth2"
)

type HeznerServerClient struct {
	token  string
	ctx    context.Context
	client *hcloud.Client
	log    *slog.Logger
}

func NewHCloudServerClient(ctx context.Context) (*HeznerServerClient, error) {

	token := os.Getenv("HCLOUD_TOKEN")
	if token == "" {
		err := fmt.Errorf("failed to get token %s", token)
		return nil, err
	}
	logger := logger.NewLogger()
	slog.SetDefault(logger.ReturnDefault())

	ctx = context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logger
	tc := oauth2.NewClient(ctx, ts)
	retryClient.HTTPClient = tc
	retryClient.Backoff = utils.NewExponentialBackoff()
	slog.Debug("Calling to server client...")
	req, err := retryablehttp.NewRequest("GET", "http://api.hetzner.cloud/v1/servers", nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err

	}
	resp, err := retryClient.Do(req)
	if err != nil {
		slog.Error("request failed after all retries", "error", err)
		return nil, err
	} else {
		slog.Info("client connected", "status", resp.Status)
		resp.Body.Close()

	}
	client := hcloud.NewClient(hcloud.WithToken(token), hcloud.WithHTTPClient(retryClient.StandardClient()))
	return &HeznerServerClient{
		token:  token,
		ctx:    ctx,
		client: client,
		log:    logger.ReturnDefault(),
	}, nil

}
func (hsc *HeznerServerClient) GetAllRunningServers() ([]*hcloud.Server, error) {
	opts := hcloud.ServerListOpts{
		Status: []hcloud.ServerStatus{"running"},
	}
	return hsc.client.Server.AllWithOpts(hsc.ctx, opts)

}
