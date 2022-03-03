package login

import "github.com/asaskevich/govalidator"

type Login struct {
	Nickname string `json:"nickname" db:"nickname" `
	Email    string `json:"email" db:"email" `
	Password string `json:"password" db:"password" valid:"required"`
	RealIP   string `json:"real_ip" db:"real_ip"`
}

func NewLogin(nickname, email, password, realIp string) *Login {
	return &Login{
		Nickname: nickname,
		Email:    email,
		Password: password,
		RealIP:   realIp,
	}
}

func (m *Login) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
