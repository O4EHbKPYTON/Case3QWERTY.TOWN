package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// Инициализация логгера
func init() {
	Log = logrus.New()

	// Настройка формата JSON
	Log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// Настройка вывода в файл с ротацией
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30,   // days
		Compress:   true, // gzip
	}
	
	Log.SetOutput(os.Stdout)
	Log.AddHook(&fileHook{writer: logFile})

	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
}

type fileHook struct {
	writer *lumberjack.Logger
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = h.writer.Write([]byte(line))
	return err
}

func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func Any(level logrus.Level, v interface{}, fields ...map[string]interface{}) {
	entry := Log.WithFields(combineFields(fields...))
	entry.Log(level, formatValue(v))
}

func DebugAny(v interface{}, fields ...map[string]interface{}) {
	Any(logrus.DebugLevel, v, fields...)
}

func InfoAny(v interface{}, fields ...map[string]interface{}) {
	Any(logrus.InfoLevel, v, fields...)
}

func WarnAny(v interface{}, fields ...map[string]interface{}) {
	Any(logrus.WarnLevel, v, fields...)
}

func ErrorAny(v interface{}, fields ...map[string]interface{}) {
	Any(logrus.ErrorLevel, v, fields...)
}

func formatValue(v interface{}) string {
	if v == nil {
		return "<nil>"
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return "<nil>"
	}

	switch v := v.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	case []byte:
		return string(v)
	}

	// Попробуем преобразовать в JSON
	if b, err := json.Marshal(v); err == nil {
		if str := string(b); str != "null" && !strings.HasPrefix(str, "{") && !strings.HasPrefix(str, "[") {
			return str
		}
	}

	// Для сложных структур используем %+v
	return fmt.Sprintf("%+v", v)
}

// combineFields объединяет несколько map полей
func combineFields(fields ...map[string]interface{}) logrus.Fields {
	result := make(logrus.Fields)
	for _, f := range fields {
		for k, v := range f {
			result[k] = v
		}
	}
	return result
}

// ================= КОНТЕКСТНОЕ ЛОГГИРОВАНИЕ =================

// WithContext возвращает логгер с контекстными полями
func WithContext(fields map[string]interface{}) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithRequestID добавляет ID запроса
func WithRequestID(requestID string) *logrus.Entry {
	return WithContext(map[string]interface{}{"request_id": requestID})
}

// WithUser добавляет информацию о пользователе
func WithUser(userID int64, username string) *logrus.Entry {
	return WithContext(map[string]interface{}{
		"user_id":  userID,
		"username": username,
	})
}
