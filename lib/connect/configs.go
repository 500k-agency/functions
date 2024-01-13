package connect

type Config struct {
	AppID         string `toml:"app_id"`
	AppSecret     string `toml:"app_secret"`
	WebhookSecret string `toml:"webhook_secret"`
	ReturnURL     string `toml:"return_url"`
}

type Configs struct {
	Stripe   Config         `toml:"stripe"`
	Tally    Config         `toml:"tally"`
	Sendgrid SendgridConfig `tomp:"sendgrid"`
}

// Configure loads the connect configs from config file
func Configure(confs Configs) {
	SetupStripe(confs.Stripe)
	SetupSendgrid(confs.Sendgrid)
	SetupTally(confs.Tally)
}
