package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Beige0905/recipe-backend/internal/model"
	"github.com/Beige0905/recipe-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET 환경변수가 설정되지 않았습니다")
	}
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("비밀번호 해싱 실패: %w", err)
	}
	user := &model.User{
		Email:       email,
		DisplayName: displayName,
	}
	if err := s.userRepo.Create(ctx, user, string(hash)); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *model.User, error) {
	user, passwordHash, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", "", nil, errors.New("이메일 또는 비밀번호가 올바르지 않습니다")
	}
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", nil, err
	}
	refreshToken, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.userRepo.DeleteRefreshToken(ctx, hashToken(refreshToken))
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, *model.User, error) {
	rt, err := s.userRepo.FindRefreshToken(ctx, hashToken(refreshToken))
	if err != nil {
		return "", nil, errors.New("유효하지 않은 리프레시 토큰")
	}
	if time.Now().After(rt.ExpiresAt) {
		_ = s.userRepo.DeleteRefreshToken(ctx, hashToken(refreshToken))
		return "", nil, errors.New("리프레시 토큰이 만료되었습니다")
	}
	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return "", nil, err
	}
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", nil, err
	}
	return accessToken, user, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) ValidateAccessToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("유효하지 않은 토큰")
	}
	return claims.UserID, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *AuthService) createRefreshToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err := s.userRepo.CreateRefreshToken(ctx, userID, hashToken(token), expiresAt); err != nil {
		return "", err
	}
	return token, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
