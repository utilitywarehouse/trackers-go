package mparticle

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mParticle/mparticle-go-sdk/events"

	"github.com/utilitywarehouse/trackers-go"
)

var _ trackers.Tracker = (*MParticleTracker)(nil)
var _ trackers.BulkTracker = (*MParticleTracker)(nil)

type MParticleTracker struct {
	environment events.Environment
	auth        events.BasicAuth
	client      *events.APIClient
}

func NewMParticleTracker(url, apiKey, apiSecret string, client *http.Client, isDev bool) *MParticleTracker {
	env := events.ProductionEnvironment
	if isDev {
		env = events.DevelopmentEnvironment
	}
	cfg := events.NewConfiguration()
	cfg.BasePath = url
	cfg.HTTPClient = client
	return &MParticleTracker{
		environment: env,
		client:      events.NewAPIClient(cfg),
		auth: events.BasicAuth{
			APIKey:    apiKey,
			APISecret: apiSecret,
		},
	}
}

func buildIdentity(mapped map[string]string) *events.UserIdentities {
	id := &events.UserIdentities{}

	for key, val := range mapped {
		switch key {
		case "Other":
			id.Other = val
		case "CustomerID":
			id.CustomerID = val
		case "Facebook":
			id.Facebook = val
		case "Twitter":
			id.Twitter = val
		case "Google":
			id.Google = val
		case "Microsoft":
			id.Microsoft = val
		case "Yahoo":
			id.Yahoo = val
		case "Email":
			id.Email = val
		case "Alias":
			id.Alias = val
		case "FacebookCustomAudienceID":
			id.FacebookCustomAudienceID = val
		case "Other2":
			id.OtherID2 = val
		case "Other3":
			id.OtherID3 = val
		case "Other4":
			id.OtherID4 = val
		}
	}

	return id
}

func (t *MParticleTracker) Track(
	ctx context.Context,
	schema trackers.SchemaInfo,
	identity trackers.Identity,
	evs []trackers.Event,
	attribs []trackers.Attribute,
) error {

	calLCtx := context.WithValue(
		ctx,
		events.ContextBasicAuth,
		t.auth,
	)
	batch := toMParticleBatch(schema, identity, evs, attribs, t.environment)
	result, err := t.client.EventsAPI.UploadEvents(calLCtx, batch)
	if result == nil || result.StatusCode != 202 {
		if gerr, ok := err.(events.GenericError); ok {
			return fmt.Errorf(
				"Error while uploading!\nstatus: %v\nresponse body: %#v",
				gerr.Error(),
				gerr.Model(),
			)
		}
		return fmt.Errorf("Unexpected error while uploading!\nerror:%#v", err)
	}

	return nil
}

func (t *MParticleTracker) BulkTrack(ctx context.Context, batches []trackers.Batch) error {

	calLCtx := context.WithValue(
		ctx,
		events.ContextBasicAuth,
		t.auth,
	)
	mpBatches := []events.Batch{}
	for _, b := range batches {
		mpBatches = append(mpBatches, toMParticleBatch(b.Schema, b.Identity, b.Events, b.Attributes, t.environment))
	}
	if len(mpBatches) == 0 {
		return nil
	}
	result, err := t.client.EventsAPI.BulkUploadEvents(calLCtx, mpBatches)
	if result == nil || result.StatusCode != 202 {
		if gerr, ok := err.(events.GenericError); ok {
			return fmt.Errorf(
				"Error while uploading!\nstatus: %v\nresponse body: %#v",
				gerr.Error(),
				gerr.Model(),
			)
		}
		return fmt.Errorf("Unexpected error while uploading!\nerror:%#v", err)
	}

	return nil
}

func toMParticleBatch(schema trackers.SchemaInfo,
	identity trackers.Identity,
	evs []trackers.Event,
	attribs []trackers.Attribute,
	env events.Environment) events.Batch {
	batch := events.Batch{
		Environment:    env,
		UserIdentities: buildIdentity(identity.Map()),
		UserAttributes: make(map[string]interface{}),
		Events:         []events.Event{},
	}
	if schema != trackers.NoSchema {
		batch.BatchContext = &events.BatchContext{
			DataPlan: &events.DataPlanContext{
				PlanID:      schema.Name(),
				PlanVersion: schema.Version(),
			},
		}
	}
	for _, x := range evs {
		customEvent := events.NewCustomEvent()
		customEvent.Data.EventName = x.Name()
		customEvent.Data.CustomEventType = events.OtherCustomEventType
		customEvent.Data.CustomAttributes = x.Payload()

		if customEvent.Data.CustomAttributes["uw.schema-name"] == "" {
			customEvent.Data.CustomAttributes["uw.schema-name"] = schema.Name()
		}
		if customEvent.Data.CustomAttributes["uw.schema-version"] == "" {
			customEvent.Data.CustomAttributes["uw.schema-version"] = strconv.FormatInt(schema.Version(), 10)
		}

		batch.Events = append(batch.Events, customEvent)
	}

	for _, x := range attribs {
		batch.UserAttributes[x.Name()] = x.Value()
	}
	return batch
}
