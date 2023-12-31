package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/500k-agency/function/api"
	"github.com/500k-agency/function/config"
	"github.com/500k-agency/function/lib/connect"
	"github.com/500k-agency/function/product"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-chi/render"
	"github.com/stripe/stripe-go/v76"

	_ "github.com/500k-agency/function/product"
)

const (
	maxStripeBodyBytes = int64(65536)
)

func setup() {
	conf, err := config.NewFromSecrets()
	if err != nil {
		log.Fatalf("main.NewFromSecrets: %v\n", err)
	}
	connect.Configure(conf.Connect)
	product.Setup(conf.Products)
}

func init() {
	functions.HTTP("PurchaseHandler", PurchaseHandler)
}

// PurchaseHandler handle incoming stripe connections
func PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	setup()

	r.Body = http.MaxBytesReader(w, r.Body, maxStripeBodyBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		render.Respond(w, r, err)
		return
	}
	defer r.Body.Close()

	// Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	event, err := connect.StripeClient.ConstructEvent(body, r.Header.Get("Stripe-Signature"))
	// Ignore Signature for now.
	if err != nil {
		render.Respond(w, r, api.ErrInvalidRequest(fmt.Errorf("Stripe ConstructEvent errored: %w", err)))
		return
	}

	ctx := context.WithValue(r.Context(), &api.ContextKey{Name: "eventType"}, event.Type)

	switch event.Type {
	case "checkout.session.completed":
		// Sent when a customer clicks the Pay or Subscribe button in Checkout, informing you of a new purchase.
		// Payment is successful and the subscription is created. Provision the subscription.
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			render.Respond(w, r, fmt.Sprintf("CheckoutSessionCompleted handler errored: %+v", err))
			return
		}
		switch session.Mode {
		case stripe.CheckoutSessionModePayment:
			if err := product.HandlePaymentCheckoutSession(ctx, session); err != nil {
				render.Respond(w, r, fmt.Sprintf("CheckoutSessionCompleted CheckoutSessionModePayment handler errored: %+v", err))
				return
			}
			// ignore other modes
		case stripe.CheckoutSessionModeSubscription:
			break
		case stripe.CheckoutSessionModeSetup:
			break
		}
	}

	// Send an HTTP response
	render.Respond(w, r, "OK")
}
