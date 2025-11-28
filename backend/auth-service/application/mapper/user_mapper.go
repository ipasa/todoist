package mapper

import (
	"github.com/todoist/backend/auth-service/application/dto"
	"github.com/todoist/backend/auth-service/domain/entity"
)

// ToUserResponseDTO converts User entity to UserResponseDTO
func ToUserResponseDTO(user *entity.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		Provider:  user.Provider,
	}
}
