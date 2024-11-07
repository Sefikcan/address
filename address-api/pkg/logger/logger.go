package logger

import (
	"github.com/sefikcan/address/address-api/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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

	var encoderCfg zapcore.EncoderConfig
	// According to dev or prod environment set encoder setting
	if l.cfg.Server.Mode == "Dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

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
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

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
