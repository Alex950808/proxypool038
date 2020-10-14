package config

import "github.com/Alex950808/proxypool038/pkg/tool"

type Source struct {
	Type    string       `json:"type" yaml:"type"`
	Options tool.Options `json:"options" yaml:"options"`
}
