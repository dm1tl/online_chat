package appmodels

type CreateRoomReq struct {
	ID       int64  `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetRoomResp struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type AddClientReq struct {
	RoomID   int64
	ClientID int64
	Username string
	Password string
}

type AddMessageReq struct {
	Content  string
	RoomID   int64
	Username string
	UserID   int64
}
