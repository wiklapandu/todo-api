package types

type TypeResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type RequestRegister struct {
	Name        string `json:"name" xml:"name" form:"name" validate:"required"`
	Email       string `json:"email" xml:"email" form:"email" validate:"required,email"`
	Pass        string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	ConfirmPass string `json:"confirmPassword" xml:"confirmPassword" form:"confirmPassword" validate:"eqfield=Pass"`
}

type RequestLogin struct {
	Email string `json:"email" xml:"email" form:"email" validate:"required,email"`
	Pass  string `json:"password" xml:"password" form:"password" validate:"required"`
}
