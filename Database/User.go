package Database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	LineId      string
	DeleteAt    bool
	NoticeRange int
}

func GetUsersFromNoticeRange(limit []int) ([]User, error) {
	var users []User
	var Query = "SELECT id, \"lineID\", \"noticeRange\" FROM users WHERE deleteat = false AND("
	var limitQuery string

	for i, v := range limit {
		if i != 0 {
			limitQuery += fmt.Sprintf("OR \"noticeRange\" = %d ", v)
		} else {
			limitQuery += fmt.Sprintf("\"noticeRange\" = %d ", v)
		}
	}
	limitQuery += ")"
	Query += limitQuery

	stmt, err := Db.Query(Query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for stmt.Next() {
		var u User
		if err := stmt.Scan(&u.Id, &u.LineId, &u.NoticeRange); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func (u *User) SendNotification() ([]Location, error) {
	return getLocationsFromUserID(u.Id)
}

func GetUserByLineID(lineID string) (*User, error) {
	var user User

	err := Db.QueryRow("SELECT id, \"lineID\", \"noticeRange\" FROM users WHERE \"lineID\" = $1 AND deleteat = false", lineID).
		Scan(&user.Id, &user.LineId, &user.NoticeRange)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用戶不存在，返回 nil
		}
		return nil, err // 其他錯誤
	}

	return &user, nil // 返回用戶記錄的指針
}
