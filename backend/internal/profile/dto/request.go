package dto

type CreateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
}
