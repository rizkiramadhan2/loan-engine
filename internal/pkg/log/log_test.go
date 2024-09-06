package log

import (
	"errors"
	"testing"
)

func TestLogAllInfo(t *testing.T) {
	Info("test info")
	Infof("%s", "test infof")
	Infoln("test infoln")
	InfoWithFields("test info with field", map[string]interface{}{})
}

func TestLogAllPrint(t *testing.T) {
	Print("test print")
	Printf("%s", "test printf")
	Println("test println")
}

func TestLogAllError(t *testing.T) {
	Error("test error")
	Errors(errors.New("some err"))
	Errorf("%s", "test errorf")
	Errorln("test errorln")
	ErrorWithFields("test error with field", map[string]interface{}{})
}

func TestLogAllWarn(t *testing.T) {
	Warn("test warn")
	Warnf("%s", "test warnf")
	Warnln("test warnln")
	WarnWithFields("test warn with field", map[string]interface{}{})
}

func TestLogAllDebug(t *testing.T) {
	Debug("test debug")
	Debugf("%s", "test debugf")
	Debugln("test debugln")
	DebugWithFields("test debug with field", map[string]interface{}{})
}
