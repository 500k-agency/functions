package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/500k-agency/function/data"
	"github.com/500k-agency/function/lib/connect"
	"github.com/500k-agency/function/lib/sendgrid"
	"github.com/stripe/stripe-go/v76"
)

var (
	ErrSessionUnpaid = errors.New("session unpaid")
)

func HandlePaymentCheckoutSession(ctx context.Context, session stripe.CheckoutSession) error {
	// if checkout is successful, send payment confirmation
	if session.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return ErrSessionUnpaid
	}

	name := data.SplitName(session.CustomerDetails.Name)
	var errs []error

	// fetch the checkout item list
	items := connect.StripeClient.GetSessionItems(session.ID)
	for _, it := range items {
		product := GetProductByID(it.Price.Product.ID)

		contact := &sendgrid.ContactRequest{
			ListIDs: product.PurchaseThankyou.ListIDs,
			Contacts: []*sendgrid.Contact{
				{
					Email:     session.CustomerDetails.Email,
					FirstName: name.FirstName,
					LastName:  name.LastName,
				},
			},
		}

		if err := connect.SendgridClient.AddContact(ctx, contact); err != nil {
			errs = append(errs, err)
		}

		err := connect.SendgridClient.Send(ctx, &sendgrid.MailRequest{
			Personalizations: []*sendgrid.MailPerson{
				{
					To: []*sendgrid.MailAddress{
						{Email: session.CustomerDetails.Email},
					},
					DynamicTemplateData: map[string]interface{}{
						"firstName":   name.FirstName,
						"productName": product.Name,
						"productUrl":  product.URL,
					},
				},
			},
			From: sendgrid.MailAddress{
				Email: "noreply@spacestationlabs.ltd",
				Name:  "Paul at Spacestation Labs",
			},
			ReplyTo: sendgrid.MailAddress{
				Email: "paul@spacestationlabs.ltd",
			},
			TemplateID:   product.PurchaseThankyou.TemplateID,
			MailSettings: &sendgrid.MailSettings{},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	var err error
	for _, v := range errs {
		err = fmt.Errorf("CheckoutSession: %+w", v)
		break
	}
	return err
}
