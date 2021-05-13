package go_example

import (
	"context"
	"errors"
	"fmt"
	"github.com/mParticle/mparticle-go-sdk/events"
	"github.com/utilitywarehouse/mparticle-data-plan/go-example/schema"
)

type MParticleTracker struct {
	Environment events.Environment
	APIKey string
	APISecret string
	Client *events.APIClient
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
	schemaName string,
	schemaVersion string,
	identity schema.Identity,
	payloads... interface{},
) error {

	mappedID := identity.Map()

	batch := events.Batch{Environment: t.Environment}

	batch.BatchContext = &events.BatchContext{
		DataPlan: &events.DataPlanContext{
			PlanID:      schemaName,
			PlanVersion: schemaVersion,
		},
	}

	batch.UserIdentities = &events.UserIdentities{
		OtherID4: mappedID.Other4,
		Email: mappedID.Email,
	}
	batch.UserAttributes = make(map[string]interface{})
	batch.Events = []events.Event{}

	for _, p := range payloads {
		switch x := p.(type) {
		case schema.Event:
			customEvent := events.NewCustomEvent()
			customEvent.Data.EventName = x.Name()
			customEvent.Data.CustomEventType = events.OtherCustomEventType
			customEvent.Data.CustomAttributes = x.Payload()
			batch.Events = append(batch.Events, customEvent)
			break
		case schema.Attribute:
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

	result, err := t.Client.EventsAPI.UploadEvents(calLCtx, batch)

	if result == nil || result.StatusCode != 202 {
		return fmt.Errorf(
			"Error while uploading!\nstatus: %v\nresponse body: %#v",
			err.(events.GenericError).Error(),
			err.(events.GenericError).Model(),
		)
	}

}

