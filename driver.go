package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mParticle/mparticle-go-sdk/events"
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

type MParticleTracker struct {
	Environment events.Environment
	APIKey      string
	APISecret   string
	Client      *events.APIClient
}

func NewMParticleTracker(APIKey, APISecret string, isDev bool) *MParticleTracker {

	env := events.ProductionEnvironment

	if isDev {
		env = events.DevelopmentEnvironment
	}

	client := events.NewAPIClient(events.NewConfiguration())

	return &MParticleTracker{
		Environment: env,
		APIKey:      APIKey,
		APISecret:   APISecret,
		Client:      client,
	}
}

func (t *MParticleTracker) Track(
	ctx context.Context,
	schema SchemaInfo,
	identity Identity,
	payloads ...interface{},
) error {

	batch := events.Batch{Environment: t.Environment}

	batch.BatchContext = &events.BatchContext{
		DataPlan: &events.DataPlanContext{
			PlanID:      schema.Name(),
			PlanVersion: schema.Version(),
		},
	}

	batch.UserIdentities = &events.UserIdentities{}

	for key, val := range identity.Map() {
		switch key {
		case "OtherID4":
			batch.UserIdentities.OtherID4 = val
			break
		case "Email":
			batch.UserIdentities.Email = val
			break
		}
	}

	batch.UserAttributes = make(map[string]interface{})
	batch.Events = []events.Event{}

	for _, p := range payloads {
		switch x := p.(type) {
		case Event:
			customEvent := events.NewCustomEvent()
			customEvent.Data.EventName = x.Name()
			customEvent.Data.CustomEventType = events.OtherCustomEventType
			customEvent.Data.CustomAttributes = x.Payload()
			batch.Events = append(batch.Events, customEvent)
			break
		case Attribute:
			batch.UserAttributes[x.Name()] = x.Value()
		default:
			return errors.New("could not convert payloads into either Event or Attribute")
		}
	}

	calLCtx := context.WithValue(
		ctx,
		events.ContextBasicAuth,
		events.BasicAuth{
			APIKey:    t.APIKey,
			APISecret: t.APISecret,
		},
	)

	spew.Dump(batch)

	result, err := t.Client.EventsAPI.UploadEvents(calLCtx, batch)

	if result == nil || result.StatusCode != 202 {
		return fmt.Errorf(
			"Error while uploading!\nstatus: %v\nresponse body: %#v",
			err.(events.GenericError).Error(),
			err.(events.GenericError).Model(),
		)
	}

	return nil
}
