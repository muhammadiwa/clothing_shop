package services

import (
    "errors"
    "github.com/yourusername/clothing-shop-api/internal/domain/models"
    "github.com/yourusername/clothing-shop-api/internal/repository"
)

type UserService struct {
    userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{
        userRepository: userRepo,
    }
}

func (s *UserService) Register(user *models.User) error {
    existingUser, _ := s.userRepository.FindByEmail(user.Email)
    if existingUser != nil {
        return errors.New("user already exists")
    }
    return s.userRepository.Create(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    return s.userRepository.FindByID(id)
}

func (s *UserService) UpdateUser(user *models.User) error {
    existingUser, err := s.userRepository.FindByID(user.ID)
    if err != nil {
        return err
    }
    if existingUser == nil {
        return errors.New("user not found")
    }
    return s.userRepository.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
    existingUser, err := s.userRepository.FindByID(id)
    if err != nil {
        return err
    }
    if existingUser == nil {
        return errors.New("user not found")
    }
    return s.userRepository.Delete(id)
}