package siti_entity

type Login struct {
	IDLogin  int    `json:"id_login"`
	IDRujuk  string `json:"id_rujuk"`
	Username string `json:"username"`
	Password string `json:"password"`
	Level    int    `json:"level"`
}

func (Login) TableName() string {
	return "t_login"
}
