package define

type DataTypeToInt int
type DataTypeToString string
type DataTypeToInt64 int64

const (
	FileDistribution     DataTypeToInt = 1
	FileCopy             DataTypeToInt = 2
	FileCopyToOp         DataTypeToInt = 3
	FileCopyToGateWay    DataTypeToInt = 4
	StrictMode           DataTypeToInt = 2
	Pause                DataTypeToInt = 2
	Stop                 DataTypeToInt = 4
	Start                DataTypeToInt = 1
	FileByte             DataTypeToInt = 4194000
	ErrorCode            DataTypeToInt = -1
	BusinessWarn         DataTypeToInt = 6
	TaskSuccess          DataTypeToInt = 1
	TaskFailed           DataTypeToInt = 2
	TaskInProgress       DataTypeToInt = 3
	TaskPartialSuccess   DataTypeToInt = 4
	FileGetting          DataTypeToInt = 1
	FileFinish           DataTypeToInt = 2
	FileError            DataTypeToInt = 3
	FileLocalUpload      DataTypeToInt = 1
	FileHostUpload       DataTypeToInt = 2
	FileDistributeUpload DataTypeToInt = 3

	NodeMachine     DataTypeToString = "1"
	StartService    DataTypeToString = "start"
	StopService     DataTypeToString = "stop"
	RedoSectorsTask DataTypeToString = "redo_sectors_task"

	GenerallyFile  DataTypeToInt64 = 1
	ProveFile      DataTypeToInt64 = 2
	HeightFile     DataTypeToInt64 = 3
	SnapshotFile   DataTypeToInt64 = 4
	MinerFile      DataTypeToInt64 = 5
	AddLocalUpload DataTypeToInt64 = 1
	AddOnline      DataTypeToInt64 = 2
	AddNodeCopy    DataTypeToInt64 = 3
)

func (d DataTypeToInt) Int() int {
	return int(d)
}

func (d DataTypeToInt64) Int64() int64 {
	return int64(d)
}

func (d DataTypeToString) String() string {
	return string(d)
}
