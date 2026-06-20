package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	moduleName := "MAIN"
	if val, ok := entry.Data["module"]; ok {
		moduleName = fmt.Sprintf("%v", val)
	}

	// b.WriteString(fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, moduleName, message))
	fmt.Fprintf(b, "[%s] [%s] [%s] %s", timestamp, level, moduleName, message)

	hasExtraData := false

	for key, value := range entry.Data {
		if key != "module" {
			if !hasExtraData {
				b.WriteString(" |")
				hasExtraData = true
			}
			// b.WriteString(fmt.Sprintf(" %s=%v", key, value))
			fmt.Fprintf(b, "%s=%v", key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil

}

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viper.GetInt("log.level")))
	log.SetFormatter(&Formatter{})

	return log
}
