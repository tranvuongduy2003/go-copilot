package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	permissioncommand "github.com/tranvuongduy2003/go-copilot/internal/application/permission/command"
	permissionquery "github.com/tranvuongduy2003/go-copilot/internal/application/permission/query"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type PermissionHandler struct {
	createPermissionHandler      *permissioncommand.CreatePermissionHandler
	updatePermissionHandler      *permissioncommand.UpdatePermissionHandler
	deletePermissionHandler      *permissioncommand.DeletePermissionHandler
	listPermissionsHandler       *permissionquery.ListPermissionsHandler
	getPermissionHandler         *permissionquery.GetPermissionHandler
	getPermissionsForRoleHandler *permissionquery.GetPermissionsForRoleHandler
	validator                    *validator.Validator
	logger                       logger.Logger
}

type PermissionHandlerParams struct {
	CreatePermissionHandler      *permissioncommand.CreatePermissionHandler
	UpdatePermissionHandler      *permissioncommand.UpdatePermissionHandler
	DeletePermissionHandler      *permissioncommand.DeletePermissionHandler
	ListPermissionsHandler       *permissionquery.ListPermissionsHandler
	GetPermissionHandler         *permissionquery.GetPermissionHandler
	GetPermissionsForRoleHandler *permissionquery.GetPermissionsForRoleHandler
	Validator                    *validator.Validator
	Logger                       logger.Logger
}

func NewPermissionHandler(params PermissionHandlerParams) *PermissionHandler {
	return &PermissionHandler{
		createPermissionHandler:      params.CreatePermissionHandler,
		updatePermissionHandler:      params.UpdatePermissionHandler,
		deletePermissionHandler:      params.DeletePermissionHandler,
		listPermissionsHandler:       params.ListPermissionsHandler,
		getPermissionHandler:         params.GetPermissionHandler,
		getPermissionsForRoleHandler: params.GetPermissionsForRoleHandler,
		validator:                    params.Validator,
		logger:                       params.Logger,
	}
}

func (handler *PermissionHandler) List(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()

	var resource *string
	if resourceStr := queryParams.Get("resource"); resourceStr != "" {
		resource = &resourceStr
	}

	query := permissionquery.ListPermissionsQuery{
		Resource: resource,
	}

	permissions, err := handler.listPermissionsHandler.Handle(request.Context(), query)
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

func (handler *PermissionHandler) Get(writer http.ResponseWriter, request *http.Request) {
	permissionID, err := handler.parsePermissionID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid permission id")
		return
	}

	query := permissionquery.GetPermissionQuery{PermissionID: permissionID}
	permission, err := handler.getPermissionHandler.Handle(request.Context(), query)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.PermissionResponse{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Code:        permission.Code,
		Description: permission.Description,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	})
}

func (handler *PermissionHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var requestBody dto.CreatePermissionRequest
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

	cmd := permissioncommand.CreatePermissionCommand{
		Resource:    requestBody.Resource,
		Action:      requestBody.Action,
		Description: requestBody.Description,
	}

	permission, err := handler.createPermissionHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	location := "/api/v1/permissions/" + permission.ID.String()
	response.CreatedWithLocation(writer, dto.PermissionResponse{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Code:        permission.Code,
		Description: permission.Description,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}, location)
}

func (handler *PermissionHandler) Update(writer http.ResponseWriter, request *http.Request) {
	permissionID, err := handler.parsePermissionID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid permission id")
		return
	}

	var requestBody dto.UpdatePermissionRequest
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

	cmd := permissioncommand.UpdatePermissionCommand{
		PermissionID: permissionID,
		Description:  requestBody.Description,
	}

	permission, err := handler.updatePermissionHandler.Handle(request.Context(), cmd)
	if err != nil {
		response.Error(writer, request, err)
		return
	}

	response.Success(writer, dto.PermissionResponse{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Code:        permission.Code,
		Description: permission.Description,
		IsSystem:    permission.IsSystem,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	})
}

func (handler *PermissionHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	permissionID, err := handler.parsePermissionID(request)
	if err != nil {
		response.BadRequest(writer, request, "invalid permission id")
		return
	}

	cmd := permissioncommand.DeletePermissionCommand{PermissionID: permissionID}
	if err := handler.deletePermissionHandler.Handle(request.Context(), cmd); err != nil {
		response.Error(writer, request, err)
		return
	}

	response.NoContent(writer)
}

func (handler *PermissionHandler) parsePermissionID(request *http.Request) (uuid.UUID, error) {
	permissionIDParam := chi.URLParam(request, "id")
	return uuid.Parse(permissionIDParam)
}
