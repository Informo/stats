package entrypoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"informo-statisticator/common"

	"github.com/matrix-org/gomatrix"
	"github.com/matrix-org/gomatrixserverlib"
)

func GetEntryPoints(homeserverURL string, accessToken string) (entrypoints []string, err error) {
	url := homeserverURL + "/_matrix/client/r0/rooms/" + common.InformoRoomID + "/state"
	url = url + "?access_token=" + accessToken

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Server responded with %d", resp.StatusCode)
		return
	}

	var state []gomatrix.Event
	if err = json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return
	}

	aliasesEvents := getAliasesEvents(state)
	entrypoints = getAliases(aliasesEvents)
	return
}

func CountEntryNodes(entrypoints []string) (count int) {
	nodes := make(map[gomatrixserverlib.ServerName]int)

	for _, entrypoint := range entrypoints {
		_, node, err := gomatrixserverlib.SplitID('#', entrypoint)
		if err != nil {
			continue
		}

		nodes[node]++
	}

	count = len(nodes)
	return
}
