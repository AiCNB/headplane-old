package hpagent

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/tale/headplane/agent/internal/i18n"
	"github.com/tale/headplane/agent/internal/tsnet"
	"github.com/tale/headplane/agent/internal/util"
	"tailscale.com/tailcfg"
)

// Represents messages from the Headplane master
type RecvMessage struct {
	NodeIDs []string
}

type SendMessage struct {
	Type string
	Data any
}

// Starts listening for messages from stdin
func FollowMaster(agent *tsnet.TSAgent) {
	log := util.GetLogger()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Bytes()

		var msg RecvMessage
		err := json.Unmarshal(line, &msg)
		if err != nil {
			log.Error("%s", i18n.Message(
				"agent.hpagent.unmarshal_failed",
				"Unable to unmarshal message: {{.Error}}",
				map[string]any{"Error": err},
			))
			log.Debug("%s", i18n.Message(
				"agent.hpagent.full_error",
				"Full Error: {{.Error}}",
				map[string]any{"Error": err},
			))
			continue
		}

		log.Debug("%s", i18n.Message(
			"agent.hpagent.received_message",
			"Received message from master: {{.Message}}",
			map[string]any{"Message": string(line)},
		))

		if len(msg.NodeIDs) == 0 {
			log.Debug("%s", i18n.Message(
				"agent.hpagent.missing_node_ids",
				"Message received had no node IDs",
				nil,
			))
			log.Debug("%s", i18n.Message(
				"agent.hpagent.full_message",
				"Full message: {{.Message}}",
				map[string]any{"Message": string(line)},
			))
			continue
		}

		// Accumulate the results since we invoke via gofunc
		results := make(map[string]*tailcfg.HostinfoView)
		mu := sync.Mutex{}
		wg := sync.WaitGroup{}

		for _, nodeID := range msg.NodeIDs {
			wg.Add(1)
			go func(nodeID string) {
				defer wg.Done()
				result, err := agent.GetStatusForPeer(nodeID)
				if err != nil {
					log.Error("%s", i18n.Message(
						"agent.hpagent.status_failed",
						"Unable to get status for node {{.NodeID}}: {{.Error}}",
						map[string]any{"NodeID": nodeID, "Error": err},
					))
					return
				}

				if result == nil {
					log.Debug("%s", i18n.Message(
						"agent.hpagent.status_missing",
						"No status for node {{.NodeID}}",
						map[string]any{"NodeID": nodeID},
					))
					return
				}

				mu.Lock()
				results[nodeID] = result
				mu.Unlock()
			}(nodeID)
		}

		wg.Wait()

		// Send the results back to the Headplane master
		log.Debug("%s", i18n.Message(
			"agent.hpagent.sending_status",
			"Sending status back to master: {{.Results}}",
			map[string]any{"Results": results},
		))
		log.Msg(&SendMessage{
			Type: "status",
			Data: results,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("%s", i18n.Message(
			"agent.hpagent.stdin_error",
			"Error reading from stdin: {{.Error}}",
			map[string]any{"Error": err},
		))
	}
}
