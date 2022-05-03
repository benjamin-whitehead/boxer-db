package api

import "github.com/benjamin-whitehead/boxer-db/m/v2/replication"

type Request struct {
	Value interface{} `json:"value"`
}

type RoleResponse struct {
	Role string `json:"role"`
}

type InvalidRequestToFollowerResponse struct {
	Message string `json:"message"`
}

type ReplicationRequest struct {
	Log []replication.Entry `json:"log"`
}
