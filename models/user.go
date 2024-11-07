package models

import (
	"context"
	"errors"
	"time"

	"github.com/aungmyozaw92/go-restapi-sta/config"
	"github.com/aungmyozaw92/go-restapi-sta/utils"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Username   string    `gorm:"size:100;not null;unique" json:"username" binding:"required"`
	Name       string    `gorm:"size:100;not null" json:"name" binding:"required"`
	Email      string    `gorm:"size:100;unique;default:null" json:"email"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Mobile     string    `gorm:"size:20" json:"mobile"`
	ImageUrl   string    `json:"image_url"`
	Password   string    `gorm:"size:255;not null" json:"password"`
	IsActive   *bool     `gorm:"not null" json:"is_active"`
	RoleId     int       `gorm:"not null;default:0" json:"role_id" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewUser struct {
	Username   string   `json:"username"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Mobile     string   `json:"mobile"`
	ImageUrl   string   `json:"image_url"`
	Password   string   `json:"password"`
	IsActive   *bool    `json:"is_active"`
	RoleId     int      `json:"role_id"`
}

type LoginInfo struct {
	Token      string   `json:"token"`
	UserId	   int      `json:"user_id"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Mobile     string   `json:"mobile"`
	ImageUrl   string   `json:"image_url"`
	Modules    []AllowedModule `json:"modules"`
}

type AllowedModule struct {
	ModuleName     string `json:"module_name"`
	AllowedActions string `json:"allowed_actions"`
}

func (result *User) PrepareGive() {
	result.Password = ""
}

// node
// returns decoded curosr string
func (s User) GetCursor() string {
	return s.CreatedAt.String()
}


func Login(ctx context.Context, username string, password string) (*LoginInfo, error) {

	db := config.GetDB()
	var err error
	var result LoginInfo

	user := User{}

	err = db.WithContext(ctx).Model(User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &result, errors.New("invalid username or password")
	}
	err = utils.ComparePassword(user.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &result, errors.New("invalid username or password")
	}

	isActive := *user.IsActive
	if !isActive {
		return &result, errors.New("user is disabled")
	}
	token, err := utils.JwtGenerate(user.ID)
	result.Token = token
	result.UserId = user.ID
	result.Name = user.Name
	result.Username = user.Username
	result.Email = user.Email
	result.Phone = user.Phone
	result.ImageUrl = user.ImageUrl

	// if user.RoleId == 0 {
	// 	return nil, errors.New("please assign role")
	// } else {
	// 	var userRole Role
	// 	if err := db.WithContext(ctx).Model(&Role{}).
	// 		Preload("RoleModules").Preload("RoleModules.Module").
	// 		Where("id = ?", user.RoleId).First(&userRole, user.RoleId).Error; err != nil {
	// 		return nil, err
	// 	}
	// 	result.Role = userRole.Name
	// 	var allowedModules []AllowedModule
	// 	for _, rm := range userRole.RoleModules {
	// 		allowedModules = append(allowedModules, AllowedModule{
	// 			ModuleName:     rm.Module.Name,
	// 			AllowedActions: rm.AllowedActions,
	// 		})
	// 	}
	// 	result.Modules = allowedModules
	// }

	if err != nil {
		return &result, err
	}

	return &result, nil
}

// destroy current session
// func Logout(ctx context.Context) (bool, error) {
// 	token, ok := utils.GetTokenFromContext(ctx)
// 	if !ok || token == "" {
// 		return false, errors.New("token is required")
// 	}

// 	// Invalidate the token by storing it in Redis with an expiration
// 	expiration := time.Hour // Match this with your JWT token's expiration time
// 	err := config.SetRedisValue(token, "invalid", expiration)
// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

func GetUser(ctx context.Context, id int) (*User, error) {

	db := config.GetDB()
	var result User

	err := db.WithContext(ctx).First(&result, id).Error

	if err != nil {
		return &result, errors.New("record not found")
	}

	result.PrepareGive()

	return &result, nil
}