package config

type Config struct {
	Client Client     `json:"client"`
	Server Server     `json:"server"`
	Result Result     `json:"result"`
	Host   HostMgConf `json:"host"`
}