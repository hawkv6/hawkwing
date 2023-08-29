package version

import "fmt"

const MAJOR uint = 0
const MINOR uint = 1
const PATCH uint = 0

func GetVersion() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR, MINOR, PATCH)
}
