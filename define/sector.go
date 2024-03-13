package define

const (
	SectorTypeCC = iota + 1
	SectorTypeDC
)

type SectorStatus uint64

const (
	AP            SectorStatus = 2
	PC1           SectorStatus = 3
	PC2           SectorStatus = 4
	Finish        SectorStatus = 5
	Failed        SectorStatus = 6
	Recover       SectorStatus = 1
	Recovering    SectorStatus = 4
	RecoverFailed SectorStatus = 2
	CC            SectorStatus = 3
	DC            SectorStatus = 4
)

const (
	Ss1MB    = 1 << 20
	Ss1GiB   = 1 << 30
	Ss512MiB = 512 << 20
	Ss32GiB  = 32 << 30
	Ss64GiB  = 64 << 30
	Ss234GiB = 234 << 30
)

func (s SectorStatus) Uint64() uint64 {
	return uint64(s)
}

type ReturnType string

const (
	AddPiece              ReturnType = "AddPiece"
	SealPreCommit1        ReturnType = "SealPreCommit1"
	SealPreCommit2        ReturnType = "SealPreCommit2"
	WaitSeed              ReturnType = "WaitSeed"
	SealCommit1           ReturnType = "SealCommit1"
	SealCommit2           ReturnType = "SealCommit2"
	FinalizeSector        ReturnType = "FinalizeSector"
	FinalizeReplicaUpdate ReturnType = "FinalizeReplicaUpdate"
	ReplicaUpdate         ReturnType = "ReplicaUpdate"
	ProveReplicaUpdate1   ReturnType = "ProveReplicaUpdate1"
	ProveReplicaUpdate2   ReturnType = "ProveReplicaUpdate2"
	ReleaseUnsealed       ReturnType = "ReleaseUnsealed"
	MoveStorage           ReturnType = "MoveStorage"
)

func (r ReturnType) String() string {
	return string(r)
}

func (r ReturnType) BeforeP2(status string) bool {
	switch status {
	case AddPiece.String(), SealPreCommit1.String(), "PreCommit1":
		return true
	}
	return false
}

const (
	//AddPiece              = "AddPiece"
	PreCommit1           = "PreCommit1"
	PreCommit2           = "PreCommit2"
	PreCommitting        = "PreCommitting"
	SubmitPreCommitBatch = "SubmitPreCommitBatch"
	PreCommitWait        = "PreCommitWait"
	PreCommitBatchWait   = "PreCommitBatchWait"
	//WaitSeed              = "WaitSeed"
	Committing            = "Committing"
	SubmitCommit          = "SubmitCommit"
	SubmitCommitAggregate = "SubmitCommitAggregate"
	CommitFinalize        = "CommitFinalize"
	CommitWait            = "CommitWait"
	CommitAggregateWait   = "CommitAggregateWait"
	Proving               = "Proving"
	Removed               = "Removed"
	//FinalizeSector        = "FinalizeSector"
)
