package user

type UserModel struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Stocks struct {
	Id    int `gorm:"primarykey" json:"id"`
	Name  string
	Stock string
}
