package repo

import (
	"github.com/ksj/car-auction/internal/model"
	"gorm.io/gorm"
)

// UserRepo はユーザーの永続化を担当するリポジトリです
type UserRepo struct{ DB *gorm.DB }

// NewUserRepo は新しい UserRepo を生成します
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{DB: db} }

// FindByEmail は指定されたメールアドレスのユーザーを検索します
// メールアドレスが重複していないか確認する際にも使用されます
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Create は新しいユーザーをデータベースに保存します
func (r *UserRepo) Create(u *model.User) error {
	return r.DB.Create(u).Error
}
