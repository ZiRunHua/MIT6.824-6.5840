package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"MIT6.824-6.5840/util/color"
	"MIT6.824-6.5840/util/config"
)

var (
	debug         = false
	logDetail     = false
	lab           = ""
	recordProfile = false
)

func init() {
	debug = config.Get("DEBUG", false)
	logDetail = config.Get("DEBUG_DETAIL", false)
	lab = strings.ToLower(config.Get("DEBUG_LAB", ""))
	recordProfile = config.Get("RECORD_PROFILE", false)
	if recordProfile {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		go func() {
			fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}
}

func Println(title string, name string, msg string, details ...interface{}) {
	if !debug {
		return
	}
	if len(details) > 0 && logDetail {
		str, _ := json.Marshal(details)
		fmt.Printf(
			"[%s] %-25s %-25s %s %s\n",
			color.Yellow(time.Now().Format("04:05.0000000")),
			color.Green(title),
			color.Blue(name),
			color.Cyan(msg),
			string(str),
		)
	} else {
		fmt.Printf(
			"[%s] %-25s %-25s %s\n",
			color.Yellow(time.Now().Format("04:05.0000000")),
			color.Green(title),
			color.Blue(name),
			color.Cyan(msg),
		)
	}
}

func Lab3(title string, name string, msg string, details ...interface{}) {
	if len(lab) > 0 && lab != "lab3" {
		return
	}
	Println(title, name, msg, details...)
}
func Lab4(title string, name string, msg string, details ...interface{}) {
	if len(lab) > 0 && lab != "lab4" {
		return
	}
	Println(title, name, msg, details...)
}

func Lab5(title string, name string, msg string, details ...interface{}) {
	if len(lab) > 0 && lab != "lab5" {
		return
	}
	Println(title, name, msg, details...)
}

func Memory() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Println(
		color.Magenta(
			fmt.Sprintf(
				"Alloc: %v KB TotalAlloc: %v KB Sys: %v KB NumGC: %v",
				memStats.Alloc/1024, memStats.TotalAlloc/1024,
				memStats.Sys/1024, memStats.NumGC,
			),
		),
	)
}

func RecordProfile(name string) {
	if !recordProfile {
		return
	}
	f, err := os.Create(fmt.Sprintf("%s_profile_%d.prof", name, time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := pprof.Lookup(name)
	if err := p.WriteTo(f, 0); err != nil {
		panic(err)
	}
}
