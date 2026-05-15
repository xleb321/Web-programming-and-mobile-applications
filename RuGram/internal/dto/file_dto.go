package dto

type UploadFileResponse struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	CreatedAt    string `json:"created_at"`
}

type FileResponse struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	CreatedAt    string `json:"created_at"`
}

type UpdateProfileRequest struct {
	DisplayName  *string `json:"display_name,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	AvatarFileID *string `json:"avatar_file_id,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
}

type ProfileResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateUserRequest struct {
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Password     *string `json:"password,omitempty"`
	DisplayName  *string `json:"display_name,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	AvatarFileID *string `json:"avatar_file_id,omitempty"`
}