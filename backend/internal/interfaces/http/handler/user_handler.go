package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/application/command"
	"github.com/tranvuongduy2003/go-copilot/internal/application/query"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type UserHandler struct {
	createUserHandler     *command.CreateUserHandler
	updateUserHandler     *command.UpdateUserHandler
	deleteUserHandler     *command.DeleteUserHandler
	changePasswordHandler *command.ChangePasswordHandler
	activateUserHandler   *command.ActivateUserHandler
	deactivateUserHandler *command.DeactivateUserHandler
	banUserHandler        *command.BanUserHandler
	getUserHandler        *query.GetUserHandler
	listUsersHandler      *query.ListUsersHandler
	validator             *validator.Validator
	logger                logger.Logger
}

func NewUserHandler(
	createUserHandler *command.CreateUserHandler,
	updateUserHandler *command.UpdateUserHandler,
	deleteUserHandler *command.DeleteUserHandler,
	changePasswordHandler *command.ChangePasswordHandler,
	activateUserHandler *command.ActivateUserHandler,
	deactivateUserHandler *command.DeactivateUserHandler,
	banUserHandler *command.BanUserHandler,
	getUserHandler *query.GetUserHandler,
	listUsersHandler *query.ListUsersHandler,
	validator *validator.Validator,
	logger logger.Logger,
) *UserHandler {
	return &UserHandler{
		createUserHandler:     createUserHandler,
		updateUserHandler:     updateUserHandler,
		deleteUserHandler:     deleteUserHandler,
		changePasswordHandler: changePasswordHandler,
		activateUserHandler:   activateUserHandler,
		deactivateUserHandler: deactivateUserHandler,
		banUserHandler:        banUserHandler,
		getUserHandler:        getUserHandler,
		listUsersHandler:      listUsersHandler,
		validator:             validator,
		logger:                logger,
	}
}

func (handler *UserHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.CreateUserRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := command.CreateUserCommand{
		Email:    requestBody.Email,
		Password: requestBody.Password,
		FullName: requestBody.FullName,
	}

	userDTO, err := handler.createUserHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	location := "/api/v1/users/" + userDTO.ID.String()
	response.CreatedWithLocation(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	}, location)
}

func (handler *UserHandler) Get(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	userQuery := query.GetUserQuery{UserID: userID}
	userDTO, err := handler.getUserHandler.Handle(request.Context(), userQuery)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	})
}

func (handler *UserHandler) List(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	page := 1
	if pageStr := queryParams.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 20
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var status *string
	if statusStr := queryParams.Get("status"); statusStr != "" {
		status = &statusStr
	}

	var search *string
	if searchStr := queryParams.Get("search"); searchStr != "" {
		search = &searchStr
	}

	var sortBy *string
	if sortByStr := queryParams.Get("sort_by"); sortByStr != "" {
		sortBy = &sortByStr
	}

	var sortOrder *string
	if sortOrderStr := queryParams.Get("sort_order"); sortOrderStr != "" {
		sortOrder = &sortOrderStr
	}

	var dateFrom *string
	if dateFromStr := queryParams.Get("date_from"); dateFromStr != "" {
		dateFrom = &dateFromStr
	}

	var dateTo *string
	if dateToStr := queryParams.Get("date_to"); dateToStr != "" {
		dateTo = &dateToStr
	}

	listQuery := query.ListUsersQuery{
		Page:      page,
		Limit:     limit,
		Status:    status,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}

	result, err := handler.listUsersHandler.Handle(request.Context(), listQuery)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	userResponses := make([]dto.UserResponse, len(result.Items))
	for i, userDTO := range result.Items {
		userResponses[i] = dto.UserResponse{
			ID:        userDTO.ID,
			Email:     userDTO.Email,
			FullName:  userDTO.FullName,
			Status:    userDTO.Status,
			CreatedAt: userDTO.CreatedAt,
			UpdatedAt: userDTO.UpdatedAt,
			DeletedAt: userDTO.DeletedAt,
		}
	}

	response.Success(writer, dto.PaginatedUsersResponse{
		Items:      userResponses,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	})
}

func (handler *UserHandler) Update(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	var requestBody dto.UpdateUserRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := command.UpdateUserCommand{
		UserID:   userID,
		FullName: requestBody.FullName,
	}

	userDTO, err := handler.updateUserHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	})
}

func (handler *UserHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	cmd := command.DeleteUserCommand{UserID: userID}
	if err := handler.deleteUserHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.NoContent(writer)
}

func (handler *UserHandler) ChangePassword(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	var requestBody dto.ChangePasswordRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := command.ChangePasswordCommand{
		UserID:          userID,
		CurrentPassword: requestBody.CurrentPassword,
		NewPassword:     requestBody.NewPassword,
	}

	if err := handler.changePasswordHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.SuccessWithMessage(writer, nil, "password changed successfully")
}

func (handler *UserHandler) Activate(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	cmd := command.ActivateUserCommand{UserID: userID}
	userDTO, err := handler.activateUserHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	})
}

func (handler *UserHandler) Deactivate(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	cmd := command.DeactivateUserCommand{UserID: userID}
	userDTO, err := handler.deactivateUserHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	})
}

func (handler *UserHandler) Ban(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	var requestBody dto.BanUserRequest
	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
		response.BadRequest(writer, request, "invalid request body")
		return
	}

	if err := handler.validator.Validate(requestBody); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			response.ValidationError(writer, request, validationErrors)
			return
		}
		response.BadRequest(writer, request, err.Error())
		return
	}

	cmd := command.BanUserCommand{
		UserID: userID,
		Reason: requestBody.Reason,
	}

	userDTO, err := handler.banUserHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.UserResponse{
		ID:        userDTO.ID,
		Email:     userDTO.Email,
		FullName:  userDTO.FullName,
		Status:    userDTO.Status,
		CreatedAt: userDTO.CreatedAt,
		UpdatedAt: userDTO.UpdatedAt,
		DeletedAt: userDTO.DeletedAt,
	})
}

func (handler *UserHandler) parseUserID(request *http.Request) (uuid.UUID, error) {
	userIDParam := chi.URLParam(request, "id")
	return uuid.Parse(userIDParam)
}
