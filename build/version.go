package build

const BuildVersion = "1.2.0"

var CurrentCommit string

func UserVersion() string {
	if CurrentCommit != "" {
		return BuildVersion + "+" + CurrentCommit
	}
	return BuildVersion
}
