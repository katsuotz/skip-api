package siti_entity

type Keahlian struct {
	IDKeahlian   int    `json:"id_keahlian"`
	NamaKeahlian string `json:"nama_keahlian"`
}

func (Keahlian) TableName() string {
	return "t_keahlian"
}
