package asa

type Config struct {
	Enabled          bool `config:"enabled"`
	IncludeTimestamp bool `config:"include_timestamp"`
}
