package vpcflow

type config struct {
	Type string `config:"type"`
}

func defaultConfig() config {
	return config{
		Type: "aws:vpcflow",
	}
}
