package dao

import (
	"basic-go/webook/internal/domain"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicated     = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	println(err)
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicate uint16 = 1062
		if me.Number == duplicate {
			return ErrDuplicated
		}
	}
	return err
}

func (dao *UserDAO) FindbyEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) EditProfile(c *gin.Context, user domain.User) error {
	var u = &User{
		NickName: user.NickName,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	}

	err := dao.db.WithContext(c).Model(&u).Where("id=?", u.ID).Updates(map[string]any{
		"nick_name": u.NickName,
		"birthday":  u.Birthday,
		"about_me":  u.AboutMe,
	}).Error
	return err
}

func (dao *UserDAO) Profile(c *gin.Context, user domain.User) (User, error) {
	var u = &User{

		ID: user.ID,
	}
	err := dao.db.WithContext(c).Select("nick_name", "birthday", "about_me").Where("id=?", u.ID).Find(&u).Error
	return *u, err
}

type User struct {
	ID       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	NickName string
	Birthday string
	AboutMe  string

	Ctime int64
	Utime int64
}
