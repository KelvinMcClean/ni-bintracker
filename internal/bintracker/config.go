package bintracker

type Config struct {
	House struct {
		Id       int64  `toml:"id"`
		Postcode string `toml:"postcode"`
		Number   int64  `toml:"number"`
	}
	Council struct {
		Name string `toml:"name"`
	}
	Calendar struct {
		ID string `toml:"id"`
	}
}
