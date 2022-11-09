package client

import (
	"errors"
	"sync"

	"DaisyClubHouse/gobang/manager/player"
	"golang.org/x/exp/slog"
)

type PlayerClientManager struct {
	clients sync.Map //  map[string]*player.Client // client_id -> Client
	idMap   sync.Map // player_id -> client_id
}

func PlayerClientManagerProvider() *PlayerClientManager {
	return &PlayerClientManager{
		// clients: make(map[string]*player.Client),
	}
}

func (m *PlayerClientManager) ClientConnected(client *player.Client) {
	m.clients.Store(client.ID, client)
	slog.Info("玩家客户端加入列表", slog.String("client_id", client.ID))
}

func (m *PlayerClientManager) ClientDisconnected(clientID string) {
	m.clients.Delete(clientID)
	slog.Info("玩家客户端离开列表", slog.String("client_id", clientID))
}

func (m *PlayerClientManager) AssociatedID(clientID, playerID string) {
	m.idMap.Store(playerID, clientID)
	slog.Info("玩家客户端ID关联", slog.String("player_id", playerID), slog.String("client_id", clientID))
}

func (m *PlayerClientManager) GetClientByClientID(clientID string) (*player.Client, error) {
	client, ok := m.clients.Load(clientID)
	if !ok {
		return nil, errors.New("client not found")
	}
	c := client.(*player.Client)
	return c, nil
}

func (m *PlayerClientManager) GetClientByPlayerID(playerID string) (*player.Client, error) {
	clientID, ok := m.idMap.Load(playerID)
	if !ok {
		return nil, errors.New("client not found")
	}

	c, err := m.GetClientByClientID(clientID.(string))
	if err != nil {
		return nil, err
	}
	return c, nil
}
