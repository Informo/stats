package entrypoints

import (
	"github.com/matrix-org/gomatrix"
)

func getAliasesEvents(state []gomatrix.Event) (aliasesEvents []gomatrix.Event) {
	aliasesEvents = make([]gomatrix.Event, 0)

	for _, stateEvent := range state {
		if stateEvent.Type == "m.room.aliases" {
			aliasesEvents = append(aliasesEvents, stateEvent)
		}
	}

	return
}

func getAliases(aliasesEvents []gomatrix.Event) (aliases []string) {
	aliases = make([]string, 0)

	for _, aliasesEvent := range aliasesEvents {
		slice, ok := aliasesEvent.Content["aliases"].([]interface{})
		if !ok {
			continue
		}

		for _, alias := range slice {
			aliasStr, ok := alias.(string)
			if !ok {
				continue
			}

			aliases = append(aliases, aliasStr)
		}
	}

	return
}
