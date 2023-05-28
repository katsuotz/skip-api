package entity

type Kelas struct {
	ID          int      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	NamaKelas   string   `gorm:"type:varchar(20)" json:"nama_kelas"`
	JurusanID   int      `gorm:"index;index:tahun_jurusan_idx" json:"jurusan_id"`
	TahunAjarID int      `gorm:"index;index:tahun_jurusan_idx" json:"tahun_ajar_id"`
	TahunAjar   string   `gorm:"->;-:migration" json:"tahun_ajar,omitempty"`
	PegawaiID   int      `gorm:"index" json:"pegawai_id"`
	Pegawai     *Pegawai `gorm:"foreignKey:pegawai_id" json:"pegawai,omitempty"`
	Base
}
