package utils

import abcitypes "github.com/tendermint/tendermint/abci/types"

func AttributeValue(events []abcitypes.Event, eventType, attrKey string) (string, bool) {
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		for _, attr := range event.Attributes {
			if string(attr.Key) == attrKey {
				return string(attr.Value), true
			}
		}
	}
	return "", false
}
