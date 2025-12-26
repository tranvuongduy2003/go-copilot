package user

type Status string

const (
	StatusPending  Status = "pending"
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBanned   Status = "banned"
)

var validStatuses = map[Status]bool{
	StatusPending:  true,
	StatusActive:   true,
	StatusInactive: true,
	StatusBanned:   true,
}

var allowedTransitions = map[Status][]Status{
	StatusPending:  {StatusActive, StatusBanned},
	StatusActive:   {StatusInactive, StatusBanned},
	StatusInactive: {StatusActive, StatusBanned},
	StatusBanned:   {},
}

func (s Status) IsValid() bool {
	return validStatuses[s]
}

func (s Status) String() string {
	return string(s)
}

func (s Status) CanTransitionTo(target Status) bool {
	allowed, exists := allowedTransitions[s]
	if !exists {
		return false
	}
	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

func (s Status) IsPending() bool {
	return s == StatusPending
}

func (s Status) IsActive() bool {
	return s == StatusActive
}

func (s Status) IsInactive() bool {
	return s == StatusInactive
}

func (s Status) IsBanned() bool {
	return s == StatusBanned
}

func ParseStatus(s string) (Status, bool) {
	status := Status(s)
	return status, status.IsValid()
}
