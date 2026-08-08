package user

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// normalize inputs here
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	// validate inputs here
	if req.Name == "" {
		return nil, ErrInvalidName
	}

	if req.Email == "" {
		return nil, ErrInvalidEmail
	}

	_, err := s.repository.GetByEmail(ctx, req.Email)
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
