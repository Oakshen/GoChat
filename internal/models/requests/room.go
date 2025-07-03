package requests

// CreateRoomRequest 创建聊天室请求
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100" vd:"len($)>0 && len($)<=100; msg:'聊天室名称长度应在1-100字符之间'"`
	Description string `json:"description" binding:"max=500" vd:"len($)<=500; msg:'聊天室描述不能超过500字符'"`
	IsPrivate   bool   `json:"is_private"` // 是否为私有聊天室
}

// UpdateRoomRequest 更新聊天室请求
type UpdateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100" vd:"len($)>0 && len($)<=100; msg:'聊天室名称长度应在1-100字符之间'"`
	Description string `json:"description" binding:"max=500" vd:"len($)<=500; msg:'聊天室描述不能超过500字符'"`
}

// JoinRoomRequest 加入聊天室请求
type JoinRoomRequest struct {
	RoomID uint `json:"room_id" binding:"required" vd:"$>0; msg:'聊天室ID必须大于0'"`
}
