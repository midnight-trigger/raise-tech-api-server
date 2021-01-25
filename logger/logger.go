package logger

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	cron "github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

func Init(runtime string) {
	Logger, err := newLogger("./log/", "log", "0 0 * * *", runtime)
	if err != nil {
		log.Fatal(err)
	}
	L = Logger
}

func newLogger(path, fileName, schedule, runtime string) (*zap.SugaredLogger, error) {
	c := cron.New()

	cronLogger := &autoDailyLogger{
		path:      path,
		fileName:  fileName,
		fileWrite: &fileWrite{},
		c:         c,
	}
	//cronLogger.setOutputFile()
	//c.AddFunc(schedule, cronLogger.setOutputFile)
	//c.Start()

	return cronLogger.getLogger(runtime)
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

func (a *autoDailyLogger) setOutputFile() {
	output, err := os.OpenFile(getPath(filepath.Join(a.path, a.fileName)), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0766)
	if a.fileWrite.file != nil {
		a.fileWrite.file.Close()
	}
	if err == nil {
		a.fileWrite.file = output
	}
}

func (a *autoDailyLogger) getLogger(runtime string) (*zap.SugaredLogger, error) {

	var config zap.Config
	var loggerConfig []byte
	var err error
	switch runtime {
	case "":
		loggerConfig, err = ioutil.ReadFile("./configs/log/test.yml")
	case "test":
		loggerConfig, err = ioutil.ReadFile("../../configs/log/test.yml")
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(loggerConfig, &config); err != nil {
		log.Fatalf("failed to unmarshal the log config file as yaml: err=%s", err)
	}

	//config.OutputPaths = append(config.OutputPaths, getPath(filepath.Join(a.path, a.fileName)))

	l, err := config.Build()
	if err != nil {
		log.Fatalf("failed to build the logger: err=%s", err)
	}

	return l.Sugar(), nil
}

func getPath(path string) string {
	return path + "_" + time.Now().Format("2006-01-02") + ".log"
}
