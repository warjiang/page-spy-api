package room

import (
	"encoding/json"
	"fmt"

	"github.com/warjiang/page-spy-api/api/event"
	"github.com/warjiang/page-spy-api/api/room"
)

func roomMessageToPackage(msg *room.Message, from *event.Address) (*event.Package, error) {
	bs, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, fmt.Errorf("room message encode failed, %w", err)
	}

	return &event.Package{
		From:       from,
		CreatedAt:  msg.CreatedAt,
		RequestId:  msg.RequestId,
		RoutingKey: msg.Type,
		Content:    bs,
	}, nil
}

func packageToRoomMessage(pkg *event.Package) (*room.Message, error) {
	content := room.NewMessageContent(pkg.RoutingKey)
	err := json.Unmarshal(pkg.Content, content)
	if err != nil {
		return nil, fmt.Errorf("raw message decode failed, %w", err)
	}

	return &room.Message{
		CreatedAt: pkg.CreatedAt,
		RequestId: pkg.RequestId,
		Type:      pkg.RoutingKey,
		Content:   content,
	}, nil
}
