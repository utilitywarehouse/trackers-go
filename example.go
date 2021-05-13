package main

import (
	"context"
	"github.com/utilitywarehouse/tracker-example/schema"
	"log"
)

func main() {

	tracker := NewMParticleTracker("key", "secret", true)

	err := tracker.Track(
		context.Background(),
		schema.SchemaName,
		schema.SchemaVersion,
		&schema.Identity{CustomerPersonId: "abc"},
		&schema.HomeInsuranceQuoteAttemptedEvent{
			QuoteId:        "abc",
			QuoteReference: "fef",
			Product:        schema.ProductContents,
			QuoteType:      schema.QuoteTypeRenewal,
		},
		schema.HomeInsuranceRenewalDateAttribute("2016-01-21"),
	)

	if err != nil {
		log.Fatalln(err)
	}

}
