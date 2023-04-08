package repository

import (
	"context"
	"gitlab.com/katsuotz/skip-api/dto"
	"gitlab.com/katsuotz/skip-api/entity"
	"gorm.io/gorm"
)

type LoginLogRepository interface {
	GetLog(ctx context.Context, page int, perPage int, search string) dto.LoginLogPagination
}

type loginLogRepository struct {
	db *gorm.DB
}

func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &loginLogRepository{db: db}
}

func (r *loginLogRepository) GetLog(ctx context.Context, page int, perPage int, search string) dto.LoginLogPagination {
	result := dto.LoginLogPagination{}
	loginLog := entity.LoginLog{}
	temp := r.db.Model(&loginLog)
	if search != "" {
		search = "%" + search + "%"
		temp.Where("nama ilike ?", search)
	}

	temp.Select("nama, foto, role, username, action, user_agent, location, login_logs.created_at").
		Joins("join users on users.id = login_logs.user_id").
		Joins("join profiles on profiles.user_id = users.id").
		Order("login_logs.created_at desc")
	temp.Offset(perPage * (page - 1)).Limit(perPage).Find(&result.Data)

	var totalItem int64
	temp.Offset(-1).Limit(-1).Count(&totalItem)
	result.Pagination.TotalItem = totalItem
	result.Pagination.Page = page
	totalPage := totalItem / int64(perPage)
	if totalItem%int64(perPage) > 0 {
		totalPage++
	}
	result.Pagination.TotalPage = totalPage
	result.Pagination.PerPage = perPage

	return result
}
