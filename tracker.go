package trackers

import "context"

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
