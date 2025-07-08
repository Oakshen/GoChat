package services

import (
	"gochat/internal/dal"
	"gochat/internal/models/entities"
)

// RoomService 聊天室服务
type RoomService struct {
	roomDAL       *dal.RoomDAL
	roomMemberDAL *dal.RoomMemberDAL
}

// NewRoomService 创建聊天室服务实例
func NewRoomService() *RoomService {
	return &RoomService{
		roomDAL:       dal.NewRoomDAL(),
		roomMemberDAL: dal.NewRoomMemberDAL(),
	}
}

// CreateRoom 创建聊天室
func (s *RoomService) CreateRoom(room *entities.Room) (*entities.Room, error) {
	return s.roomDAL.Create(room)
}

// GetRoomByID 根据ID获取聊天室
func (s *RoomService) GetRoomByID(roomID uint) (*entities.Room, error) {
	return s.roomDAL.GetByID(roomID)
}

// SearchRoomsByBlurName 模糊搜索聊天室
func (s *RoomService) SearchRoomsByBlurName(name string) ([]*entities.Room, error) {
	return s.roomDAL.GetByBlurName(name)
}

// GetUserRooms 获取用户加入的聊天室
func (s *RoomService) GetUserRooms(userID uint) ([]*entities.Room, error) {
	return s.roomDAL.GetUserRooms(userID)
}

// JoinRoom 用户加入聊天室
func (s *RoomService) JoinRoom(userID, roomID uint) error {
	// 先检查用户是否已经在聊天室中
	isMember, err := s.roomMemberDAL.IsMember(roomID, userID)
	if err != nil {
		return err
	}

	// 如果用户已经在聊天室中，返回成功（幂等操作）
	if isMember {
		return nil
	}

	// 用户不在聊天室中，创建新的成员记录
	member := &entities.RoomMember{
		RoomID: roomID,
		UserID: userID,
		Role:   "member",
	}
	return s.roomMemberDAL.Create(member)
}

// LeaveRoom 用户离开聊天室
func (s *RoomService) LeaveRoom(userID, roomID uint) error {
	return s.roomMemberDAL.Delete(roomID, userID)
}

// CanUserJoinRoom 检查用户是否可以加入聊天室
func (s *RoomService) CanUserJoinRoom(userID, roomID uint) (bool, error) {
	// 检查聊天室是否存在
	room, err := s.roomDAL.GetByID(roomID)
	if err != nil {
		return false, err
	}

	// 如果是私聊，需要验证权限
	if room.IsPrivate {
		return s.roomMemberDAL.IsMember(roomID, userID)
	}

	// 公开聊天室，任何人都可以加入
	return true, nil
}

// GetRoomMembers 获取聊天室成员
func (s *RoomService) GetRoomMembers(roomID uint) ([]*entities.RoomMember, error) {
	return s.roomMemberDAL.GetByRoomID(roomID)
}

// UpdateRoom 更新聊天室信息
func (s *RoomService) UpdateRoom(roomID uint, name, description string) (*entities.Room, error) {
	return s.roomDAL.Update(roomID, name, description)
}

// DeleteRoom 删除聊天室
func (s *RoomService) DeleteRoom(roomID uint) error {
	return s.roomDAL.Delete(roomID)
}
