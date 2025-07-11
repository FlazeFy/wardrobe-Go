package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"wardrobe/config"
	"wardrobe/models"
	"wardrobe/models/others"
	"wardrobe/repositories"
	"wardrobe/utils"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService interface {
	BasicRegister(userReq models.User) (*string, error)
	BasicSignOut(token string) error
	BasicLogin(loginReq others.LoginRequest) (*string, *string, error)
	GoogleRegister(code string) (*string, error)
	GetMyProfile(userID uuid.UUID, role string) (*others.MyProfile, error)
}

type authService struct {
	userRepo        repositories.UserRepository
	adminRepo       repositories.AdminRepository
	googleTokenRepo repositories.GoogleTokenRepository
	redisClient     *redis.Client
}

func NewAuthService(userRepo repositories.UserRepository, adminRepo repositories.AdminRepository, googleTokenRepo repositories.GoogleTokenRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:        userRepo,
		adminRepo:       adminRepo,
		googleTokenRepo: googleTokenRepo,
		redisClient:     redisClient,
	}
}

func (s *authService) BasicRegister(userReq models.User) (*string, error) {
	// Repo : Find By Email
	user, err := s.userRepo.FindByUsernameOrEmail(userReq.Username, userReq.Email)
	if user != nil || err != gorm.ErrRecordNotFound {
		if user != nil {
			return nil, errors.New("username or email has already been used")
		}

		return nil, err
	}

	// Hashing
	user, err = utils.HashPassword(userReq, userReq.Password)
	if err != nil {
		return nil, errors.New("failed hashing password")
	}

	// Service : Create User
	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	return &token, nil
}

func (s *authService) GoogleRegister(code string) (*string, error) {
	var role string

	// Token Exchange
	tokenGoogle, err := config.GetGoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("cant exchange token")
	}

	// Google Client
	client := config.GetGoogleOAuthConfig().Client(context.Background(), tokenGoogle)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	defer resp.Body.Close()

	// Decode Google Response
	var googleUser others.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, errors.New("failed to decode user info")
	}

	// Repo : Find By Email
	cleanUsername := utils.EmailToUsername(googleUser.Email)
	userCheck, err := s.userRepo.FindByUsernameOrEmail(cleanUsername, googleUser.Email)
	if userCheck != nil || err != gorm.ErrRecordNotFound {
		if userCheck != nil {
			return nil, errors.New("username or email has already been used")
		}

		return nil, err
	}

	// Repo : Create User
	user := &models.User{
		Password:       "GOOGLE_SIGN_IN",
		Email:          googleUser.Email,
		Username:       cleanUsername,
		TelegramUserId: nil,
	}
	user, err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	role = "user"

	// Repo : Create Google Token
	googleToken := &models.GoogleToken{
		AccessToken:  tokenGoogle.AccessToken,
		RefreshToken: tokenGoogle.AccessToken,
		Expiry:       tokenGoogle.Expiry,
	}
	err = s.googleTokenRepo.CreateGoogleToken(googleToken, user.ID)
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	// JWT Token Generate
	token, err := utils.GenerateToken(user.ID, role)
	if err != nil {
		return nil, errors.New("failed generating token")
	}

	return &token, nil
}

func (s *authService) BasicLogin(loginReq others.LoginRequest) (*string, *string, error) {
	var role string
	// Model
	var account others.Account

	// Repo : Find User By Email
	user, err := s.userRepo.FindByEmail(loginReq.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}
	if user != nil {
		role = "user"
		account = user
	}

	if account == nil {
		// Repo : Find Admin By Email
		admin, err := s.adminRepo.FindByEmail(loginReq.Email)
		if err != nil {
			return nil, nil, err
		}
		role = "admin"
		account = admin
	}

	// Utils : Check Password
	if err := utils.CheckPassword(account, loginReq.Password); err != nil {
		return nil, nil, errors.New("invalid password")
	}

	// Utils : JWT Token Generate
	token, err := utils.GenerateToken(account.GetID(), role)
	if err != nil {
		return nil, nil, errors.New("failed generating token")
	}

	return &token, &role, nil
}

func (s *authService) BasicSignOut(authHeader string) error {
	// Clean Bearer
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return errors.New("invalid authorization header")
	}

	// Token Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("missing exp in token")
	}

	// Check If Token Expired
	expTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expTime)
	if duration <= 0 {
		return errors.New("token already expired")
	}

	// Save token to Redis blacklist
	err = s.redisClient.Set(context.Background(), tokenString, "blacklisted", duration).Err()
	if err != nil {
		return errors.New("failed to blacklist token")
	}

	return nil
}

func (s *authService) GetMyProfile(userID uuid.UUID, role string) (*others.MyProfile, error) {
	if role == "admin" {
		return s.adminRepo.FindById(userID)
	} else { // user
		return s.userRepo.FindById(userID)
	}
}
