package appmodels

type BackupMessage struct {
	Content  string `db:"content"`
	RoomID   int64  `db:"room_id"`
	Username string `db:"username"`
	UserID   int64  `db:"client_id"`
}

type BackupMessages map[int64][]BackupMessage

type BackupRoom struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
