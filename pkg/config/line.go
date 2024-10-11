package config

type LineConfig struct {
	ClientID      string `yaml:"channel_id"`
	ChannelSecret string `yaml:"channel_secret"`
	ChannelToken  string `yaml:"channel_access_token"`
}
