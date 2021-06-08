package user

type ParamsLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
