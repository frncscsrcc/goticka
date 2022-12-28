package version

import "fmt"

type Version struct {
	Major uint
	Minor uint
	Patch uint
	Rel   string
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Rel)
}

var version Version = Version{0, 0, 1, "alpha"}

func GetVersion() Version {
	return version
}
