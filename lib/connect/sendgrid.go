package connect

import (
	"context"

	"github.com/500k-agency/function/lib/sendgrid"
)

type Sendgrid struct {
	Client  *sendgrid.Client
	Sandbox bool
}

type SendgridConfig struct {
	Config
	Sandbox bool `toml:"sandbox" env:"SANDBOX"`
}

var SendgridClient *Sendgrid

func SetupSendgrid(conf SendgridConfig) *Sendgrid {
	client, _ := sendgrid.NewClient(
		nil,
		sendgrid.WithApp(conf.AppID, conf.AppSecret),
		sendgrid.WithDebug(true),
	)
	SendgridClient = &Sendgrid{
		Client:  client,
		Sandbox: conf.Sandbox,
	}
	return SendgridClient
}

func (s *Sendgrid) AddContact(ctx context.Context, v *sendgrid.ContactRequest) error {
	if s.Sandbox {
		return nil
	}
	if _, _, err := s.Client.Contact.Upsert(ctx, v); err != nil {
		return err
	}
	return nil
}

func (s *Sendgrid) Send(ctx context.Context, v *sendgrid.MailRequest) error {
	if s.Sandbox {
		v.MailSettings.SandboxMode = sendgrid.NewSetting(true)
	}
	if _, err := s.Client.Mail.Send(ctx, v); err != nil {
		return err
	}
	return nil
}
