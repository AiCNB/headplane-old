package tsnet

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tale/headplane/agent/internal/config"
	"github.com/tale/headplane/agent/internal/i18n"
	"github.com/tale/headplane/agent/internal/util"
	"tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

// Wrapper type so we can add methods to the server.
type TSAgent struct {
	*tsnet.Server
	Lc *tailscale.LocalClient
	ID string
}

// Creates a new tsnet agent and returns an instance of the server.
func NewAgent(cfg *config.Config) *TSAgent {
	log := util.GetLogger()

	dir, err := filepath.Abs(cfg.WorkDir)
	if err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.tsnet.abs_path_failed",
			"Failed to get absolute path: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.tsnet.workdir_create_failed",
			"Cannot create agent work directory: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	server := &tsnet.Server{
		Dir:        dir,
		Hostname:   cfg.Hostname,
		ControlURL: cfg.TSControlURL,
		AuthKey:    cfg.TSAuthKey,
		Logf:       func(string, ...any) {}, // Disabled by default
		UserLogf:   log.Info,
	}

	if cfg.Debug {
		server.Logf = log.Debug
	}

	return &TSAgent{server, nil, ""}
}

// Starts the tsnet agent and sets the node ID.
func (s *TSAgent) Connect() {
	log := util.GetLogger()

	// Waits until the agent is up and running.
	status, err := s.Up(context.Background())
	if err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.tsnet.tailnet_connect_failed",
			"Failed to connect to Tailnet: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	s.Lc, err = s.LocalClient()
	if err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.tsnet.local_client_failed",
			"Failed to initialize local Tailscale client: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	id, err := status.Self.PublicKey.MarshalText()
	if err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.tsnet.marshal_failed",
			"Failed to marshal public key: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	log.Info("%s", i18n.Message(
		"agent.tsnet.connected",
		"Connected to Tailnet (PublicKey: {{.PublicKey}})",
		map[string]any{"PublicKey": status.Self.PublicKey},
	))
	s.ID = string(id)
}

// Shuts down the tsnet agent.
func (s *TSAgent) Shutdown() {
	s.Close()
}
