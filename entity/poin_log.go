package entity

type PoinLog struct {
	ID           int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title        string     `gorm:"type:varchar(200)" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Penanganan   string     `gorm:"type:text" json:"penanganan,omitempty"`
	TindakLanjut string     `gorm:"type:text" json:"tindak_lanjut,omitempty"`
	Poin         float64    `json:"poin"`
	PoinBefore   float64    `json:"poin_before"`
	PoinAfter    float64    `json:"poin_after"`
	Type         string     `gorm:"index;type:varchar(20)" json:"type"`
	File         string     `gorm:"type:text" json:"file"`
	PegawaiID    int        `gorm:"index" json:"pegawai_id"`
	Pegawai      *Pegawai   `gorm:"foreignKey:pegawai_id" json:"pegawai,omitempty"`
	PoinSiswaID  int        `gorm:"index" json:"poin_siswa_id"`
	PoinSiswa    *PoinSiswa `gorm:"foreignKey:poin_siswa_id" json:"poin_siswa,omitempty"`
	DataPoinID   int        `gorm:"index" json:"data_poin_id"`
	DataPoin     *DataPoin  `gorm:"foreignKey:data_poin_id" json:"data_poin,omitempty"`
	Base
}

func (PoinLog) TableName() string {
	return "poin_log"
}
