package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Response struct {
		Code    int         `json:"code"`
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message"`
	}
	ResponseInfo struct {
		Status  int    `yaml:"status"`
		Code    int    `yaml:"code"`
		Message string `yaml:"message"`
	}
	UserResponse struct {
		Created        ResponseInfo `yaml:"created"`
		DuplicateEmail ResponseInfo `yaml:"duplicate_email"`
		UpdateFailed   ResponseInfo `yaml:"update_failed"`
		CreateFailed   ResponseInfo `yaml:"create_failed"`
		DeleteFailed   ResponseInfo `yaml:"delete_failed"`
		UserNotFound   ResponseInfo `yaml:"user_not_found"`
	}
	AuthResponse struct {
		UserLocked    ResponseInfo `yaml:"user_locked"`
		EmailNotExist ResponseInfo `yaml:"email_not_exist"`
		PasswordWrong ResponseInfo `yaml:"password_wrong"`
		Unauthorized  ResponseInfo `yaml:"unauthorized"`
		TokenInvalid  ResponseInfo `yaml:"token_invalid"`
	}
	MovieResponse struct {
		Created      ResponseInfo `yaml:"created"`
		DeleteFailed ResponseInfo `yaml:"delete_failed"`
		NotFound     ResponseInfo `yaml:"not_found"`
	}
	NormalResponse struct {
		Success        ResponseInfo `yaml:"success"`
		NotFound       ResponseInfo `yaml:"not_found"`
		TimeOut        ResponseInfo `yaml:"timeout"`
		BadRequest     ResponseInfo `yaml:"bad_request"`
		Internal       ResponseInfo `yaml:"internal"`
		PermissionDeny ResponseInfo `yaml:"permission_deny"`
	}

	AllResponse struct {
		NormalResponse NormalResponse `yaml:"normal"`
		UserResponse   UserResponse   `yaml:"user"`
		AuthResponse   AuthResponse   `yaml:"auth"`
		MovieResponse  MovieResponse  `yaml:"movie"`
	}
)

func ResponseJson(w http.ResponseWriter, data interface{}, resinfo ResponseInfo) {

	res := &Response{}
	res.Code = resinfo.Code
	res.Message = resinfo.Message
	if data != "" {
		res.Data = data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resinfo.Status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func Load() *AllResponse {

	path := os.Getenv("STATUS_FILE_PATH")
	if path == "" {
		path = "configs/status.yml"
	}
	yamlFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	all := &AllResponse{}

	if err := yaml.Unmarshal(yamlFile, all); err != nil {
		panic(err)
	}
	return all
}

func Normal() NormalResponse {
	return Load().NormalResponse
}
func User() UserResponse {
	return Load().UserResponse
}
func Auth() AuthResponse {
	return Load().AuthResponse
}
func Movie() MovieResponse {
	return Load().MovieResponse
}
