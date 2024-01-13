package waitlist

type Waitlist struct {
	Config
}

// Config holds all the configuration fields needed within the application
type Config struct {
	Name    string   `toml:"name"`
	FormID  string   `toml:"form_id"` // tally form id
	ListIDs []string `toml:"list_ids"`
}

var (
	waitlist *Waitlist
)

func Setup(conf Config) {
	waitlist = &Waitlist{
		Config: conf,
	}
}

func GetListByID(listId string) *Waitlist {
	return waitlist
}
