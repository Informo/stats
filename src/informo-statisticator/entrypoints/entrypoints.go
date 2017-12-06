package entrypoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"informo-statisticator/common"

	"github.com/matrix-org/gomatrix"
	"github.com/matrix-org/gomatrixserverlib"
)

func GetEntryPoints(homeserverURL string, accessToken string) (entrypoints []string, err error) {
	entrypoints = make([]string, 0)

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

	defer resp.Body.Close()

	var state []gomatrix.Event
	if err = json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return
	}

	aliasesEvents := getAliasesEvents(state)
	aliases := getAliases(aliasesEvents)

	var active bool
	for _, alias := range aliases {
		if active, err = checkEntrypointActive(homeserverURL, alias, accessToken); err != nil {
			return
		} else if active {
			entrypoints = append(entrypoints, alias)
		}
	}

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

func checkEntrypointActive(homeserverURL string, alias string, accessToken string) (active bool, err error) {
	url := homeserverURL + "/_matrix/client/r0/directory/room/" + url.PathEscape(alias)
	url = url + "?access_token=" + accessToken

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	defer resp.Body.Close()

	var body struct {
		RoomID string `json:"room_id"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return
	}

	active = (body.RoomID == common.InformoRoomID)
	return
}
