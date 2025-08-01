package room

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/warjiang/page-spy-api/api/event"
	"github.com/warjiang/page-spy-api/api/room"
	localRpc "github.com/warjiang/page-spy-api/rpc"
)

func NewRemoteRpcRoomManager(addressManager *localRpc.AddressManager,
	rpcManager *localRpc.RpcManager,
	event event.EventEmitter,
	localRoomManager *LocalRoomManager) *RemoteRpcRoomManager {

	return &RemoteRpcRoomManager{
		BasicManager:     *NewBasicManager(),
		AddressManager:   addressManager,
		rpcManager:       rpcManager,
		event:            event,
		localRoomManager: localRoomManager,
	}
}

type RemoteRpcRoomManager struct {
	BasicManager
	AddressManager   *localRpc.AddressManager
	rpcManager       *localRpc.RpcManager
	event            event.EventEmitter
	localRoomManager *LocalRoomManager
}

func (r *RemoteRpcRoomManager) getRpcByAddress(address *event.Address) (*localRpc.RpcClient, error) {
	if address == nil {
		return nil, fmt.Errorf("get rpc address is nil")
	}
	rpc := r.rpcManager.GetRpcByAddress(address)
	if rpc == nil {
		return nil, fmt.Errorf("rpc client %s not found", address.MachineID)
	}

	return rpc, nil
}

func (r *RemoteRpcRoomManager) Start() {
	r.start()
	log.Info("remote rpc room manager start")
}

func (r *RemoteRpcRoomManager) GetRooms(ctx context.Context) ([]room.RemoteRoom, error) {
	rooms := make([]room.RemoteRoom, 0)
	for _, c := range r.rpcManager.GetRpcList() {
		req := NewRpcLocalRoomManagerRequest()
		res := NewRpcLocalRoomManagerResponse()
		err := c.Call(ctx, "LocalRpcRoomManager.GetRooms", req, res)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, res.GetRooms()...)
	}

	return rooms, nil
}

func (r *RemoteRpcRoomManager) GetRoomsByGroup(ctx context.Context, tags map[string]string) ([]room.RemoteRoom, error) {
	rooms := make([]room.RemoteRoom, 0)
	for _, c := range r.rpcManager.GetRpcList() {
		req := NewRpcLocalRoomManagerRequest()
		req.Tags = tags
		res := NewRpcLocalRoomManagerResponse()
		err := c.Call(ctx, "LocalRpcRoomManager.GetRoomsByGroup", req, res)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, res.GetRooms()...)
	}

	return rooms, nil
}

func (r *RemoteRpcRoomManager) ListRooms(ctx context.Context, tags map[string]string) ([]*room.Info, error) {
	rooms, err := r.GetRoomsByGroup(ctx, tags)
	if err != nil {
		return nil, err
	}

	infos := make([]*room.Info, 0)
	for _, r := range rooms {
		i := r.GetInfo()
		i.Secret = "-"
		infos = append(infos, i)
	}

	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].CreatedAt.After(infos[j].CreatedAt)
	})

	return infos, nil
}

func (r *RemoteRpcRoomManager) CreateConnection() *room.Connection {
	address := r.AddressManager.GeneratorConnectionAddress()
	return &room.Connection{
		Address:   address,
		CreatedAt: time.Now(),
	}
}
func (r *RemoteRpcRoomManager) CreateLocalRoom(ctx context.Context, info *room.Info) (room.Room, error) {
	return r.localRoomManager.CreateRoom(ctx, info)
}

func (r *RemoteRpcRoomManager) CreateRoom(ctx context.Context, info *room.Info) (room.Room, error) {
	if r.AddressManager.IsSelfMachineAddress(info.Address) {
		return r.CreateLocalRoom(ctx, info)
	}

	return r.CreateRemoteRoom(ctx, info)
}

func (r *RemoteRpcRoomManager) UpdateRoomOption(ctx context.Context, info *room.Info) (*room.Info, error) {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return nil, err
	}
	err = rpcClient.Call(ctx, "LocalRpcRoomManager.UpdateRoomOption", req, res)
	if err != nil {
		return nil, err
	}
	return res.Room.Info, nil
}

func (r *RemoteRpcRoomManager) CreateRemoteRoom(ctx context.Context, info *room.Info) (room.Room, error) {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return nil, err
	}
	err = rpcClient.Call(ctx, "LocalRpcRoomManager.CreateRoom", req, res)
	if err != nil {
		return nil, err
	}
	return res.Room, nil
}

func (r *RemoteRpcRoomManager) GetRoomUsers(ctx context.Context, info *room.Info) ([]*room.Connection, error) {
	room, err := r.GetRoom(ctx, info)
	if err != nil {
		return nil, err
	}

	if room == nil || reflect.ValueOf(room).IsNil() {
		return nil, fmt.Errorf("room %s not found", info.Address.ID)
	}

	return room.GetRoomUsers(), nil
}

func (r *RemoteRpcRoomManager) GetRoom(ctx context.Context, info *room.Info) (room.Room, error) {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return nil, err
	}
	err = rpcClient.Call(ctx, "LocalRpcRoomManager.GetRoom", req, res)
	if err != nil {
		return nil, err
	}
	return res.Room, nil
}

func (r *RemoteRpcRoomManager) RemoveRoom(ctx context.Context, info *room.Info) error {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return err
	}

	return rpcClient.Call(ctx, "LocalRpcRoomManager.RemoveRoom", req, res)
}

func (r *RemoteRpcRoomManager) LeaveRoom(ctx context.Context, info *room.Info, connection *room.Connection) error {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	req.Connection = connection
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return err
	}

	return rpcClient.Call(ctx, "LocalRpcRoomManager.LeaveRoom", req, res)
}

func (r *RemoteRpcRoomManager) ForceJoinRoom(ctx context.Context, connection *room.Connection, opt *room.Info, roomOpt *room.Info) (room.RemoteRoom, error) {
	rm, err := r.JoinRoom(ctx, connection, opt)
	if err != nil {
		re, ok := err.(*room.Error)
		if !ok || re.Code != room.RoomNotFoundError {
			return nil, err
		}

		_, err = r.CreateRoom(ctx, roomOpt)
		if err != nil {
			return nil, err
		}

		rm, err = r.JoinRoom(ctx, connection, opt)
		if err != nil {
			return nil, fmt.Errorf("force create room error %w", err)
		}
	}

	return rm, nil
}

func (r *RemoteRpcRoomManager) JoinRoom(ctx context.Context, connection *room.Connection, opt *room.Info) (room.RemoteRoom, error) {
	room, err := r.GetRoom(ctx, opt)
	if err != nil {
		return nil, err
	}

	remoteRoom, err := NewRemoteRoom(connection, opt, r.event, room)
	if err != nil {
		return nil, err
	}

	err = remoteRoom.Start(ctx)
	if err != nil {
		return nil, err
	}

	err = r.joinRoom(ctx, connection, opt)
	if err != nil {
		return nil, err
	}

	return remoteRoom, nil
}

func (r *RemoteRpcRoomManager) joinRoom(ctx context.Context, connection *room.Connection, info *room.Info) error {
	req := NewRpcLocalRoomManagerRequest()
	req.Info = info
	req.Connection = connection
	res := NewRpcLocalRoomManagerResponse()
	rpcClient, err := r.getRpcByAddress(info.Address)
	if err != nil {
		return err
	}

	err = rpcClient.Call(ctx, "LocalRpcRoomManager.JoinRoom", req, res)
	if err != nil {
		return err
	}

	return nil
}
