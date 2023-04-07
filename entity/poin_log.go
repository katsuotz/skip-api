package entity

type PoinLog struct {
	ID          int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string     `gorm:"type:varchar(50)" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Poin        float64    `json:"poin"`
	PoinBefore  float64    `json:"poin_before"`
	PoinAfter   float64    `json:"poin_after"`
	Type        string     `gorm:"type:varchar(20)" json:"type"`
	GuruID      int        `gorm:"index" json:"guru_id"`
	Guru        *Guru      `gorm:"foreignKey:guru_id" json:"guru,omitempty"`
	PoinSiswaID int        `gorm:"index" json:"poin_siswa_id"`
	PoinSiswa   *PoinSiswa `gorm:"foreignKey:poin_siswa_id" json:"poin_siswa,omitempty"`
	Base
}

func (PoinLog) TableName() string {
	return "poin_log"
}
