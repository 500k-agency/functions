package connect

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"github.com/stripe/stripe-go/v76/webhook"
)

// Stripe config struct with exposed methods needed
type Stripe struct {
	*client.API

	client *client.API
	config Config
}

var (
	StripeClient *Stripe
)

// SetupStripe sets up stripe with the credentials given
func SetupStripe(confs Config) *Stripe {
	sc := &client.API{}
	sc.Init(confs.AppSecret, &stripe.Backends{
		API: stripe.GetBackendWithConfig(
			stripe.APIBackend,
			&stripe.BackendConfig{
				LeveledLogger: &stripe.LeveledLogger{
					Level: stripe.LevelInfo,
				},
			},
		),
		Connect: stripe.GetBackend(stripe.ConnectBackend),
		Uploads: stripe.GetBackend(stripe.UploadsBackend),
	})
	StripeClient = &Stripe{
		API:    sc,
		client: sc,
		config: confs,
	}
	return StripeClient
}

// CreateCustomerBillingPortal creates customer portal to auto manage billing and subscriptions
func (s *Stripe) CreateCustomerBillingPortal(customerID, returnUrl string) (*stripe.BillingPortalSession, error) {
	if returnUrl == "" {
		returnUrl = s.config.ReturnURL
	}
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(returnUrl),
	}
	return s.client.BillingPortalSessions.New(params)
}

// ConstructEvent validates stripe webhook secret is authentic
func (s *Stripe) ConstructEvent(body []byte, header string) (stripe.Event, error) {
	return webhook.ConstructEventWithOptions(
		body,
		header,
		s.config.WebhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
}

// VerifySession checks if the stripe session id is valid
func (s *Stripe) VerifySession(sessionID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{}
	return s.client.CheckoutSessions.Get(sessionID, params)
}

// Sessio
func (s *Stripe) GetSessionItems(sessionID string) []*stripe.LineItem {
	params := &stripe.CheckoutSessionListLineItemsParams{
		Session: stripe.String(sessionID),
	}

	var items []*stripe.LineItem
	iter := s.client.CheckoutSessions.ListLineItems(params)
	for iter.Next() {
		items = append(items, iter.LineItem())
	}
	return items
}
