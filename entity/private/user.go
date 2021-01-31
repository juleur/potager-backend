package private

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	ID             int
	Username       string
	Description    string
	Email          string
	Password       string
	RefreshToken   string
	RegisteredAt   time.Time
	VerifiedAt     *sql.NullString
	LastLoggedInAt time.Time
	UserFarmer     UserFarmer
	UserPermission UserPermission
	State          string
}

////////////////////////////////////////////////
type UserPermission struct {
	FailedLoginAttempts int
	LockedUntil         *sql.NullString
	Status              int
}

////////////////////////////////////////////////
type UserFarmer struct {
	ImgUrl            string
	Description       string
	Commune           string
	Coordonnees       float64
	CreatedAt         time.Time
	UpdatedAt         *sql.NullString
	TemporaryDisabled bool
}

func (u User) IsVerified() bool {
	return u.VerifiedAt.Valid
}

func (u User) LockedUntil() string {
	datetime := u.UserPermission.LockedUntil.String
	if len(datetime) == 0 {
		return "00:00:00 00h00"
	}
	min, err := strconv.ParseInt(datetime[14:16], 10, 64)
	if err != nil {
		return fmt.Sprintf("%s/%s/%s %sh%s", datetime[8:10], datetime[5:7], datetime[0:4], datetime[11:13], datetime[14:16])
	} else if min < 9 {
		return fmt.Sprintf("%s/%s/%s %sh0%d", datetime[8:10], datetime[5:7], datetime[0:4], datetime[11:13], min+1)
	} else {
		return fmt.Sprintf("%s/%s/%s %sh%d", datetime[8:10], datetime[5:7], datetime[0:4], datetime[11:13], min+1)
	}
}
