package entity

type ScoreSiswa struct {
	ID           int     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	SiswaKelasID int     `gorm:"index" json:"siswa_kelas_id"`
	Score        float64 `json:"score"`
	Base
}

func (ScoreSiswa) TableName() string {
	return "score_siswa"
}
