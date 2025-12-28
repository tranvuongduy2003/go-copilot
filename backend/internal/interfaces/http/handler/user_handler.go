package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	usercommand "github.com/tranvuongduy2003/go-copilot/internal/application/user/command"
	userquery "github.com/tranvuongduy2003/go-copilot/internal/application/user/query"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type UserHandler struct {
	createUserHandler        *usercommand.CreateUserHandler
	updateUserHandler        *usercommand.UpdateUserHandler
	deleteUserHandler        *usercommand.DeleteUserHandler
	changePasswordHandler    *usercommand.ChangePasswordHandler
	activateUserHandler      *usercommand.ActivateUserHandler
	deactivateUserHandler    *usercommand.DeactivateUserHandler
	banUserHandler           *usercommand.BanUserHandler
	assignRoleToUserHandler  *usercommand.AssignRoleToUserHandler
	revokeRoleFromUserHandler *usercommand.RevokeRoleFromUserHandler
	setUserRolesHandler      *usercommand.SetUserRolesHandler
	getUserHandler           *userquery.GetUserHandler
	listUsersHandler         *userquery.ListUsersHandler
	getUserRolesHandler      *userquery.GetUserRolesHandler
	getUserPermissionsHandler *userquery.GetUserPermissionsHandler
	validator                *validator.Validator
	logger                   logger.Logger
}

type UserHandlerParams struct {
	CreateUserHandler         *usercommand.CreateUserHandler
	UpdateUserHandler         *usercommand.UpdateUserHandler
	DeleteUserHandler         *usercommand.DeleteUserHandler
	ChangePasswordHandler     *usercommand.ChangePasswordHandler
	ActivateUserHandler       *usercommand.ActivateUserHandler
	DeactivateUserHandler     *usercommand.DeactivateUserHandler
	BanUserHandler            *usercommand.BanUserHandler
	AssignRoleToUserHandler   *usercommand.AssignRoleToUserHandler
	RevokeRoleFromUserHandler *usercommand.RevokeRoleFromUserHandler
	SetUserRolesHandler       *usercommand.SetUserRolesHandler
	GetUserHandler            *userquery.GetUserHandler
	ListUsersHandler          *userquery.ListUsersHandler
	GetUserRolesHandler       *userquery.GetUserRolesHandler
	GetUserPermissionsHandler *userquery.GetUserPermissionsHandler
	Validator                 *validator.Validator
	Logger                    logger.Logger
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		createUserHandler:         params.CreateUserHandler,
		updateUserHandler:         params.UpdateUserHandler,
		deleteUserHandler:         params.DeleteUserHandler,
		changePasswordHandler:     params.ChangePasswordHandler,
		activateUserHandler:       params.ActivateUserHandler,
		deactivateUserHandler:     params.DeactivateUserHandler,
		banUserHandler:            params.BanUserHandler,
		assignRoleToUserHandler:   params.AssignRoleToUserHandler,
		revokeRoleFromUserHandler: params.RevokeRoleFromUserHandler,
		setUserRolesHandler:       params.SetUserRolesHandler,
		getUserHandler:            params.GetUserHandler,
		listUsersHandler:          params.ListUsersHandler,
		getUserRolesHandler:       params.GetUserRolesHandler,
		getUserPermissionsHandler: params.GetUserPermissionsHandler,
		validator:                 params.Validator,
		logger:                    params.Logger,
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

	cmd := usercommand.CreateUserCommand{
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

	userQuery := userquery.GetUserQuery{UserID: userID}
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

	listQuery := userquery.ListUsersQuery{
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

	cmd := usercommand.UpdateUserCommand{
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

	cmd := usercommand.DeleteUserCommand{UserID: userID}
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

	cmd := usercommand.ChangePasswordCommand{
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

	cmd := usercommand.ActivateUserCommand{UserID: userID}
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

	cmd := usercommand.DeactivateUserCommand{UserID: userID}
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

	cmd := usercommand.BanUserCommand{
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

func (handler *UserHandler) GetRoles(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	query := userquery.GetUserRolesQuery{UserID: userID}
	roles, err := handler.getUserRolesHandler.Handle(request.Context(), query)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, roleDTO := range roles {
		roleResponses[i] = dto.RoleResponse{
			ID:            roleDTO.ID,
			Name:          roleDTO.Name,
			DisplayName:   roleDTO.DisplayName,
			Description:   roleDTO.Description,
			PermissionIDs: roleDTO.PermissionIDs,
			IsSystem:      roleDTO.IsSystem,
			IsDefault:     roleDTO.IsDefault,
			Priority:      roleDTO.Priority,
			CreatedAt:     roleDTO.CreatedAt,
			UpdatedAt:     roleDTO.UpdatedAt,
		}
	}

	response.Success(writer, roleResponses)
}

func (handler *UserHandler) SetRoles(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	var requestBody dto.SetUserRolesRequest
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

	cmd := usercommand.SetUserRolesCommand{
		UserID:  userID,
		RoleIDs: requestBody.RoleIDs,
	}

	userDTO, err := handler.setUserRolesHandler.Handle(request.Context(), cmd)
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

func (handler *UserHandler) AssignRole(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	cmd := usercommand.AssignRoleToUserCommand{
		UserID: userID,
		RoleID: roleID,
	}

	userDTO, err := handler.assignRoleToUserHandler.Handle(request.Context(), cmd)
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

func (handler *UserHandler) RevokeRole(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	cmd := usercommand.RevokeRoleFromUserCommand{
		UserID: userID,
		RoleID: roleID,
	}

	userDTO, err := handler.revokeRoleFromUserHandler.Handle(request.Context(), cmd)
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

func (handler *UserHandler) GetPermissions(writer http.ResponseWriter, request *http.Request) {
	userID, err := handler.parseUserID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid user id")
		return
	}

	query := userquery.GetUserPermissionsQuery{UserID: userID}
	permissions, err := handler.getUserPermissionsHandler.Handle(request.Context(), query)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	permissionResponses := make([]dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		permissionResponses[i] = dto.PermissionResponse{
			ID:          permission.ID,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Code:        permission.Code,
			Description: permission.Description,
			IsSystem:    permission.IsSystem,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		}
	}

	response.Success(writer, permissionResponses)
}

func (handler *UserHandler) parseUserID(request *http.Request) (uuid.UUID, error) {
	userIDParam := chi.URLParam(request, "id")
	return uuid.Parse(userIDParam)
}

func (handler *UserHandler) parseRoleID(request *http.Request) (uuid.UUID, error) {
	roleIDParam := chi.URLParam(request, "roleId")
	return uuid.Parse(roleIDParam)
}
