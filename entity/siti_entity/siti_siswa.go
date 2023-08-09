package siti_entity

type Siswa struct {
	Nis   string `json:"nis"`
	IDBio int    `json:"id_bio"`
}

func (Siswa) TableName() string {
	return "t_siswa"
}
