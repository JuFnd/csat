package delivery_auth

import "net/http"

type IApi interface {
	SendResponse(w http.ResponseWriter, response Response)
    Signin(w http.ResponseWriter, r *http.Request)
    Signup(w http.ResponseWriter, r *http.Request)
    LogoutSession(w http.ResponseWriter, r *http.Request)
	AuthAccept(w http.ResponseWriter, r *http.Request)
}