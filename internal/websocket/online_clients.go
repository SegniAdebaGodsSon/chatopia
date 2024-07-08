package websocket

import (
	"sync"
)

type OnlineClients struct {
	clients sync.Map
}

func NewOnlineClientsStore() *OnlineClients {
	return &OnlineClients{}
}

func (activeClients *OnlineClients) StoreClient(sessionID string, client *Client) {
	activeClients.clients.Store(sessionID, client)
}

func (activeClients *OnlineClients) GetClient(sessionID string) (*Client, bool) {
	client, exists := activeClients.clients.Load(sessionID)
	if !exists {
		return nil, false
	}

	return client.(*Client), true
}

func (activeClients *OnlineClients) DeleteClient(sessionID string) {
	activeClients.clients.Delete(sessionID)
}

func (activeClients *OnlineClients) FindMatchingClient(sessionID string) *Client {
	currentClient, exists := activeClients.GetClient(sessionID)
	if !exists {
		return nil
	}

	var bestMatch *Client
	maxCommonInterests := -1

	activeClients.clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		if client.SessionID != currentClient.SessionID &&
			client.ChatType == currentClient.ChatType &&
			client.ChatPartner == nil && currentClient.ChatPartner == nil {

			commonInterests := countCommonInterests(client.Interests, currentClient.Interests)
			if commonInterests > maxCommonInterests {
				bestMatch = client
				maxCommonInterests = commonInterests

				// maximum common interests
				if maxCommonInterests == 3 {
					return false
				}
			}
		}
		return true
	})

	// check if both clients are still available
	if bestMatch != nil && currentClient.ChatPartner == nil && bestMatch.ChatPartner == nil {
		return bestMatch
	}

	var anyClient *Client
	activeClients.clients.Range(func(_, value interface{}) bool {
		client := value.(*Client)
		if client.SessionID != currentClient.SessionID && currentClient.ChatType == currentClient.ChatType {
			anyClient = client
			return false
		}
		return true
	})

	if anyClient != nil && currentClient.ChatPartner == nil && anyClient.ChatPartner == nil {
		return anyClient
	}

	return nil
}

func hasCommonInterests(interests1, interests2 []string) bool {
	for interest1 := range interests1 {
		for interest2 := range interests2 {
			if interest1 == interest2 {
				return true
			}
		}
	}
	return false
}

func countCommonInterests(interests1, interests2 []string) int {
	commonCount := 0
	for interest1 := range interests1 {
		for interest2 := range interests2 {
			if interest1 == interest2 {
				commonCount++
			}
		}
	}
	return commonCount
}

/*
func hasCommonInterests(interests1, interests1 []string) bool {

}
*/