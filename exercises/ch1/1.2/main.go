// Echo1.2 prints one arg per line
package main

import (
	"fmt"
	"os"
)

func main() {
	for ind, arg := range os.Args {
		fmt.Println(ind, arg)
	}
}
