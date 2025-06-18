package socket

import (
	"github.com/warjiang/page-spy-api/config"
	"github.com/warjiang/page-spy-api/event"
	"github.com/warjiang/page-spy-api/logger"
	"github.com/warjiang/page-spy-api/room"
	"github.com/warjiang/page-spy-api/rpc"
	"github.com/warjiang/page-spy-api/util"
)

func NewManager(config *config.Config, rpcManager *rpc.RpcManager, addressManager *rpc.AddressManager) (*room.RemoteRpcRoomManager, error) {
	localEvent := event.NewLocalEventEmitter(addressManager, rpcManager)
	localRoomManager := room.NewLocalRoomManager(localEvent, addressManager, int64(config.GetMaxRoomNumber()))
	localRoomManager.Start()
	_, err := event.NewRpcEventEmitter(localEvent, rpcManager)
	if err != nil {
		return nil, err
	}

	_, err = room.NewLocalRpcRoomManager(localRoomManager, rpcManager)
	if err != nil {
		return nil, err
	}

	manager := room.NewRemoteRpcRoomManager(addressManager, rpcManager, localEvent, localRoomManager)
	manager.Start()
	logger.Log().Infof("start rpc server %s successful", addressManager.GetSelfMachineID())
	logger.Log().Infof("local ip %s:%s", util.GetLocalIP(), config.Port)
	return manager, nil
}
