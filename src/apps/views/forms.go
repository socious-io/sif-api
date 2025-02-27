package views

type ProjectForm struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type AuthForm struct {
	RedirectURL string `json:"redirect_url" validate:"required"`
}

type SessionForm struct {
	Code string `json:"code" validate:"required"`
}
