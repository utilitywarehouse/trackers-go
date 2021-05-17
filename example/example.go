package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/utilitywarehouse/trackers-go"
	"github.com/utilitywarehouse/trackers-go/example/schema"
	"github.com/utilitywarehouse/trackers-go/mparticle"
)

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	tracker := mparticle.NewMParticleTracker("url", "key", "secret", httpClient, true)

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
