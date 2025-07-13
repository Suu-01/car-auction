package repo

import (
	"github.com/ksj/car-auction/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct{ DB *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{DB: db} }

// Email 중복 검사
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u *model.User) error {
	return r.DB.Create(u).Error
}
