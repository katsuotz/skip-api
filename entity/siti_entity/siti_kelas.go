package siti_entity

type Kelas struct {
	IDKelas          int    `json:"id_kelas"`
	NamaKelas        string `json:"nama_kelas"`
	IDTahunPelajaran int    `json:"id_tahun_pelajaran"`
	Tingkat          int    `json:"tingkat"`
	IDKeahlian       int    `json:"id_keahlian"`
	IDGuru           int    `json:"id_guru"`
}

func (Kelas) TableName() string {
	return "t_kelas"
}
