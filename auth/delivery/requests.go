package delivery_auth

type SignupRequest struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Name      string `json:"name"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type EditProfileRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Photo    []byte `json:"photo"`
	Password string `json:"password"`
}
