package dto

import "gitlab.com/katsuotz/skip-api/entity"

type PoinLogResponse struct {
	ID          int     `json:"id"`
	Poin        float64 `json:"poin"`
	PoinBefore  float64 `json:"poin_before"`
	PoinAfter   float64 `json:"poin_after"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	GuruID      int     `json:"guru_id"`
	Nip         string  `json:"nip"`
	NamaGuru    string  `json:"nama_guru"`
	entity.Base
}

type PoinLogPagination struct {
	Data       []PoinLogResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

type PoinLogSiswaByKelas struct {
	Kelas entity.Kelas      `json:"kelas"`
	Data  []PoinLogResponse `json:"data"`
}

type PoinLogCountResponse struct {
	Nama  string `json:"nama"`
	Nis   string `json:"nis"`
	Type  string `json:"type"`
	Total *int   `json:"total,omitempty"`
	entity.Base
}

type PoinLogCountPagination struct {
	Data       []PoinLogCountResponse `json:"data"`
	Pagination Pagination             `json:"pagination"`
}

type PoinLogCountGraphResponse struct {
	Total *int `json:"total,omitempty"`
	Month int  `json:"month"`
	Year  int  `json:"year"`
	entity.Base
}
