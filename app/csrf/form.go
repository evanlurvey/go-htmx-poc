package csrf

type FormData struct {
	CSRFToken string `form:"csrf_token"`
}

func (c FormData) GetCSRFToken() string {
	return c.CSRFToken
}
