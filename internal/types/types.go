// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type CreateTaskRequest struct {
}

type CreateTaskResponse struct {
	ID int `json:"id"`
}

type DownloadTaskRequest struct {
	ID int `path:"id"`
}

type DownloadTaskResponse struct {
	State int `json:"state"`
}

type GetTaskRequest struct {
	ID int `path:"id"`
}

type GetTaskResponse struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	State    int    `json:"state"`
}

type LoginRequest struct {
	ID int `json:"id"`
}

type LoginResponse struct {
}

type TranslateRequest struct {
	ID int `path:"id"`
}

type TranslateResponse struct {
}

type UserRequest struct {
}

type UserResponse struct {
}
