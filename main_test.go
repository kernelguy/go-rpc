package gorpc

import (
    "flag"
    "fmt"
    "os"
    "path"
    "runtime"
    "testing"
	log "github.com/Sirupsen/logrus"
)

func file_line() string {
    _, fileName, fileLine, ok := runtime.Caller(2)
    var s string
    if ok {
        s = fmt.Sprintf("%s:%d", path.Base(fileName), fileLine)
    } else {
        s = ""
    }
    return s
}

func beginTest(name string) {
	log.SetLevel(log.DebugLevel)
	log.Infof("***** %s ***** %s", name, file_line())
}

func endTest() {
	log.Info("***** End *****\n")
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.Info("***** STARTING TEST *****")
	r := m.Run()
	log.Info("***** TEST FINISHED *****")
	os.Exit(r)
}
