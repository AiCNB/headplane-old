package tsnet

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tale/headplane/agent/internal/i18n"
	"github.com/tale/headplane/agent/internal/util"
	"tailscale.com/tailcfg"
	"tailscale.com/types/key"

	"go4.org/mem"
)

// Returns the raw hostinfo for a peer based on node ID.
func (s *TSAgent) GetStatusForPeer(id string) (*tailcfg.HostinfoView, error) {
	log := util.GetLogger()

	if !strings.HasPrefix(id, "nodekey:") {
		log.Debug("%s", i18n.Message(
			"agent.tsnet.invalid_node_prefix",
			"Node ID with missing prefix: {{.NodeID}}",
			map[string]any{"NodeID": id},
		))
		return nil, fmt.Errorf("%s", i18n.Message(
			"agent.tsnet.error.invalid_node_id",
			"invalid node ID: {{.NodeID}}",
			map[string]any{"NodeID": id},
		))
	}

	log.Debug("%s", i18n.Message(
		"agent.tsnet.query_peer",
		"Querying status of peer: {{.NodeID}}",
		map[string]any{"NodeID": id},
	))
	status, err := s.Lc.Status(context.Background())
	if err != nil {
		log.Debug("%s", i18n.Message(
			"agent.tsnet.status_failed",
			"Failed to get status: {{.Error}}",
			map[string]any{"Error": err},
		))
		return nil, fmt.Errorf("%s", i18n.Message(
			"agent.tsnet.error.get_status",
			"failed to get status: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	// We need to convert from 64 char hex to 32 byte raw.
	bytes, err := hex.DecodeString(id[8:])
	if err != nil {
		log.Debug("%s", i18n.Message(
			"agent.tsnet.decode_failed",
			"Failed to decode hex: {{.Error}}",
			map[string]any{"Error": err},
		))
		return nil, fmt.Errorf("%s", i18n.Message(
			"agent.tsnet.error.decode_hex",
			"failed to decode hex: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	raw := mem.B(bytes)
	if raw.Len() != 32 {
		log.Debug("%s", i18n.Message(
			"agent.tsnet.invalid_length",
			"Invalid node ID length: {{.Length}}",
			map[string]any{"Length": raw.Len()},
		))
		return nil, fmt.Errorf("%s", i18n.Message(
			"agent.tsnet.error.invalid_length",
			"invalid node ID length: {{.Length}}",
			map[string]any{"Length": raw.Len()},
		))
	}

	nodeKey := key.NodePublicFromRaw32(raw)
	peer := status.Peer[nodeKey]
	if peer == nil {
		// Check if we are on Self.
		if status.Self.PublicKey == nodeKey {
			peer = status.Self
		} else {
			log.Debug("%s", i18n.Message(
				"agent.tsnet.peer_not_found",
				"Peer not found in status: {{.NodeID}}",
				map[string]any{"NodeID": id},
			))
			return nil, nil
		}
	}

	ip := peer.TailscaleIPs[0].String()
	whois, err := s.Lc.WhoIs(context.Background(), ip)
	if err != nil {
		log.Debug("%s", i18n.Message(
			"agent.tsnet.whois_failed",
			"Failed to get whois: {{.Error}}",
			map[string]any{"Error": err},
		))
		return nil, fmt.Errorf("%s", i18n.Message(
			"agent.tsnet.error.whois",
			"failed to get whois: {{.Error}}",
			map[string]any{"Error": err},
		))
	}

	log.Debug("%s", i18n.Message(
		"agent.tsnet.got_whois",
		"Got whois for peer {{.NodeID}}: {{.Whois}}",
		map[string]any{"NodeID": id, "Whois": whois},
	))
	return &whois.Node.Hostinfo, nil
}
