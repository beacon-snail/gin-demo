package services

import (
	"errors"

	"gin-mysql-demo/dto"
	"gin-mysql-demo/middleware"
	"gin-mysql-demo/models"
	"gin-mysql-demo/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// 检查用户名
	existing, _ := s.repo.FindByUsername(req.Username)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱
	existingEmail, _ := s.repo.FindByUsername(req.Email)
	if existingEmail != nil && existingEmail.ID > 0 {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
	}
	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status == 0 {
		return nil, errors.New("user is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	s.repo.UpdateLastLogin(user.ID)

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

func (s *UserService) GetUser(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return s.toUserResponse(user), nil
}

func (s *UserService) GetUsers(params map[string]string) ([]dto.UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *s.toUserResponse(&user))
	}
	return responses, total, nil
}

func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}

func (s *UserService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Age:       user.Age,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
