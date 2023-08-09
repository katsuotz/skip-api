package dto

import "gitlab.com/katsuotz/skip-api/entity"

type SyncPagination struct {
	Data       []entity.Sync `json:"data"`
	Pagination Pagination    `json:"pagination"`
}
