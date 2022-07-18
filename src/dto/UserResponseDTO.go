package dto

type UserResponseDTO struct {
	ID             int
	Auth0ID        string
	FirstName      string
	LastName       string
	Email          string
	PhoneNumber    string
	Gender         int
	Username       string
	DateOfBirth    float32
	Biography      string
	Education      string
	WorkExperience string
	Skills         string
	Interests      string
	Public         bool
}
