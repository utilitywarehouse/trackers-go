package main

import (
	"context"
	"github.com/utilitywarehouse/trackers-go"
	"github.com/utilitywarehouse/trackers-go/example/schema"
	"github.com/utilitywarehouse/trackers-go/mparticle"
	"log"
)

func main() {

	tracker := mparticle.NewMParticleTracker("key", "secret", true)

	err := tracker.Track(
		context.Background(),
		schema.Info,
		&schema.Identity{CustomerPersonId: "abc"},
		[]trackers.Event{
			&schema.HomeInsuranceQuoteAttemptedEvent{
				QuoteId:        "abc",
				QuoteReference: "fef",
				Product:        schema.ProductContents,
				QuoteType:      schema.QuoteTypeRenewal,
			},
		},
		[]trackers.Attribute{
			schema.HomeInsuranceRenewalDateAttribute("2016-01-21"),
		},
	)

	if err != nil {
		log.Fatalln(err)
	}

}
