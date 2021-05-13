package go_example

import "context"

func main() {

	tracker := NewMParticleTracker("key", "secret", true)

	tracker.Track(
		context.Background(),
		schema.SchemaName,
		schema.SchemaVersion,
		&schema.Identity{CustomerPersonId: "abc"},
		&schema.HomeInsuranceQuoteAttemptedEvent{}
	)

}