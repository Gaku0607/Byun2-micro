package models

//會員參數
type Member struct {
	Name string `json:"name" form:"name"`

	Banlancer float64 `json:"banlancer" form:"banlancer"`

	IsSeller bool `json:"is_seller" form:"is_seller"`
}

//註冊參數
type Registry struct {
	Name string `json:"name" form:"name"`

	Pwd string `json:"pwd" form:"pwd"`

	RePwd string `json:"re_pwd" form:"re_pwd"`

	Email string `json:"email" form:"email"`
}

//登入參數
type Login struct {
	Name string `json:"name" form:"name"`

	Pwd string `json:"pwd" form:"pwd"`
}

//修改用戶參數
type MotidyMember struct {
	Email string `json:"email" form:"email"`

	Pwd string `json:"pwd" form:"pwd"`
}

//註冊驗證參數
type Code struct {
	Code string `json:"code" form:"code"`
}

//用戶ＩＤ
type MemberId struct {
	Id int64 `json:"id" form:"id"`
}
