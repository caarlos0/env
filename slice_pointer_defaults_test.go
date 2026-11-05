package env

import "testing"

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
