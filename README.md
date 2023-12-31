# 500k.agency toolkit

## Thank you email:

Sets up a cloud function as webhook to send thank you email on stripe
  purchase.

### Sendgrid

1. Create an API key with the following permissions

- Mail send (full access)
- Marketing (full access)

2. Upload templates from the `templates/` directory


### Cloudfunction

1. create a user managed service account
    https://cloud.google.com/iam/docs/service-accounts-create#iam-service-accounts-create-console

2. from secrets manager, give permission to the new service function.

3. Deploy the function with `make deploy`

### Stripe

1. Register the webhook
2. Update config with webhook secret
3. Redeploy with `make deploy`

### Testing

1. Update function.conf
2. Run cloud function via `make run`
3. `stripe listen --forward-to localhost:8080`
4. `stripe trigger checkout.session.completed` to generate mock products
5. Run command `export PRODUCT_PRICE=<price_id>; export CUSTOMER_EMAIL=<customer_email>; stripe fixtures fixtures/checkout.session.completed.json`
