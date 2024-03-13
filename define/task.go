package define

import "time"

type SectorFileType int

var PathTypes = []SectorFileType{FTUnsealed, FTSealed, FTCache, FTUpdate, FTUpdateCache}

const (
	FTUnsealed SectorFileType = 1 << iota
	FTSealed
	FTCache
	FTUpdate
	FTUpdateCache
)

func (ft SectorFileType) String() string {
	switch ft {
	case FTUnsealed:
		return "unsealed"
	case FTSealed:
		return "sealed"
	case FTCache:
		return "cache"
	case FTUpdate:
		return "update"
	case FTUpdateCache:
		return "update-cache"
	}
	return ""
}

type TaskType string

const (
	TTDataCid    TaskType = "seal/v0/datacid"
	TTAddPiece   TaskType = "seal/v0/addpiece"
	TTPreCommit1 TaskType = "seal/v0/precommit/1"
	TTPreCommit2 TaskType = "seal/v0/precommit/2"
	TTCommit1    TaskType = "seal/v0/commit/1"
	TTCommit2    TaskType = "seal/v0/commit/2"

	TTFinalize         TaskType = "seal/v0/finalize"
	TTFinalizeUnsealed TaskType = "seal/v0/finalizeunsealed"

	TTFetch     TaskType = "seal/v0/fetch"
	TTFetchLong TaskType = "seal/v0/fetchLong" //zcjs
	TTUnseal    TaskType = "seal/v0/unseal"

	TTReplicaUpdate         = "seal/v0/replicaupdate"
	TTProveReplicaUpdate1   = "seal/v0/provereplicaupdate/1"
	TTProveReplicaUpdate2   = "seal/v0/provereplicaupdate/2"
	TTRegenSectorKey        = "seal/v0/regensectorkey"
	TTFinalizeReplicaUpdate = "seal/v0/finalize/replicaupdate"

	TTDownloadSector = "seal/v0/download/sector"

	TTGenerateWindowPoSt  = "post/v0/windowproof"
	TTGenerateWinningPoSt = "post/v0/winningproof"
)

func (t TaskType) String() string {
	return string(t)
}

type SealMonitorType int

const (
	ApMonitor SealMonitorType = iota
	P1Monitor SealMonitorType = iota
	P2Monitor SealMonitorType = iota
	WSMonitor SealMonitorType = iota
	C2Monitor SealMonitorType = iota
)

func (s SealMonitorType) String() string {
	switch s {
	case ApMonitor:
		return AddPiece.String()
	case P1Monitor:
		return SealPreCommit1.String()
	case P2Monitor:
		return SealPreCommit2.String()
	case WSMonitor:
		return WaitSeed.String()
	case C2Monitor:
		return SealCommit2.String()
	}
	return ""
}
func (s SealMonitorType) Int() int {
	return int(s)
}
func (s SealMonitorType) TimeOut(ssize uint64) time.Duration {
	switch s {
	case ApMonitor:
		if ssize > Ss32GiB {
			return time.Minute * 30
		}
		return time.Minute * 10
	case P1Monitor:
		if ssize > Ss32GiB {
			return time.Hour * 10
		}
		return time.Hour * 5
	case P2Monitor:
		if ssize > Ss32GiB {
			return time.Minute * 30
		}
		return time.Minute * 20
	case WSMonitor:
		return time.Minute * 80
	case C2Monitor:
		return time.Minute * 20
	}
	return time.Hour * 12
}
