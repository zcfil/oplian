package oplocal

import (
	"oplian/service/pb"
	"time"
)

var (
	WaitRecord     bool // Prevent duplicate CC creation
	ImportCarCount int32
	APCount        = 1                                                            //Default concurrency
	PreCount1      int                                                            //Limit the number of AP+P1
	PreCount2      int                                                            //Limit p2 quantity
	Tasking        = Running{Run: make(map[string][]*pb.Task)}                    //Be performing a task
	OpWorkers      = LocalWorkerInfo{info: make(map[string]*pb.OpenWindow)}       //Local worker information
	Preselect      = Preselection{Pre1: 1, Pre2: 1, LastTime: time.Now(), Pre: 1} //Preselected buffer

	SealSectorPath = make(SealSectorDisk) //Sealed sector cache hard disk
	PathSealCount  = make(DiskSealCount)  //Number of hard disk sealing fans
	DiskCount      uint64
	IsDiskCount    = make(map[uint64]struct{})
	UnseledClearTime = time.Now()
)
