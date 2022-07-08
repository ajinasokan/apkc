package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// run (re)installs apk and runs the main activity
func run() {
	// get list of all devices
	cmd := exec.Command(adbPath, "devices")
	out, err := cmd.Output()
	if err != nil {
		LogF("run", err, string(out))
	}

	devices := strings.Split(strings.TrimSpace(string(out)), "\n")[1:]
	if len(devices) == 0 {
		LogF("run", "no devices found")
	}

	// if there are more than one device prompt user to pick one
	var deviceNum int64 = 1
	if len(devices) > 1 {
		for i, d := range devices {
			fmt.Printf("[%d] %s\n", i+1, d)
		}
		deviceNum, err = strconv.ParseInt(prompt("Choose device to run app (1):", "1"), 10, 64)
		if err != nil {
			LogF("run", "invalid number")
		}
	}

	device := strings.Split(devices[deviceNum-1], "\t")
	serial := device[0]

	// install apk in the device
	LogI("run", "installing in", device[1]+"("+device[0]+")")

	cmd = exec.Command(adbPath, "-s", serial, "install", "-r", filepath.Join("build", "app.apk"))
	out, err = cmd.CombinedOutput()
	if err != nil {
		LogF("run", err, string(out))
	}

	// collect package name and activity name from manifest
	mb, err := ioutil.ReadFile("AndroidManifest.xml")
	if err != nil {
		LogF("run", string(out))
	}
	manifest := Manifest{}
	xml.Unmarshal(mb, &manifest)
	pkg := manifest.Package
	activityName := manifest.Application.MainActivity()
	if activityName == "" {
		LogF("run", "couldn't find main activity")
	}

	// launch activity
	LogI("run", "launching main activity")
	cmd = exec.Command(adbPath, "shell", "am", "start", "-W", "-S", "-n", fmt.Sprintf("%s/%s", pkg, activityName))
	out, err = cmd.CombinedOutput()
	if err != nil {
		LogF("run", err, string(out))
	}

	// get pid of app
	cmd = exec.Command(adbPath, "shell", "pidof", "-s", pkg)
	out, err = cmd.Output()
	if err != nil {
		LogF("run", err, string(out))
	}
	pid := strings.TrimSpace(string(out))

	// stream logcat for the pid
	cmd = exec.Command(adbPath, "logcat", "-v", "color", "--pid", pid)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		LogF("run", err, string(out))
	}
}
