package api

type Request struct {
	Value interface{} `json:"value"`
}

type RoleResponse struct {
	Role string `json:"role"`
}

type InvalidRequestToFollowerResponse struct {
	Message string `json:"message"`
}
