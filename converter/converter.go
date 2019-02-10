package converter

import (
	"fmt"
)

// Converter is an interface to allow convert into files different formats
type Converter interface {
	Convert() error
	WithOutputDir(string) error
	WithWebFiles([]string) error
}

// ConverterFormat output format type supported by the application
type ConverterFormat string

const (
	//EpubFormat output format epub
	EpubFormat ConverterFormat = "epub"
)

type Option func(Converter) error

// NewConverter get the strategy converter to convert the files into the format selected
func NewConverter(format string, opts ...Option) (Converter, error) {
	if ok := isAllowedConverterFormat(format); !ok {
		return nil, fmt.Errorf("Format %s is not allowed", format)
	}

	var c Converter
	var err error
	switch format {
	case string(EpubFormat):
		c, err = NewPandocConverter()
	}
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithOutputDir(outputDir string) Option {
	return func(c Converter) error {
		return c.WithOutputDir(outputDir)
	}
}

func WithWebFiles(webFiles []string) Option {
	return func(c Converter) error {
		return c.WithWebFiles(webFiles)
	}
}

func allowedConverterFormat() map[string]ConverterFormat {
	return map[string]ConverterFormat{
		string(EpubFormat): EpubFormat,
	}
}

func isAllowedConverterFormat(format string) bool {
	if _, ok := allowedConverterFormat()[format]; !ok {
		return false
	}

	return true
}
