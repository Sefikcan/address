package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sefikcan/address-consumer/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

// Logger interface, defines log functions
// With this function, we can write logs different levels
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// logger struct, implements Logger interface
// this struct contains configuration settings and zap logger object
type logger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

func (l *logger) newElasticsearch(level zapcore.Level) zapcore.Core {
	cfg := elasticsearch.Config{
		Addresses: []string{l.cfg.Logger.ElasticsearchUrl},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create elasticsearch client: %v", err))
	}

	writeSyncer := zapcore.AddSync(newElasticsearchWriter(es, l.cfg.Logger.IndexName))
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	return zapcore.NewCore(encoder, writeSyncer, level)
}

type elasticSearchWriter struct {
	es    *elasticsearch.Client
	index string
}

func (e elasticSearchWriter) Write(p []byte) (n int, err error) {
	req := esapi.IndexRequest{
		Index:      e.index,
		DocumentID: fmt.Sprintf("%d", time.Now().UnixNano()),
		Body:       bytes.NewReader(p),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.IsError() {
		return 0, fmt.Errorf("elasticsearch index request failed: %s", res.String())
	}

	return len(p), nil
}

func newElasticsearchWriter(es *elasticsearch.Client, index string) *elasticSearchWriter {
	return &elasticSearchWriter{
		es:    es,
		index: index,
	}
}

// Above functions, implements logger interface functions
// In this way, we can write log messages by different log level

func (l *logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args)
}

func (l *logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args)
}

func (l *logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args)
}

func (l *logger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *logger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// loggerLevelMap, converts string value to zapcore.level type
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLogLevel, return zapcore.level based on the log level in config
func (l *logger) getLogLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}

// InitLogger initialize zap logger
func (l *logger) InitLogger() {
	logLevel := l.getLogLevel(l.cfg)        // gets log level from config
	logWriter := zapcore.AddSync(os.Stderr) // routes log output to standard error output

	encoderCfg := zap.NewProductionEncoderConfig()

	// encoder configuration and specification to json or console format
	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg) // write log in console format
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg) // write log in json format
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder // set time format ISO8601
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	elasticSearch := l.newElasticsearch(logLevel)
	tee := zapcore.NewTee(core, elasticSearch)
	logger := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	// Start SugaredLogger instance
	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err) // Logger logs sync error
	}
}

// NewLogger function create a new logger instance and return this object
func NewLogger(cfg *config.Config) Logger {
	return &logger{
		cfg: cfg,
	}
}
