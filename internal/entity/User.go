package entity

type User struct {
	Id       string `json:"id" format:"uuid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Leaves   string `json:"leaves"`
}
