package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/caarlos0/env/v11"
)

const (
	fileConfigEnvName = "CONFIG"
	fileConfigArgName = "c"
)

// Handler интерфейс парсера конфигурации.
type Handler[T any] interface {
	Handle(*T) (*T, error)
}

// FileHandler конфиг из файла.
type FileHandler[T any] struct {
	next      Handler[T]
	arguments []string
}

// Handle обработчик.
func (f *FileHandler[T]) Handle(c *T) (*T, error) {
	var (
		err    error
		path   string
		exists bool
	)

	if path, exists = os.LookupEnv(fileConfigEnvName); !exists {
		for i, argv := range f.arguments {
			if strings.HasPrefix(argv, "-"+fileConfigArgName+"=") {
				path = strings.TrimPrefix(argv, "-"+fileConfigArgName+"=")
				break
			} else if argv == "-"+fileConfigArgName && len(f.arguments) > i+1 {
				path = f.arguments[i+1]
				break
			}
		}
	}

	if path == "" {
		return f.next.Handle(c)
	}

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var content []byte
	if content, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, c); err != nil {
		return nil, err
	}

	if f.next != nil {
		return f.next.Handle(c)
	}
	return c, nil
}

// NewFileHandler конструктор.
func NewFileHandler[T any](next Handler[T], arguments []string) *FileHandler[T] {
	return &FileHandler[T]{
		next:      next,
		arguments: arguments,
	}
}

// EnvHandler конфиг из переменных среды.
type EnvHandler[T any] struct {
	next Handler[T]
}

// Handle обработчик.
func (e *EnvHandler[T]) Handle(c *T) (*T, error) {
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	if e.next != nil {
		return e.next.Handle(c)
	}
	return c, nil
}

// NewEnvHandler конструктор.
func NewEnvHandler[T any](next Handler[T]) *EnvHandler[T] {
	return &EnvHandler[T]{
		next: next,
	}
}

// ArgumentsHandler конфиг из аргументов.
type ArgumentsHandler[T any] struct {
	next      Handler[T]
	arguments []string
}

// Handle обработчик.
func (a *ArgumentsHandler[T]) Handle(c *T) (*T, error) {
	var (
		p   *arg.Parser
		err error
	)
	p, err = arg.NewParser(arg.Config{
		IgnoreEnv:     true,
		IgnoreDefault: false,
	}, c)
	if err != nil {
		return nil, err
	}

	p.MustParse(a.arguments)

	if a.next != nil {
		return a.next.Handle(c)
	}
	return c, nil
}

// NewArgumentsHandler конструктор.
func NewArgumentsHandler[T any](next Handler[T], arguments []string) *ArgumentsHandler[T] {
	return &ArgumentsHandler[T]{
		next:      next,
		arguments: arguments,
	}
}
