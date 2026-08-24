package dto

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}
