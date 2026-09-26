package v2

import (
	"fmt"
)

func ExampleMountURL() {
	fmt.Println(MountURL("", "", "localhost", "27017"))
	fmt.Println(MountURL("user", "pass", "localhost", "27017"))
	// Output:
	// mongodb://localhost/27017
	// mongodb://user:pass@localhost/27017
}
