package dto
type UpdateUserProfileRequest struct {
	UserID      uint   `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
    NextOfKinName string `json:"next_of_kin_name"`
	NextOfKinRelationship string `json:"next_of_kin_relationship"`
	Address     string `json:"address"`
	ProfilePicture string `json:"profile_picture"`
	
}
type UpdateUserProfileResponse struct {
	Message string `json:"message"`
}
