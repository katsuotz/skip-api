package siti_entity

type TahunPelajaran struct {
	IDTahunPelajaran int    `json:"id_tahun_pelajaran"`
	TahunPelajaran   string `json:"tahun_pelajaran"`
	Aktif            bool   `json:"aktif"`
}

func (TahunPelajaran) TableName() string {
	return "t_tahun_pelajaran"
}
