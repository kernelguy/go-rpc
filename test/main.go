package test

import (
    "flag"
    "fmt"
    "os"
    "path"
    "runtime"
    "testing"
	log "github.com/Sirupsen/logrus"
)

func file_line(line int) string {
    _, fileName, fileLine, ok := runtime.Caller(line)
    var s string
    if ok {
        s = fmt.Sprintf("%s:%d", path.Base(fileName), fileLine)
    } else {
        s = ""
    }
    return s
}

func Begin(name string, line int) {
	log.Infof("***** %s ***** %s", name, file_line(line))
}

func End() {
	log.Info("***** End *****\n")
}

func Init(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}
}

func Run(m *testing.M) {
	log.Info("***** STARTING TEST *****")
	r := m.Run()
	log.Info("***** TEST FINISHED *****")
	os.Exit(r)
}
