package dto

import "gitlab.com/katsuotz/skip-api/entity"

type SiswaKelasPoin struct {
	ID           int `json:"id"`
	SiswaKelasID int `json:"siswa_kelas_id"`
	entity.SiswaKelas
	entity.PoinSiswa
}
