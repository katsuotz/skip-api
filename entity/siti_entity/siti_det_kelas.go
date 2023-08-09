package siti_entity

type DetKelas struct {
	IDDetKelas int    `json:"id_det_kelas"`
	IDKelas    int    `json:"id_kelas"`
	Nis        string `json:"nis"`
}

func (DetKelas) TableName() string {
	return "t_det_kelas"
}

type DetKelasJoin struct {
	DetKelas
	Siswa
	Bio
	Login
	Nis string `json:"nis"`
}
