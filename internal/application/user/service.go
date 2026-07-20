package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"os"

	"payme/internal/application/user/dto"
	"payme/pkg/utils"

	"gorm.io/gorm"
)

type UserService interface {
	UpdateUserProfile(ctx context.Context, req dto.UpdateUserProfileRequest, userID uint, file multipart.File, handler *multipart.FileHeader) (dto.UpdateUserProfileResponse, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) UpdateUserProfile(ctx context.Context, req dto.UpdateUserProfileRequest, userID uint, file multipart.File, handler *multipart.FileHeader) (dto.UpdateUserProfileResponse, error) {
	var user User
	s.db.Where("id = ?", userID).First(&user)
	user.FullName = req.FullName
	user.PhoneNumber = req.PhoneNumber
	user.DateOfBirth = req.DateOfBirth
	user.NextOfKinName = req.NextOfKinName
	user.NextOfKinRelationship = req.NextOfKinRelationship
	user.Address = req.Address


	// Validate file type
	contentType := handler.Header.Get("Content-Type")
	fmt.Println(contentType)
	ext, err := utils.ValidateImageFile(contentType)
	if err != nil {
		return dto.UpdateUserProfileResponse{}, errors.New("invalid image")
	}
	// Create file on server
	tempFile, err := os.CreateTemp("../../uploads/products", "*"+ext)
	if err != nil {
		return dto.UpdateUserProfileResponse{}, errors.New("could not save file")
	}
	defer tempFile.Close()

	//  Copy file (SAFE way)
	//Copies bytes from the uploaded file to the temporary file on the server
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return dto.UpdateUserProfileResponse{}, errors.New("failed to save file")
	}
	s.db.Save(&user)
	return dto.UpdateUserProfileResponse{}, nil
}
