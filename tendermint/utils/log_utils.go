package utils

import (
	"github.com/rs/zerolog"
	zelog "github.com/rs/zerolog/log"
	"log"
	"os"
)

var (
	WarningStdOutLogger *log.Logger
	InfoStdOutLogger *log.Logger
	ErrorStdOutLogger *log.Logger
	FatalStdOutLogger *log.Logger
)

func InitLogConsole(){
	WarningStdOutLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoStdOutLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorStdOutLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalStdOutLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetupZeroLogFile(logFilePath string) * os.File{
	EnsurePath(logFilePath, DefaultDirMod)
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	fileLog, err := os.OpenFile(logFilePath, mod, DefaultFileMod)
	if err != nil{
		panic(err)
	}

	//defer func() {
	//	if fileLog != nil {
	//		fileLog.Close()
	//	}
	//}()

	zelog.Logger = zelog.Output(zerolog.ConsoleWriter{Out: fileLog})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return fileLog
}