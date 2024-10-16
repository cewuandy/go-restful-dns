package domain

type Options struct {
	HttpAddr           string `default:"0.0.0.0" usage:"[Server Mode] Restful API address"`
	HttpPort           uint   `default:"8081" usage:"[Server Mode] Restful API port"`
	DnsAddr            string `default:"0.0.0.0" usage:"DNS address"`
	DnsPort            uint   `default:"53" usage:"DNS port"`
	UpstreamForwarders string `default:"1.1.1.1:53" usage:"DNS upstream forwarders, e.g. 1.1.1.1:53,8.8.8.8:53"`
	RedisAddr          string `default:"" usage:"Redis address"`
	RedisPassword      string `default:"" usage:"Redis password"`
	RedisMasterName    string `default:"" usage:"Redis master"`
	RedisSentinelHost  string `default:"" usage:"Sentinel host"`
	RedisSentinelPort  uint   `default:"" usage:"Sentinel port"`
}
