package log

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
)

func TestLog(t *testing.T) { TestingT(t) }

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) SetUpTest(c *C) {
	// reset global loggers chain before every test
	loggers = []Logger{}
}

func (s *LogSuite) TestInit(c *C) {
	Init(newTestLogger("log1"), newTestLogger("log2"))
	c.Assert(len(loggers), Equals, 2)
	c.Assert(typeOf(loggers[0]), Equals, "*log.testLogger")
	c.Assert(typeOf(loggers[1]), Equals, "*log.testLogger")
}

func (s *LogSuite) TestInitWithConfig(c *C) {
	InitWithConfig(Config{Console, "info"}, Config{Syslog, "info"})
	c.Assert(len(loggers), Equals, 2)
	c.Assert(typeOf(loggers[0]), Equals, "*log.consoleLogger")
	c.Assert(typeOf(loggers[1]), Equals, "*log.sysLogger")
}

func (s *LogSuite) TestNewLogger(c *C) {
	l, err := NewLogger(Config{Console, "info"})
	c.Assert(err, IsNil)
	c.Assert(typeOf(l), Equals, "*log.consoleLogger")

	l, err = NewLogger(Config{Syslog, "warn"})
	c.Assert(err, IsNil)
	c.Assert(typeOf(l), Equals, "*log.sysLogger")

	l, err = NewLogger(Config{UDPLog, "error"})
	c.Assert(err, IsNil)
	c.Assert(typeOf(l), Equals, "*log.udpLogger")

	l, err = NewLogger(Config{"SuperDuperLogger", "info"})
	c.Assert(err, NotNil)
	c.Assert(l, IsNil)
}

func (s *LogSuite) TestDebugf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Debugf("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "DEBUG hello world\n")
	c.Assert(logger2.b.String(), Equals, "DEBUG hello world\n")
}

func (s *LogSuite) TestInfof(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Infof("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "INFO hello world\n")
	c.Assert(logger2.b.String(), Equals, "INFO hello world\n")
}

func (s *LogSuite) TestWarningf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Warningf("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "WARN hello world\n")
	c.Assert(logger2.b.String(), Equals, "WARN hello world\n")
}

func (s *LogSuite) TestErrorf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Errorf("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "ERROR hello world\n")
	c.Assert(logger2.b.String(), Equals, "ERROR hello world\n")
}

func typeOf(o interface{}) string {
	return reflect.TypeOf(o).String()
}

// testLogger helps in tests.
type testLogger struct {
	id string
	b  *bytes.Buffer
}

func newTestLogger(id string) *testLogger {
	return &testLogger{id, &bytes.Buffer{}}
}

func (l *testLogger) Writer(sev Severity) io.Writer {
	return l.b
}

func (l *testLogger) FormatMessage(sev Severity, caller *CallerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s\n", sev, fmt.Sprintf(format, args...))
}
