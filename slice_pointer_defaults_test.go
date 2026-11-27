package env

import "testing"

func TestInitialSliceEntriesWithoutEnvironment(t *testing.T) {
	type server struct {
		Name string `env:"NAME,required"`
		Port int    `env:"PORT" envDefault:"8080"`
	}
	cfg := struct {
		Servers []server `envPrefix:"SERVERS_"`
	}{Servers: []server{{}}}
	opts := Options{Environment: map[string]string{}}

	err := ParseWithOptions(&cfg, opts)
	isErrorWithMessage(t, err, `env: required environment variable "SERVERS_0_NAME" is not set`)

	opts.Environment["SERVERS_0_NAME"] = "api"
	isNoErr(t, ParseWithOptions(&cfg, opts))
	isEqual(t, server{Name: "api", Port: 8080}, cfg.Servers[0])
}

func TestDefaultsInInitialSliceWithoutEnvironment(t *testing.T) {
	type server struct {
		Port int `env:"PORT" envDefault:"8080"`
	}
	servers := []server{{}, {Port: 9000}}
	cfg := struct {
		Servers *[]server `envPrefix:"SERVERS_"`
	}{Servers: &servers}
	opts := Options{Environment: map[string]string{}, SetDefaultsForZeroValuesOnly: true}

	isNoErr(t, ParseWithOptions(&cfg, opts))
	isEqual(t, []server{{Port: 8080}, {Port: 9000}}, *cfg.Servers)
}

func TestSlicePointerPreservesInitialValues(t *testing.T) {
	type item struct {
		Name string `env:"NAME"`
		Port int    `env:"PORT"`
	}
	items := []item{{Name: "first", Port: 8080}, {Name: "second", Port: 9090}}
	cfg := struct {
		Items *[]item `envPrefix:"ITEMS_"`
	}{Items: &items}
	err := ParseWithOptions(&cfg, Options{Environment: map[string]string{"ITEMS_0_NAME": "updated"}})
	isNoErr(t, err)
	isEqual(t, []item{{Name: "updated", Port: 8080}, {Name: "second", Port: 9090}}, *cfg.Items)
}
