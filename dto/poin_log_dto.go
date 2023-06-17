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
	PegawaiID   int     `json:"pegawai_id"`
	Nip         string  `json:"nip"`
	NamaPegawai string  `json:"nama_pegawai"`
	Nama        *string `json:"nama,omitempty"`
	Nis         *string `json:"nis,omitempty"`
	Foto        *string `json:"foto,omitempty"`
	File        *string `json:"file,omitempty"`
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
	Nama  *string `json:"nama,omitempty"`
	Title *string `json:"title,omitempty"`
	Nis   *string `json:"nis,omitempty"`
	Type  string  `json:"type"`
	Total *int    `json:"total,omitempty"`
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
