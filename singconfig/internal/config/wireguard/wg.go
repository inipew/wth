package wireguard

type WireGuardOutbound struct{
	Type            string   `json:"type,omitempty"`
	Tag             string   `json:"tag,omitempty"`
	DomainStrategy            string     `json:"domain_strategy,omitempty"`
	LocalAddress              []string   `json:"local_address,omitempty"`
	PrivateKey                string     `json:"private_key,omitempty"`
	Server                    string     `json:"server,omitempty"`
	ServerPort                int        `json:"server_port,omitempty"`
	PeerPublicKey             string     `json:"peer_public_key,omitempty"`
	MTU                       int        `json:"mtu,omitempty"`
	URL                       string     `json:"url,omitempty"`
}