package mparticle

import (
	"context"
	"fmt"
	"github.com/mParticle/mparticle-go-sdk/events"
	"github.com/utilitywarehouse/trackers-go"
	"strconv"
)

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

func buildIdentity(mapped map[string]string) *events.UserIdentities {
	id := &events.UserIdentities{}

	for key, val := range mapped {
		switch key {
		case "Other":
			id.Other = val
			break
		case "CustomerID":
			id.CustomerID = val
			break
		case "Facebook":
			id.Facebook = val
			break
		case "Twitter":
			id.Twitter = val
			break
		case "Google":
			id.Google = val
			break
		case "Microsoft":
			id.Microsoft = val
			break
		case "Yahoo":
			id.Yahoo = val
			break
		case "Email":
			id.Email = val
			break
		case "Alias":
			id.Alias = val
			break
		case "FacebookCustomAudienceID":
			id.FacebookCustomAudienceID = val
			break
		case "OtherID2":
			id.OtherID2 = val
			break
		case "OtherID3":
			id.OtherID3 = val
			break
		case "OtherID4":
			id.OtherID4 = val
			break
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

	batch := events.Batch{Environment: t.Environment}

	batch.BatchContext = &events.BatchContext{
		DataPlan: &events.DataPlanContext{
			PlanID:      schema.Name(),
			PlanVersion: schema.Version(),
		},
	}

	batch.UserIdentities = buildIdentity(identity.Map())

	batch.UserAttributes = make(map[string]interface{})
	batch.Events = []events.Event{}

	for _, x := range evs {
		customEvent := events.NewCustomEvent()
		customEvent.Data.EventName = x.Name()
		customEvent.Data.CustomEventType = events.OtherCustomEventType
		customEvent.Data.CustomAttributes = x.Payload()

		//attach some event ID data as custom flags
		customEvent.Data.CustomAttributes["uw.schema-name"] = schema.Name()
		customEvent.Data.CustomAttributes["uw.schema-version"] = strconv.FormatInt(schema.Version(), 10)

		batch.Events = append(batch.Events, customEvent)
	}

	for _, x := range attribs {
		batch.UserAttributes[x.Name()] = x.Value()
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

	return nil
}
