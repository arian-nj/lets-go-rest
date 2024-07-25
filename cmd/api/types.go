package main

type ApiError struct {
	Error string `json:"error"`
}
type envelope map[string]interface{}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
