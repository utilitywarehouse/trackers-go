package trackers

import (
	"context"

	"github.com/google/uuid"
)

type SchemaInfo interface {
	Name() string
	Version() int64
}

type Event interface {
	Payload() map[string]string
	Name() string
}

type Attribute interface {
	Name() string
	Value() interface{}
}

type Identity interface {
	Map() map[string]string
}

type Tracker interface {
	Track(
		ctx context.Context,
		schema SchemaInfo,
		identity Identity,
		events []Event,
		attributes []Attribute,
	) error
}

type BulkTracker interface {
	BulkTrack(ctx context.Context, batches []Batch) error
}

type Batch struct {
	Schema     SchemaInfo
	Identity   Identity
	Events     []Event
	Attributes []Attribute
}

var (
	accountNSUUID = uuid.NewSHA1(uuid.UUID{}, []byte("customer"))
	personNSUUID  = uuid.NewSHA1(uuid.UUID{}, []byte("person"))
)

func CustomerPersonIDFromAccountID(accountID string) string {
	return uuid.NewSHA1(personNSUUID, []byte(accountID+"-1")).String()
}

func CustomerPersonIDFromAccountNumber(accountNumber string) string {
	return CustomerPersonIDFromAccountID(hashAccountID(accountNumber))
}

func hashAccountID(accountNumber string) string {
	return uuid.NewSHA1(accountNSUUID, []byte(accountNumber)).String()
}
