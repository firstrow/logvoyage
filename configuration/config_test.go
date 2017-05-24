package configuration

import "testing"

func TestReadConfiguration(t *testing.T) {
	cfg := ReadConf("resources/logvoyage.yml")
	if cfg.Indexes.User != "logvoyage_user" {
		t.Errorf("Expect that Indexes.User to be eql: \"logvoyage_user\" and not: \"%s\"", cfg.Indexes.User)
	}
}
