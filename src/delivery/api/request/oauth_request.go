package request

type OAuthPostAuthentication struct {
	AuthenticationType string `json:"authentication_type" validate:"required"`
	AuthenticationID   string `json:"authentication_id" validate:"required"`
	Password           string `json:"password" validate:"required"`
}

type OAuthPostCallBack struct {
	AccessToken      string `json:"authentication_type" validate:"required"`
	AuthenticationID string `json:"authentication_id" validate:"required"`
	Password         string `json:"password" validate:"required"`
}
