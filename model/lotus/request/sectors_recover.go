package request

type SectorID struct {
	Miner  uint64
	Number uint64
}

type PieceInfo struct {
	Size     int64
	PieceCID string
}

type P1Run struct {
	PieceInfo       []PieceInfo
	St              SectorTicket
	TaskDetailId    uint64
	SectorRecoverId uint64
}
type P2Run struct {
	PreCommit1Out   []byte
	Sector          SectorRef
	TaskDetailId    uint64
	SectorRecoverId uint64
}
type SectorTicket struct {
	Sector    SectorRef
	Ticket    []byte
	PieceCid  string
	PieceSize int
}

type SectorRef struct {
	ID              SectorID
	ProofType       int64
	TaskDetailId    uint64
	SectorRecoverId uint64
}

type RunTaskParam struct {
	Miner       uint64
	Number      uint64
	ProofType   int64
	PieceCid    string
	PieceSize   int
	CarPath     string
	Unsealed    string
	Sealed      string
	Cache       string
	Update      string
	UpdateCache string
	Ticket      []byte
	Piece       []PieceInfo
	Phase1Out   []byte
	Size        int64
}
