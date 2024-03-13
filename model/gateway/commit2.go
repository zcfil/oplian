package gateway

type SectorID struct {
	Miner  uint64
	Number uint64
}

type ProofResult struct {
	Proof []byte `json:"proof"`
	Err   string `json:"err"`
}

type SectorRef struct {
	ID        SectorID
	ProofType int64
}

type SealerParam struct {
	Sector    SectorRef
	Phase1Out []byte
	Status    int
	ID        SectorID
}
