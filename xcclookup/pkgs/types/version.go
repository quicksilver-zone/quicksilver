package types

import "encoding/json"

var (
	VERSION             = "xcclookup/development"
	QUICKSILVER_VERSION = ""
	COMMIT              = ""
)

type Version struct {
	Version            string
	QuicksilverVersion string
	QuicksilverCommit  string
}

func GetVersion() ([]byte, error) {
	v := Version{
		Version:            VERSION,
		QuicksilverVersion: QUICKSILVER_VERSION,
		QuicksilverCommit:  COMMIT,
	}
	jsonOut, err := json.Marshal(v)
	return jsonOut, err
}
