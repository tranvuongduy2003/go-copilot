package permission

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var resourceNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type Resource struct {
	value string
}

func NewResource(value string) (Resource, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return Resource{}, shared.NewValidationError("resource", "resource cannot be empty")
	}
	if len(normalized) > 100 {
		return Resource{}, shared.NewValidationError("resource", "resource cannot exceed 100 characters")
	}
	if !resourceNameRegex.MatchString(normalized) {
		return Resource{}, shared.NewValidationError("resource", "resource must be lowercase alphanumeric with underscores, starting with a letter")
	}
	return Resource{value: normalized}, nil
}

func (r Resource) String() string {
	return r.value
}

func (r Resource) Equals(other Resource) bool {
	return r.value == other.value
}

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionList   Action = "list"
	ActionManage Action = "manage"
	ActionAssign Action = "assign"
	ActionAdmin  Action = "admin"
)

var standardActions = map[Action]bool{
	ActionCreate: true,
	ActionRead:   true,
	ActionUpdate: true,
	ActionDelete: true,
	ActionList:   true,
	ActionManage: true,
	ActionAssign: true,
	ActionAdmin:  true,
}

func NewAction(value string) (Action, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", shared.NewValidationError("action", "action cannot be empty")
	}
	if len(normalized) > 100 {
		return "", shared.NewValidationError("action", "action cannot exceed 100 characters")
	}
	action := Action(normalized)
	if !standardActions[action] {
		if !resourceNameRegex.MatchString(normalized) {
			return "", shared.NewValidationError("action", "custom action must be lowercase alphanumeric with underscores")
		}
	}
	return action, nil
}

func (a Action) String() string {
	return string(a)
}

func (a Action) IsStandard() bool {
	return standardActions[a]
}

type PermissionCode struct {
	resource Resource
	action   Action
}

func NewPermissionCode(resource Resource, action Action) PermissionCode {
	return PermissionCode{
		resource: resource,
		action:   action,
	}
}

func ParsePermissionCode(code string) (PermissionCode, error) {
	parts := strings.Split(code, ":")
	if len(parts) != 2 {
		return PermissionCode{}, shared.NewValidationError("permission_code", "permission code must be in format 'resource:action'")
	}

	resource, err := NewResource(parts[0])
	if err != nil {
		return PermissionCode{}, fmt.Errorf("invalid resource in permission code: %w", err)
	}

	action, err := NewAction(parts[1])
	if err != nil {
		return PermissionCode{}, fmt.Errorf("invalid action in permission code: %w", err)
	}

	return PermissionCode{
		resource: resource,
		action:   action,
	}, nil
}

func (pc PermissionCode) String() string {
	return fmt.Sprintf("%s:%s", pc.resource.String(), pc.action.String())
}

func (pc PermissionCode) Resource() Resource {
	return pc.resource
}

func (pc PermissionCode) Action() Action {
	return pc.action
}

func (pc PermissionCode) Equals(other PermissionCode) bool {
	return pc.resource.Equals(other.resource) && pc.action == other.action
}
