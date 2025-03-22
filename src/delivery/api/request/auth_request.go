package request

type AuthPostAuthentication struct {
	AuthenticationType string `json:"authentication_type" validate:"required"`
	AuthenticationID   string `json:"authentication_id" validate:"required"`
	Password           string `json:"password" validate:"required"`
	Pin                string `json:"pin"`
	FcmToken           string `json:"fcm_token"`
	Code               string `json:"code"`
	SourceAppID        string `json:"source_app_id"`
}
