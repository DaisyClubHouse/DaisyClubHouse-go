package service

type RoomDomainServiceImpl interface {
	GetRoomList()
}

// RoomDomainService 房间领域服务
type RoomDomainService struct {
}

func NewRoomDomainService() RoomDomainServiceImpl {
	return &RoomDomainService{}
}

func (service *RoomDomainService) GetRoomList() {

}
