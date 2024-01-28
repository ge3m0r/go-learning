package domain

type User struct {
	ID       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	NickName string
	AboutMe  string
	Birthday string
	Ctime    int64
	Utime    int64
}

func (u User) ValidateEmail() string {
	return u.Email
}
