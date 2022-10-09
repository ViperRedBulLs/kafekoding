package handlers

type ChatAction struct {
	Type   string `form:"type"`
	RoomID uint   `form:"room_id"`
}

type AdminAction struct {
	TableName string `form:"table_name"`
	FieldID   uint   `form:"field_id"`
}
