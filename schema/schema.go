package schema

const SchemaName = "insurance"
const SchemaVersion = 1

type Event interface {
	Payload() map[string]string
	Name() string
}

type Attribute interface {
	Name() string
	Value() interface{}
}

type Product string
const (
	ProductBuildings         Product = "ProductBuildings"
	ProductContents          Product = "ProductContents"
	ProductBuildingsContents Product = "ProductBuildingsContents"
)


type QuoteType string
const (
	QuoteTypeNewBusiness       QuoteType = "QuoteTypeNewBusiness"
	QuoteTypeRenewal           QuoteType = "QuoteTypeRenewal"
	QuoteTypeMidTermAdjustment QuoteType = "QuoteTypeMidTermAdjustment"
)

type HomeInsuranceQuoteAttemptedEvent struct {
	QuoteId        string    `json:"quote_id,omitempty"`
	QuoteReference string    `json:"quote_reference,omitempty"`
	Product        Product   `json:"product,omitempty"`
	QuoteType      QuoteType `json:"quote_type,omitempty"`
}

func (e *HomeInsuranceQuoteAttemptedEvent) Name() string {
	return "home-insurance-quote-attempted-event"
}

func (e *HomeInsuranceQuoteAttemptedEvent) Payload() map[string]string {
	return map[string]string {
		"QuoteId": e.QuoteId,
		"QuoteReference": e.QuoteReference,
		"Product": string(e.Product),
		"QuoteType": string(e.QuoteType),
	}
}

type HomeInsuranceRenewalDateAttribute struct {
	attribute string
}

func (a *HomeInsuranceRenewalDateAttribute) Name() string {
	return "insurance.home.renewalAt"
}

func (a *HomeInsuranceRenewalDateAttribute) Value() interface{} {
	return a.attribute
}

type Identity struct {
	CustomerPersonId string `json:"customerPersonId,omitempty"` // a UUID of a person, generated as UUIDv5 of account_number until we resolve the people ID problem
	Email string `json:"email,omitempty"` // email of the person logging in
}

type MappedIdentity struct {
	Other4 string `json:"Other4,omitempty"` // a UUID of a person, generated as UUIDv5 of account_number until we resolve the people ID problem
	Email string `json:"Email,omitempty"` // email of the person logging in
}

// Map - returns mapped identity
func (o *Identity) Map() *MappedIdentity {
	return &MappedIdentity{
		Other4: o.CustomerPersonId,
		Email: o.Email,
	}
}
