package siti_entity

type Guru struct {
	IDGuru           int    `json:"id_guru"`
	IDBio            int    `json:"id_bio"`
	GRNRS            string `json:"grnrs"`
	Nip              string `json:"nip"`
	KategoriKaryawan string `json:"kategori_karyawan"`
	Aktif            int    `json:"aktif"`
}

func (Guru) TableName() string {
	return "t_guru"
}

type GuruBio struct {
	Guru
	Bio
	Login
}
