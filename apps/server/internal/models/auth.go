package models

import (
	authulamodels "github.com/Authula/authula/models"
)

type GetMeResponse struct {
	User    *authulamodels.User    `json:"user"`
	Session *authulamodels.Session `json:"session"`
}

type SignOutRequest struct {
	SessionID  *string `json:"session_id,omitempty"`
	SignOutAll bool    `json:"sign_out_all,omitempty"`
}

type SignOutResponse struct {
	Message string `json:"message"`
}
