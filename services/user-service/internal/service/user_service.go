package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/username/dist-ecommerce-go/pkg/common"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/models"
	"github.com/username/dist-ecommerce-go/services/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, fullName string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	ValidateUser(ctx context.Context, email, password string) (*models.User, error)
}

type userService struct {
	repo  repository.UserRepository
	cache *common.Cache
}

func NewUserService(repo repository.UserRepository, cache *common.Cache) UserService {
	return &userService{repo: repo, cache: cache}
}

func (s *userService) CreateUser(ctx context.Context, email, password, fullName string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	if s.cache != nil {
		val, err := s.cache.Get(ctx, "user:"+id)
		if err == nil {
			var user models.User
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return &user, nil
			}
		}
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		data, _ := json.Marshal(user)
		_ = s.cache.Set(ctx, "user:"+id, data, time.Hour)
	}

	return user, nil
}

func (s *userService) ValidateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
