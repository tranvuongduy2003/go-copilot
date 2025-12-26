package validator

import (
	"net"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Validator struct {
	validate *validator.Validate
	messages map[string]string
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Field+": "+err.Message)
	}
	return "validation failed:\n  - " + strings.Join(msgs, "\n  - ")
}

var (
	instance *Validator
	once     sync.Once
)

func New() *Validator {
	v := &Validator{
		validate: validator.New(),
		messages: defaultMessages(),
	}

	v.registerCustomValidations()
	return v
}

func Get() *Validator {
	once.Do(func() {
		instance = New()
	})
	return instance
}

func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return v.translateErrors(validationErrs)
}

func (v *Validator) ValidateVar(field interface{}, tag string) error {
	err := v.validate.Var(field, tag)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return v.translateErrors(validationErrs)
}

func (v *Validator) RegisterMessage(tag string, message string) {
	v.messages[tag] = message
}

func (v *Validator) translateErrors(errs validator.ValidationErrors) ValidationErrors {
	var result ValidationErrors

	for _, err := range errs {
		field := toSnakeCase(err.Field())
		tag := err.Tag()
		param := err.Param()

		message := v.getMessage(tag, field, param)

		result = append(result, ValidationError{
			Field:   field,
			Tag:     tag,
			Value:   formatValue(err.Value()),
			Message: message,
		})
	}

	return result
}

func (v *Validator) getMessage(tag, field, param string) string {
	if msg, ok := v.messages[tag]; ok {
		msg = strings.ReplaceAll(msg, "{field}", field)
		msg = strings.ReplaceAll(msg, "{param}", param)
		return msg
	}
	return field + " failed on " + tag + " validation"
}

func (v *Validator) registerCustomValidations() {
	v.validate.RegisterValidation("email_dns", validateEmailWithDNS)
	v.validate.RegisterValidation("phone", validatePhone)
	v.validate.RegisterValidation("password", validatePassword)
	v.validate.RegisterValidation("uuid_format", validateUUID)
	v.validate.RegisterValidation("username", validateUsername)
}

func validateEmailWithDNS(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	if domain == "" {
		return false
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		_, err = net.LookupIP(domain)
		return err == nil
	}

	return true
}

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}

	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	return phoneRegex.MatchString(cleaned)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true
	}

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateUUID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	_, err := uuid.Parse(value)
	return err == nil
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,29}$`)

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true
	}

	return usernameRegex.MatchString(username)
}

func defaultMessages() map[string]string {
	return map[string]string{
		"required":    "{field} is required",
		"email":       "{field} must be a valid email address",
		"email_dns":   "{field} must be a valid email address with an existing domain",
		"phone":       "{field} must be a valid phone number",
		"password":    "{field} must be at least 8 characters with uppercase, lowercase, number, and special character",
		"uuid_format": "{field} must be a valid UUID",
		"username":    "{field} must start with a letter and contain only letters, numbers, and underscores (3-30 characters)",
		"min":         "{field} must be at least {param} characters",
		"max":         "{field} must be at most {param} characters",
		"len":         "{field} must be exactly {param} characters",
		"gte":         "{field} must be greater than or equal to {param}",
		"lte":         "{field} must be less than or equal to {param}",
		"gt":          "{field} must be greater than {param}",
		"lt":          "{field} must be less than {param}",
		"eq":          "{field} must be equal to {param}",
		"ne":          "{field} must not be equal to {param}",
		"oneof":       "{field} must be one of: {param}",
		"url":         "{field} must be a valid URL",
		"uri":         "{field} must be a valid URI",
		"alpha":       "{field} must contain only letters",
		"alphanum":    "{field} must contain only letters and numbers",
		"numeric":     "{field} must be a valid number",
		"boolean":     "{field} must be a boolean",
		"contains":    "{field} must contain '{param}'",
		"excludes":    "{field} must not contain '{param}'",
		"startswith":  "{field} must start with '{param}'",
		"endswith":    "{field} must end with '{param}'",
		"datetime":    "{field} must be a valid datetime in format {param}",
		"json":        "{field} must be valid JSON",
		"jwt":         "{field} must be a valid JWT",
		"uuid":        "{field} must be a valid UUID",
		"uuid4":       "{field} must be a valid UUID v4",
		"ip":          "{field} must be a valid IP address",
		"ipv4":        "{field} must be a valid IPv4 address",
		"ipv6":        "{field} must be a valid IPv6 address",
		"cidr":        "{field} must be a valid CIDR notation",
		"mac":         "{field} must be a valid MAC address",
		"hostname":    "{field} must be a valid hostname",
		"fqdn":        "{field} must be a valid FQDN",
		"unique":      "{field} must contain unique values",
		"ascii":       "{field} must contain only ASCII characters",
		"printascii":  "{field} must contain only printable ASCII characters",
		"base64":      "{field} must be valid base64",
		"hexadecimal": "{field} must be a valid hexadecimal",
		"lowercase":   "{field} must be lowercase",
		"uppercase":   "{field} must be uppercase",
		"eqfield":     "{field} must be equal to {param}",
		"nefield":     "{field} must not be equal to {param}",
		"gtfield":     "{field} must be greater than {param}",
		"gtefield":    "{field} must be greater than or equal to {param}",
		"ltfield":     "{field} must be less than {param}",
		"ltefield":    "{field} must be less than or equal to {param}",
	}
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		if len(val) > 50 {
			return val[:50] + "..."
		}
		return val
	default:
		return ""
	}
}

func IsValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}

func GetValidationErrors(err error) (ValidationErrors, bool) {
	errs, ok := err.(ValidationErrors)
	return errs, ok
}
