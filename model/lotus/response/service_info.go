package response

type ServerInfo struct {
	Id    uint64 `json:"id"`
	Ip    string `json:"ip"`
	Token string `json:"token"`
	Actor string `json:"actor"`
}
