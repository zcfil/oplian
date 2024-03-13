package response

import "oplian/model/lotus/request"

type CarFile struct {
	Path      string `json:"path"`
	FileName  string `json:"fileName"`
	PieceCid  string `json:"pieceCid"`
	PieceSize int    `json:"pieceSize"`
}

type FileInfo struct {
	Ip      string `json:"ip"`
	Path    string `json:"path"`
	FileNum int    `json:"fileNum"`
}

type WorkerOpList struct {
	Id         uint   `json:"id"`
	OpId       string `json:"opId"`
	Ip         string `json:"ip"`
	ServerName string `json:"serverName"`
}

type WorkerCarFiles struct {
	RelationId  int    `json:"relationId"`
	FileName    string `json:"fileName"`
	FileIndex   int    `json:"fileIndex"`
	FileStr     string ` json:"fileStr"`
	CarFileName string `json:"carFileName"`
	PieceCid    string `json:"pieceCid"`
	PieceSize   int    `json:"pieceSize"`
	CarSize     int    `json:"carSize"`
	DataCid     string `json:"dataCid"`
	InputDir    string `json:"inputDir"`
}

type FileParam struct {
	FileIndex int
	FileSort  int
	FileName  string
	FilePath  string
	FileSize  int64
	FileStr   string
}

type SectorCids struct {
	Unsealed string
	Sealed   string
}

type Files struct {
	Files []FileParam `json:"files"`
}

type ApRes struct {
	Code  int
	PInfo []request.PieceInfo
}

type P1Res struct {
	Code int
	Out  []byte
}

type P2Res struct {
	Code int
	Cid  SectorCids
}

type MoveRes struct {
	Code   int
	Number uint64
}

type ApResData struct {
	Code int   `json:"Code"`
	Data ApRes `json:"Data"`
}

type P1ResData struct {
	Code int   `json:"Code"`
	Data P1Res `json:"Data"`
}

type P2ResData struct {
	Code int   `json:"Code"`
	Data P2Res `json:"Data"`
}

type MoveResData struct {
	Code int     `json:"Code"`
	Data MoveRes `json:"Data"`
}
