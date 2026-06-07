package service

import (
	"context"
	"errors"
	"time"

	"log/slog"

	"firebase.google.com/go/v4/auth"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	"github.com/WillieBam/support_copilot/backend/types/models"
	"gorm.io/gorm"
)

type authService struct {
	userRepo     interfaces.IUserRepository
	firebaseRepo interfaces.IFirebaseRepository
}

type AuthServiceParam struct {
	UserRepo     interfaces.IUserRepository
	FirebaseRepo interfaces.IFirebaseRepository
}

func New(asp AuthServiceParam) interfaces.IAuthService {
	return &authService{
		userRepo:     asp.UserRepo,
		firebaseRepo: asp.FirebaseRepo,
	}
}

// IsUserExists is to verify the user whether is in database.
// If user exists, return true, else false.
func (s *authService) IsUserExists(ctx context.Context, idToken string) bool {
	var uid *auth.Token
	var err error
	if uid, err = s.firebaseRepo.VerifyIDToken(ctx, idToken); err != nil {
		slog.Error("unable to verify user", "err", err)
		return false
	}

	if _, userNotFoundErr := s.userRepo.GetUserByFirebaseUID(ctx, uid.UID); userNotFoundErr != nil {
		slog.Error("unable to find user in database", "err", userNotFoundErr)
		return false
	}
	return true
}

// LoginOrRegister verifies the Firebase token and syncs the user into our PostgreSQL DB.
func (s *authService) LoginOrRegister(ctx context.Context, idToken string) (*models.User, error) {
	// 1. Verify token with Firebase Repository
	token, err := s.firebaseRepo.VerifyIDToken(ctx, idToken)
	if err != nil {
		slog.Error("firebase token verification failed", "err", err)
		return nil, errors.New("invalid credentials")
	}

	user, err := s.userRepo.GetUserByFirebaseUID(ctx, token.UID)

	if err != nil {
		// Check if error is due to record not existing
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("new user detected, executing database sync/seed", "uid", token.UID)

			// Extract details safely from Firebase claims
			email, _ := token.Claims["email"].(string)
			name, _ := token.Claims["name"].(string)

			newUser := &models.User{
				FirebaseUID: token.UID,
				Email:       email,
				DisplayName: name,
				CreatedAt:   time.Now(),
				Scope:       "engineer", // default fallback scope
			}

			if createErr := s.userRepo.CreateUser(ctx, newUser); createErr != nil {
				slog.Error("failed to create user in database", "err", createErr)
				return nil, errors.New("failed to provision user session")
			}

			return newUser, nil
		}

		slog.Error("database query anomaly during login lookup", "err", err)
		return nil, err
	}

	return user, nil
}

/*
newUser := &models.User{
				FirebaseUID: uid,
				Email:       email,
				DisplayName: name,
				CreatedAt:   time.Now(),
				Scope:       "engineer",
			}

			if createErr := s.userRepo.CreateUser(ctx, newUser); createErr != nil {
				return nil, createErr
			}

			user = newUser
*/
/*
 // if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		log.Errorf("user not found: %v", err)
	// 	} else {
	// 		return nil, err
	// 	}
	// }
*/
