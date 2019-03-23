package venom

import (
    "os"
    "log"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
    testIO := []struct {
        tc   string
        lgr  LoggingInterface
        tipe interface{}
    }{
        {
            tc:   "should be able to set a no-op logging interface",
            lgr:  &TestLogger{},
            tipe: &TestLogger{},
        },
        {
            tc:   "should be able set to a new io.Writer (stderr)",
            lgr:  log.New(os.Stderr, "", 0),
            tipe: &log.Logger{},
        },
    }

    for _, test := range testIO {
        t.Run(test.tc, func(t *testing.T) {
            l := NewLogable()
            defaultLog := l.GetLogger().(*log.Logger)
            assert.IsType(t, &log.Logger{}, defaultLog)
            
            l.SetLogger(test.lgr)
            customLog := l.GetLogger()
            assert.IsType(t, test.tipe, customLog)
        })
    }
}


