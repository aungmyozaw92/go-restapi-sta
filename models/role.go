package models

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/aungmyozaw92/go-restapi-sta/config"
	"github.com/aungmyozaw92/go-restapi-sta/utils"
	"gorm.io/gorm"
)


type Role struct {
	ID          int           `gorm:"primary_key" json:"id"`
	Name        string        `gorm:"index;size:100;not null" json:"name" binding:"required"`
	RoleModules []*RoleModule `gorm:"foreignKey:RoleId"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewRole struct {
	Name           string              `json:"name" binding:"required"`
	AllowedModules []*NewAllowedModule `json:"allowed_modules"`
}

type NewAllowedModule struct {
	ModuleID       int    `json:"moduleId"`
	AllowedActions string `json:"allowedActions"`
}

func extractModuleActions(s string) []string {
	return strings.Split(strings.ToLower(s), ";")
}

func CreateRole(ctx context.Context, input *NewRole) (*Role, error) {

	// check duplicate
	if err := utils.ValidateUnique[Role](ctx, "name", input.Name, 0); err != nil {
		return nil, err
	}

	roleModules, err := mapRoleModules(ctx, input.AllowedModules)
	if err != nil {
		return nil, err
	}

	role := Role{
		Name:        input.Name,
		RoleModules: roleModules,
	}
	db := config.GetDB()
	// tx := db.Begin()
	err = db.WithContext(ctx).Create(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func UpdateRole(ctx context.Context, id int, input *NewRole) (*Role, error) {

	// check role exists
	if err := utils.ValidateResourceId[Role](ctx, id); err != nil {
		return nil, err
	}

	// check duplicate
	if err := utils.ValidateUnique[Role](ctx, "name", input.Name, id); err != nil {
		return nil, err
	}
	roleModules, err := mapRoleModules(ctx, input.AllowedModules)
	if err != nil {
		return nil, err
	}

	role := Role{
		ID:         id,
		Name:       input.Name,
	}

	db := config.GetDB()
	tx := db.Begin()

	// full replace, delete excluded
	err = tx.WithContext(ctx).Model(&role).
		Session(&gorm.Session{FullSaveAssociations: true, SkipHooks: true}).
		Association("RoleModules").Unscoped().Replace(roleModules)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.WithContext(ctx).Model(&role).Updates(map[string]interface{}{
		"Name": input.Name,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &role, tx.Commit().Error
}

func DeleteRole(ctx context.Context, id int) (*Role, error) {

	var role Role
	db := config.GetDB()
	
	err := db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}

	// don't allow if a user is using the role
	count, err := utils.ResourceCountWhere[User](ctx, "role_id = ?", id)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("role has been used")
	}

	tx := db.Begin()
	// delete role
	err = tx.WithContext(ctx).Select("RoleModules").Delete(&role).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &role, tx.Commit().Error
}

func mapRoleModules(ctx context.Context, input []*NewAllowedModule) ([]*RoleModule, error) {

	availabeModuleActions := make(map[int]string, 0) // moduleId:actions
	var modules []Module
	db := config.GetDB()
	if err := db.WithContext(ctx).Find(&modules).Error; err != nil {
		return nil, err
	}
	for _, m := range modules {
		availabeModuleActions[m.ID] = m.Actions
	}

	var roleModules []*RoleModule
	for _, permission := range input {

		availableActionsString, ok := availabeModuleActions[permission.ModuleID]
		if !ok || availableActionsString == "" {
			return nil, errors.New("module_id not found")
		}
		availableActions := extractModuleActions(availableActionsString)
		inputActions := extractModuleActions(permission.AllowedActions)
		for _, action := range inputActions {
			if !slices.Contains(availableActions, action) {
				return nil, errors.New("invalid module action")
			}
		}

		roleModules = append(roleModules, &RoleModule{
			ModuleId:       permission.ModuleID,
			AllowedActions: permission.AllowedActions,
		})
	}
	return roleModules, nil
}

func GetRole(ctx context.Context, id int) (*Role, error) {

	db := config.GetDB()
	var role Role
	if err := db.WithContext(ctx).Preload("RoleModules").Preload("RoleModules.Module").Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
	// return GetResource[Role](ctx, id)

}