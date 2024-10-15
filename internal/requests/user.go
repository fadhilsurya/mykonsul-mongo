package requests

type ReqUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ReqLogin struct {
	Email string `json:"email"`
}
