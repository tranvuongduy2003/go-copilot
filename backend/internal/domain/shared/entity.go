package shared

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	id uuid.UUID
}

func NewEntity() Entity {
	return Entity{id: uuid.New()}
}

func NewEntityWithID(id uuid.UUID) Entity {
	return Entity{id: id}
}

func (e Entity) ID() uuid.UUID {
	return e.id
}

func (e Entity) Equals(other Entity) bool {
	return e.id == other.id
}

func (e Entity) IsZero() bool {
	return e.id == uuid.Nil
}

type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() uuid.UUID
}

type BaseDomainEvent struct {
	aggregateID uuid.UUID
	occurredAt  time.Time
	eventType   string
}

func NewBaseDomainEvent(aggregateID uuid.UUID, eventType string) BaseDomainEvent {
	return BaseDomainEvent{
		aggregateID: aggregateID,
		occurredAt:  time.Now().UTC(),
		eventType:   eventType,
	}
}

func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseDomainEvent) AggregateID() uuid.UUID {
	return e.aggregateID
}

type AggregateRoot struct {
	Entity
	domainEvents []DomainEvent
}

func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{
		Entity:       NewEntity(),
		domainEvents: make([]DomainEvent, 0),
	}
}

func NewAggregateRootWithID(id uuid.UUID) AggregateRoot {
	return AggregateRoot{
		Entity:       NewEntityWithID(id),
		domainEvents: make([]DomainEvent, 0),
	}
}

func (ar *AggregateRoot) AddDomainEvent(event DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, event)
}

func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}

func (ar *AggregateRoot) ClearDomainEvents() {
	ar.domainEvents = make([]DomainEvent, 0)
}

func (ar *AggregateRoot) PopDomainEvents() []DomainEvent {
	events := ar.domainEvents
	ar.ClearDomainEvents()
	return events
}
