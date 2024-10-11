package config

type LineConfig struct {
	ClientID      string `yaml:"channel_id"`
	ChannelSecret string `yaml:"channel_secret"`
	AccessToken   string `yaml:"channel_access_token"`
}
