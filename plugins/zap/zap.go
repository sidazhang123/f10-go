package zap

import (
	"github.com/sidazhang123/f10-go/basic"
	"github.com/sidazhang123/f10-go/basic/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

var (
	l                              *Logger
	sp                             = string(filepath.Separator)
	errWS, infoWS, debugWS, warnWS zapcore.WriteSyncer
	debugConsoleWS                 = zapcore.Lock(os.Stdout)
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type Logger struct {
	sync.RWMutex
	*zap.Logger
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func init() {
	l = &Logger{Opts: &Options{}}
	basic.Register(initLogger)
}

func GetLogger() *Logger { return l }

func initLogger() {

	l.Lock()
	defer l.Unlock()
	if l.inited {
		l.Info("[InfoLogger] Logger initialized.")
		return
	}
	l.loadCfg()
	l.init()
	l.Info("[InfoLogger] Zap initializing completed.")
	l.inited = true

}

func (l *Logger) init() {
	l.setSync()
	var err error

	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.Logger.Sync()
}

func (l *Logger) loadCfg() {
	err := config.C().Path("zap", l.Opts)
	if err != nil {
		panic(err)
	}

	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
	} else {
		l.zapConfig = zap.NewProductionConfig()
	}

	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.ErrorOutputPaths = []string{"stderr"}
	}
	// logs stored to "logs" under the app's running directory
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}
	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}
	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}
	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}
	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 3
	}
	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 50
	}

	l.Opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

}

func (l *Logger) setSync() {
	f := func(fileName string) zapcore.WriteSyncer {
		return zapcore.AddSync(
			&lumberjack.Logger{
				Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fileName,
				MaxSize:    l.Opts.MaxSize,
				MaxAge:     l.Opts.MaxAge,
				MaxBackups: l.Opts.MaxBackups,
				LocalTime:  true,
				Compress:   true,
			})
	}
	errWS = f(l.Opts.ErrorFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	warnWS = f(l.Opts.WarnFileName)
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
	errPriority := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {

			return lvl > zapcore.WarnLevel && zapcore.ErrorLevel >= l.zapConfig.Level.Level()
		})
	warnPriority := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {

			return lvl == zapcore.WarnLevel && zapcore.WarnLevel >= l.zapConfig.Level.Level()
		})
	infoPriority := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {
			return lvl == zapcore.InfoLevel && zapcore.InfoLevel >= l.zapConfig.Level.Level()
		})
	debugPriority := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {

			return lvl == zapcore.DebugLevel && zapcore.DebugLevel >= l.zapConfig.Level.Level()
		})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),

		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
	}

	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})

}
