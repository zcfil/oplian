package define

const (
	WindowedPostWallet  = "windowedPostWallet"
	DiskIo              = "diskIo"
	Date                = "date"
	SoftwarePackage     = "softwarePackage"
	SoftwarePath        = "softwarePath"
	SpeedUpFile         = "speedUpFile"
	ProofParam          = "proofParam"
	LotusHeight         = "lotusHeight"
	WindowedPostBalance = "/usr/local/sbin/lotus-miner actor control list | grep post | awk '{print $5}'| sed 's/\\x1b\\[[^\\x1b]*m//g'|awk  -F. '{print $1}'"
	DiskIoRes           = "dmesg  | grep I/O | grep dev | awk '{print $6}' | awk -F \",\" '{print $1}' | uniq"
	DateRes             = "date \"+%Y-%m-%d %H:%M:%S\""
	DuSh                = "du -sh * "
	LotusHeightReq      = "/usr/local/sbin/lotus  sync  status   |grep   'Height diff' |uniq  |awk  '{print $NF}'|sort  -nr |head  -n1"
)
