package models

type NameResponse struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Origin    string `json:"origin"`
}

type FullProfileResponse struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	LastName   string `json:"lastname"`
	FirstName  string `json:"firstname"`
	Country    string `json:"country"`
	BirthDate  string `json:"birth_date"`
	ProfileStr string `json:"profile_str"`
}