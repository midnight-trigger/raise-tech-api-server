package logger

import (
	"log"
	"os"

	"github.com/midnight-trigger/todo/configs"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var L *zap.SugaredLogger

func Init() {
	Logger, err := newLogger("./log/", "log", "0 0 * * *")
	if err != nil {
		log.Fatal(err)
	}
	L = Logger
}

func newLogger(path, fileName, schedule string) (*zap.SugaredLogger, error) {
	c := cron.New()

	cronLogger := &autoDailyLogger{
		path:      path,
		fileName:  fileName,
		fileWrite: &fileWrite{},
		c:         c,
	}

	return cronLogger.getLogger()
}

type fileWrite struct {
	file *os.File
}

func (f *fileWrite) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *fileWrite) Close() error {
	return f.file.Close()
}

type autoDailyLogger struct {
	path      string
	fileName  string
	fileWrite *fileWrite
	c         *cron.Cron
}

func (a *autoDailyLogger) getLogger() (*zap.SugaredLogger, error) {

	var config zap.Config
	var loggerConfig []byte
	var err error
	if configs.IsProd() {
		loggerConfig = prodYaml
	} else {
		loggerConfig = devYaml
	}

	if err := yaml.Unmarshal(loggerConfig, &config); err != nil {
		log.Fatalf("failed to unmarshal the log config file as yaml: err=%s", err)
	}

	l, err := config.Build()
	if err != nil {
		log.Fatalf("failed to build the logger: err=%s", err)
	}

	return l.Sugar(), nil
}
