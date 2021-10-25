package offtz

import (
	"fmt"
	"io/ioutil"
	"regexp/syntax"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
	// "golang.org/x/sys/windows/registry"
)

// zoneDirs adapted from https://golang.org/src/time/zoneinfo_unix.go
// https://golang.org/doc/install/source#environment
// list of available GOOS as of 10th Feb 2017
// android, darwin, dragonfly, freebsd, linux, netbsd, openbsd, plan9, solaris,windows
var zoneDirs = map[string]string{
	"android":   "/system/usr/share/zoneinfo/",
	"darwin":    "/usr/share/zoneinfo/",
	"dragonfly": "/usr/share/zoneinfo/",
	"freebsd":   "/usr/share/zoneinfo/",
	"linux":     "/usr/share/zoneinfo/",
	"netbsd":    "/usr/share/zoneinfo/",
	"openbsd":   "/usr/share/zoneinfo/",
	// "plan9":"/adm/timezone/", -- no way to test this platform
	"solaris": "/usr/share/lib/zoneinfo/",
	"windows": `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Time Zones\`,
}

var timeZones = make(map[int][]string)

func init() {
	if runtime.GOOS == "nacl" || runtime.GOOS == "" || runtime.GOOS == "windows" {
		panic("unsupported platform")
	}

	readTZFile("")
}

// readTZFile ... read timezone file and append into timeZones slice
func readTZFile(path string) {
	files, _ := ioutil.ReadDir(zoneDirs[runtime.GOOS] + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			readTZFile(path + "/" + f.Name())
		} else {
			tz := (path + "/" + f.Name())[1:]

			// convert string to rune
			tzRune, _ := utf8.DecodeRuneInString(tz[:1])

			if syntax.IsWordChar(tzRune) { // filter out entry that does not start with A-Za-z such as +VERSION
				loc, err := time.LoadLocation(tz)
				if err != nil {
					fmt.Println(err)
				}
				_, offset := time.Now().In(loc).Zone()

				appendIfMissing(timeZones[offset], tz)
			}
		}
	}

}

// appendIfMissing ... appends an element if not available in a slice
func appendIfMissing(list []string, str string) []string {
	for _, v := range list {
		if v == str {
			return list
		}
	}
	return append(list, str)
}
