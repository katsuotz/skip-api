package entity

type ScoreLog struct {
	ID           int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title        string     `gorm:"type:varchar(50)" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Score        float64    `json:"score"`
	Type         string     `gorm:"type:varchar(20)" json:"type"`
	GuruID       int        `gorm:"index" json:"guru_id"`
	Guru         Guru       `gorm:"foreignKey:guru_id" json:"guru"`
	ScoreSiswaID int        `gorm:"index" json:"score_siswa_id"`
	ScoreSiswa   ScoreSiswa `gorm:"foreignKey:score_siswa_id" json:"score_siswa"`
	Base
}

func (ScoreLog) TableName() string {
	return "score_log"
}
