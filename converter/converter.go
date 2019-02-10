package converter

import (
	"fmt"
)

// Converter is an interface to allow convert into files different formats
type Converter interface {
	Convert(opts ...[]string) error
}

// ConverterFormat output format type supported by the application
type ConverterFormat string

const (
	//EpubFormat output format epub
	EpubFormat ConverterFormat = "epub"
)

// NewConverter get the strategy converter to convert the files into the format selected
func NewConverter(format string) (Converter, error) {
	if ok := isAllowedConverterFormat(format); !ok {
		return nil, fmt.Errorf("Format %s is not allowed", format)
	}

	switch format {
	case string(EpubFormat):
		return NewPandocConverter()
	}

	return nil, fmt.Errorf("Not converted strategy created for %s", format)
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
