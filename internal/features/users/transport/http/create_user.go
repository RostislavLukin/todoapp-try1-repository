package users_transport_http

import (
	"net/http"
	"todoapp-try1/internal/core/domain"
	"todoapp-try1/internal/core/logger"
	"todoapp-try1/internal/core/transport/http/request"
	"todoapp-try1/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required=true,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Version     int     `json:"version"`
	ID          int     `json:"id"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	log := logger.FromContext(r.Context())

	responseHandler := response.NewHTTPResponseHandler(rw, log)

	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			"decode and validate HTTP request",
			err,
		)
		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.userService.CreateUser(r.Context(), userDomain)

	if err != nil {
		responseHandler.ErrorResponse(
			"failed create user",
			err,
		)
		return
	}

	response := dtoFromDomain(userDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(
		dto.FullName,
		dto.PhoneNumber,
	)
}

func dtoFromDomain(domain domain.User) CreateUserResponse {
	return CreateUserResponse{
		FullName:    domain.FullName,
		PhoneNumber: domain.PhoneNumber,
		Version:     domain.Version,
		ID:          domain.ID,
	}
}
