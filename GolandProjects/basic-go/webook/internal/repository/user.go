package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicated
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindbyEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		NickName: u.NickName,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
	}

}

func (repo *UserRepository) EditProfile(c *gin.Context, user domain.User) error {
	err := repo.dao.EditProfile(c, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) Profile(c *gin.Context, user domain.User) (domain.User, error) {
	u, err := repo.dao.Profile(c, user)
	return repo.toDomain(u), err
}
