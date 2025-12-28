package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	rolecommand "github.com/tranvuongduy2003/go-copilot/internal/application/role/command"
	rolequery "github.com/tranvuongduy2003/go-copilot/internal/application/role/query"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type RoleHandler struct {
	createRoleHandler               *rolecommand.CreateRoleHandler
	updateRoleHandler               *rolecommand.UpdateRoleHandler
	deleteRoleHandler               *rolecommand.DeleteRoleHandler
	assignPermissionToRoleHandler   *rolecommand.AssignPermissionToRoleHandler
	removePermissionFromRoleHandler *rolecommand.RemovePermissionFromRoleHandler
	setRolePermissionsHandler       *rolecommand.SetRolePermissionsHandler
	listRolesHandler                *rolequery.ListRolesHandler
	getRoleHandler                  *rolequery.GetRoleHandler
	getUsersWithRoleHandler         *rolequery.GetUsersWithRoleHandler
	validator                       *validator.Validator
	logger                          logger.Logger
}

type RoleHandlerParams struct {
	CreateRoleHandler               *rolecommand.CreateRoleHandler
	UpdateRoleHandler               *rolecommand.UpdateRoleHandler
	DeleteRoleHandler               *rolecommand.DeleteRoleHandler
	AssignPermissionToRoleHandler   *rolecommand.AssignPermissionToRoleHandler
	RemovePermissionFromRoleHandler *rolecommand.RemovePermissionFromRoleHandler
	SetRolePermissionsHandler       *rolecommand.SetRolePermissionsHandler
	ListRolesHandler                *rolequery.ListRolesHandler
	GetRoleHandler                  *rolequery.GetRoleHandler
	GetUsersWithRoleHandler         *rolequery.GetUsersWithRoleHandler
	Validator                       *validator.Validator
	Logger                          logger.Logger
}

func NewRoleHandler(params RoleHandlerParams) *RoleHandler {
	return &RoleHandler{
		createRoleHandler:               params.CreateRoleHandler,
		updateRoleHandler:               params.UpdateRoleHandler,
		deleteRoleHandler:               params.DeleteRoleHandler,
		assignPermissionToRoleHandler:   params.AssignPermissionToRoleHandler,
		removePermissionFromRoleHandler: params.RemovePermissionFromRoleHandler,
		setRolePermissionsHandler:       params.SetRolePermissionsHandler,
		listRolesHandler:                params.ListRolesHandler,
		getRoleHandler:                  params.GetRoleHandler,
		getUsersWithRoleHandler:         params.GetUsersWithRoleHandler,
		validator:                       params.Validator,
		logger:                          params.Logger,
	}
}

func (handler *RoleHandler) List(writer http.ResponseWriter, request *http.Request) {
	query := rolequery.ListRolesQuery{}
	roles, err := handler.listRolesHandler.Handle(request.Context(), query)
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

func (handler *RoleHandler) Get(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	query := rolequery.GetRoleQuery{RoleID: roleID}
	roleDTO, err := handler.getRoleHandler.Handle(request.Context(), query)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.RoleWithPermissionsResponse{
		ID:            roleDTO.ID,
		Name:          roleDTO.Name,
		DisplayName:   roleDTO.DisplayName,
		Description:   roleDTO.Description,
		Permissions:   roleDTO.Permissions,
		PermissionIDs: roleDTO.PermissionIDs,
		IsSystem:      roleDTO.IsSystem,
		IsDefault:     roleDTO.IsDefault,
		Priority:      roleDTO.Priority,
		CreatedAt:     roleDTO.CreatedAt,
		UpdatedAt:     roleDTO.UpdatedAt,
	})
}

func (handler *RoleHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.CreateRoleRequest
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

	cmd := rolecommand.CreateRoleCommand{
		Name:          requestBody.Name,
		DisplayName:   requestBody.DisplayName,
		Description:   requestBody.Description,
		PermissionIDs: requestBody.PermissionIDs,
	}

	roleDTO, err := handler.createRoleHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	location := "/api/v1/roles/" + roleDTO.ID.String()
	response.CreatedWithLocation(writer, dto.RoleResponse{
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
	}, location)
}

func (handler *RoleHandler) Update(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	var requestBody dto.UpdateRoleRequest
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

	cmd := rolecommand.UpdateRoleCommand{
		RoleID:      roleID,
		DisplayName: requestBody.DisplayName,
		Description: requestBody.Description,
	}

	roleDTO, err := handler.updateRoleHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.RoleResponse{
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
	})
}

func (handler *RoleHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	cmd := rolecommand.DeleteRoleCommand{RoleID: roleID}
	if err := handler.deleteRoleHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.NoContent(writer)
}

func (handler *RoleHandler) SetPermissions(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	var requestBody dto.SetRolePermissionsRequest
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

	cmd := rolecommand.SetRolePermissionsCommand{
		RoleID:        roleID,
		PermissionIDs: requestBody.PermissionIDs,
	}

	roleDTO, err := handler.setRolePermissionsHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.RoleResponse{
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
	})
}

func (handler *RoleHandler) AddPermission(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	permissionID, err := handler.parsePermissionID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid permission id")
		return
	}

	cmd := rolecommand.AssignPermissionToRoleCommand{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	roleDTO, err := handler.assignPermissionToRoleHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.RoleResponse{
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
	})
}

func (handler *RoleHandler) RemovePermission(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	permissionID, err := handler.parsePermissionID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid permission id")
		return
	}

	cmd := rolecommand.RemovePermissionFromRoleCommand{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	roleDTO, err := handler.removePermissionFromRoleHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.RoleResponse{
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
	})
}

func (handler *RoleHandler) GetUsersWithRole(writer http.ResponseWriter, request *http.Request) {
	roleID, err := handler.parseRoleID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid role id")
		return
	}

	query := rolequery.GetUsersWithRoleQuery{RoleID: roleID}
	users, err := handler.getUsersWithRoleHandler.Handle(request.Context(), query)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, userDTO := range users {
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

	response.Success(writer, userResponses)
}

func (handler *RoleHandler) parseRoleID(request *http.Request) (uuid.UUID, error) {
	roleIDParam := chi.URLParam(request, "id")
	return uuid.Parse(roleIDParam)
}

func (handler *RoleHandler) parsePermissionID(request *http.Request) (uuid.UUID, error) {
	permissionIDParam := chi.URLParam(request, "permissionId")
	return uuid.Parse(permissionIDParam)
}
