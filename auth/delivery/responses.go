package delivery_auth

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

type ProfileResponse struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Photo     string `json:"photo"`
	BirthDate string `json:"birthday"`
}

type AuthCheckResponse struct {
	Login string `json:"login"`
}
