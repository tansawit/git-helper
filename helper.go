package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func nilableString(s *string, maxLength int) string {
	if s == nil {
		return ""
	} else if len(*s) > maxLength && maxLength != 0 {
		return strings.Replace(*s, `"`, `'`, -1)[:maxLength]
	}
	return strings.Replace(*s, `"`, `'`, -1)
}
