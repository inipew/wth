package wireguard

type WireGuardOutbound struct {
    Type            string   `json:"type,omitempty"`
    Tag             string   `json:"tag,omitempty"`
    DomainStrategy  string   `json:"domain_strategy,omitempty"`
    LocalAddress    []string `json:"local_address,omitempty"`
    PrivateKey      string   `json:"private_key,omitempty"`
    Server          string   `json:"server,omitempty"`
    ServerPort      int      `json:"server_port,omitempty"`
    PeerPublicKey   string   `json:"peer_public_key,omitempty"`
    MTU             int      `json:"mtu,omitempty"`
}

func BuildNewWireGuardOutbound(ipv6, privateKey string) *WireGuardOutbound {
    return newWireguardOutbound(
        "warp-out",
        "",
        []string{"172.16.0.2/32", ipv6},
        "engage.cloudflareclient.com",
        2048,
        privateKey,
        "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
    )
}

func newWireguardOutbound(tag, domainStrategy string, localAddress []string, server string, serverPort int, privateKey, peerPublicKey string) *WireGuardOutbound {
    return &WireGuardOutbound{
        Type:           "wireguard",
        Tag:            tag,
        DomainStrategy: domainStrategy,
        LocalAddress:   localAddress,
        PrivateKey:     privateKey,
        Server:         server,
        ServerPort:     serverPort,
        PeerPublicKey:  peerPublicKey,
        MTU:            1280,
    }
}