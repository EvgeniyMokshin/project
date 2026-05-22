package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level определяет уровни логирования
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger — структура логгера с настройками
type Logger struct {
	level   Level
	console io.Writer
	file    io.Writer
	mu      sync.Mutex
}

// NewLogger создаёт новый экземпляр логгера
func NewLogger(level Level, console io.Writer, file io.Writer) *Logger {
	return &Logger{
		level:   level,
		console: console,
		file:    file,
	}
}

// InitFileLogger инициализирует логгер с записью в файл
func InitFileLogger(level Level, logPath string) (*Logger, error) {
	// Создаём директорию, если её нет
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	// Открываем файл для записи логов (создаём, если не существует)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return NewLogger(level, os.Stdout, file), nil
}

// Debug логирует сообщение уровня DEBUG
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.log(DEBUG, format, args...)
}

// Info логирует сообщение уровня INFO
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level > INFO {
		return
	}
	l.log(INFO, format, args...)
}

// Warn логирует сообщение уровня WARN
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level > WARN {
		return
	}
	l.log(WARN, format, args...)
}

// Error логирует сообщение уровня ERROR
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// log — внутренняя функция для записи логов
func (l *Logger) log(level Level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := l.levelToString(level)
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, levelStr, message)

	// Пишем в консоль, если она настроена
	if l.console != nil {
		l.console.Write([]byte(logEntry))
	}

	// Пишем в файл, если он настроен
	if l.file != nil {
		l.file.Write([]byte(logEntry))
	}
}

// levelToString преобразует уровень логирования в строку
func (l *Logger) levelToString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Close закрывает файл логов, если он открыт
func (l *Logger) Close() error {
	if closer, ok := l.file.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
