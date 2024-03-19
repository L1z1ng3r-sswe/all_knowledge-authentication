package models

type Response struct {
	Err      error
	ErrKey   string
	ErrMsg   string
	Code     int
	FileInfo string
}
