package postgres

import (
	"context"

	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.IUserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *userRepository) GetUserByFirebaseUID(ctx context.Context, firebaseUid string) (*models.User, error) {
	var user models.User
	err := u.db.Where("firebase_uid = ?", firebaseUid).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) UpsertUser(ctx context.Context, user *models.User) error {
	return u.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"firebase_uid", "display_name"}),
	}).Create(user).Error
}

func (u *userRepository) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	var users []models.User
	searchPattern := "%" + query + "%"
	err := u.db.WithContext(ctx).
		Where("email ILIKE ? OR display_name ILIKE ?", searchPattern, searchPattern).
		Limit(limit).
		Find(&users).Error
	return users, err
}
