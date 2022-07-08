package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	findSDKs()

	switch os.Args[1] {
	case "doctor":
		doctor()
	case "create":
		create()
	case "build":
		build()
	case "run":
		run()
	case "clean":
		clean()
	}

	LogV("apkc", "finished in", time.Since(start))
}

func printHelp() {
	fmt.Print(`
commands:

doctor - check if Android sdk is accessible
create - create a new project
build  - build apk
run    - build apk and run the available device
clean  - delete build/ dir
`)
}
