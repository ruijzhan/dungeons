package resolve

type Response struct {
	Host string `json:"host"`
	IPs  any    `json:"ips"`
}
