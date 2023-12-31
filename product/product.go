package product

type Product struct {
	Config
}

// Config holds all the configuration fields needed within the application
type Config struct {
	Name             string      `toml:"name"`
	StripeID         string      `toml:"stripe_id"`
	URL              string      `toml:"url"`
	PurchaseThankyou EmailConfig `toml:"purchase_thankyou"`
}

type EmailConfig struct {
	ListIDs    []string `toml:"list_ids"`
	TemplateID string   `toml:"template_id"`
}

var (
	productCatalogue = map[string]Product{}
)

func Setup(confs []Config) {
	for _, v := range confs {
		productCatalogue[v.StripeID] = Product{
			Config: v,
		}
	}
}

func GetProductByID(productId string) Product {
	return productCatalogue[productId]
}
