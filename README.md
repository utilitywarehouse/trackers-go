# trackers-go

Go implementations of event tracking.

## Mparticle Usage 
You will need:
- a configured [input feed](https://docs.mparticle.com/guides/getting-started/create-an-input/) for mparticle for your api key and secret.
- the url which in most cases is `https://s2s.mparticle.com/v2`
- your types [generated](https://github.com/utilitywarehouse/tracking-types-gen/) from your [data plan](https://github.com/utilitywarehouse/analytics-contracts).

First install using `go get github.com/utilitywarehouse/trackers-go`

Then you need to import the tracker and use your configured mparticle input feed. 
```go
httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	tracker := mparticle.NewMParticleTracker("url", "key", "secret", httpClient, true)
```

Then an event can be tracked by calling the `track` method.
```go
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
```
