package shared

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return Email{}, NewValidationError("email", "email cannot be empty")
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, NewValidationError("email", "invalid email format")
	}

	if !emailRegex.MatchString(value) {
		return Email{}, NewValidationError("email", "invalid email format")
	}

	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Value() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

func (e Email) IsZero() bool {
	return e.value == ""
}

type PhoneNumber struct {
	countryCode string
	number      string
}

var phoneRegex = regexp.MustCompile(`^\d{4,15}$`)

func NewPhoneNumber(countryCode, number string) (PhoneNumber, error) {
	countryCode = strings.TrimSpace(countryCode)
	number = strings.TrimSpace(number)

	cleanNumber := regexp.MustCompile(`\D`).ReplaceAllString(number, "")

	if countryCode == "" {
		return PhoneNumber{}, NewValidationError("phone", "country code cannot be empty")
	}

	if !strings.HasPrefix(countryCode, "+") {
		countryCode = "+" + countryCode
	}

	countryCodeDigits := strings.TrimPrefix(countryCode, "+")
	if len(countryCodeDigits) < 1 || len(countryCodeDigits) > 3 {
		return PhoneNumber{}, NewValidationError("phone", "invalid country code")
	}

	if cleanNumber == "" {
		return PhoneNumber{}, NewValidationError("phone", "phone number cannot be empty")
	}

	if !phoneRegex.MatchString(cleanNumber) {
		return PhoneNumber{}, NewValidationError("phone", "invalid phone number format")
	}

	return PhoneNumber{
		countryCode: countryCode,
		number:      cleanNumber,
	}, nil
}

func (p PhoneNumber) String() string {
	return p.countryCode + p.number
}

func (p PhoneNumber) CountryCode() string {
	return p.countryCode
}

func (p PhoneNumber) Number() string {
	return p.number
}

func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.countryCode == other.countryCode && p.number == other.number
}

func (p PhoneNumber) IsZero() bool {
	return p.countryCode == "" && p.number == ""
}

type Pagination struct {
	page  int
	limit int
}

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

func NewPagination(page, limit int) Pagination {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return Pagination{page: page, limit: limit}
}

func (p Pagination) Page() int {
	return p.page
}

func (p Pagination) Limit() int {
	return p.limit
}

func (p Pagination) Offset() int {
	return (p.page - 1) * p.limit
}

func (p Pagination) TotalPages(totalCount int64) int {
	if totalCount == 0 {
		return 0
	}
	pages := int(totalCount) / p.limit
	if int(totalCount)%p.limit > 0 {
		pages++
	}
	return pages
}

func (p Pagination) HasNext(totalCount int64) bool {
	return p.page < p.TotalPages(totalCount)
}

func (p Pagination) HasPrev() bool {
	return p.page > 1
}

type DateRange struct {
	from *string
	to   *string
}

func NewDateRange(from, to *string) DateRange {
	return DateRange{from: from, to: to}
}

func (d DateRange) From() *string {
	return d.from
}

func (d DateRange) To() *string {
	return d.to
}

func (d DateRange) HasFrom() bool {
	return d.from != nil && *d.from != ""
}

func (d DateRange) HasTo() bool {
	return d.to != nil && *d.to != ""
}

func (d DateRange) IsEmpty() bool {
	return !d.HasFrom() && !d.HasTo()
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

func (s SortOrder) String() string {
	return string(s)
}

func (s SortOrder) SQL() string {
	if s == SortOrderDesc {
		return "DESC"
	}
	return "ASC"
}

func NewSortOrder(s string) SortOrder {
	order := SortOrder(strings.ToLower(s))
	if !order.IsValid() {
		return SortOrderAsc
	}
	return order
}

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPaginatedResult[T any](items []T, total int64, pagination Pagination) PaginatedResult[T] {
	if items == nil {
		items = make([]T, 0)
	}
	return PaginatedResult[T]{
		Items:      items,
		Total:      total,
		Page:       pagination.Page(),
		Limit:      pagination.Limit(),
		TotalPages: pagination.TotalPages(total),
		HasNext:    pagination.HasNext(total),
		HasPrev:    pagination.HasPrev(),
	}
}

type FullName struct {
	value string
}

func NewFullName(value string) (FullName, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return FullName{}, NewValidationError("full_name", "full name cannot be empty")
	}
	if len(value) < 2 {
		return FullName{}, NewValidationError("full_name", "full name must be at least 2 characters")
	}
	if len(value) > 255 {
		return FullName{}, NewValidationError("full_name", "full name must not exceed 255 characters")
	}
	return FullName{value: value}, nil
}

func (f FullName) String() string {
	return f.value
}

func (f FullName) Value() string {
	return f.value
}

func (f FullName) Equals(other FullName) bool {
	return f.value == other.value
}

func (f FullName) IsZero() bool {
	return f.value == ""
}

type PasswordHash struct {
	value string
}

func NewPasswordHash(hashedValue string) (PasswordHash, error) {
	if hashedValue == "" {
		return PasswordHash{}, NewValidationError("password", "password hash cannot be empty")
	}
	return PasswordHash{value: hashedValue}, nil
}

func (p PasswordHash) String() string {
	return p.value
}

func (p PasswordHash) Value() string {
	return p.value
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("password", "password must be at least 8 characters")
	}
	if len(password) > 128 {
		return NewValidationError("password", "password must not exceed 128 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':\",./<>?", char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return NewValidationError("password", fmt.Sprintf("password must contain at least one %s", strings.Join(missing, ", ")))
	}

	return nil
}
