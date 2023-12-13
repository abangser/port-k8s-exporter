package config

import (
	"flag"
	"github.com/port-labs/port-k8s-exporter/pkg/goutils"
	"github.com/port-labs/port-k8s-exporter/pkg/port"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
	"os"
	"slices"
	"strings"
)

var keys []string

func prepareEnvKey(key string) string {
	newKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))

	if slices.Contains(keys, newKey) {
		klog.Fatalf("Application Error : Found duplicate config key: %s", newKey)
	}

	keys = append(keys, newKey)
	return newKey
}

func NewString(v *string, key string, defaultValue string, description string) {
	value := goutils.GetStringEnvOrDefault(prepareEnvKey(key), defaultValue)
	flag.StringVar(v, key, value, description)
}

func NewUInt(v *uint, key string, defaultValue uint, description string) {
	value := uint(goutils.GetUintEnvOrDefault(prepareEnvKey(key), uint64(defaultValue)))
	flag.UintVar(v, key, value, description)
}

type FileNotFoundError struct {
	s string
}

func (e *FileNotFoundError) Error() string {
	return e.s
}

func GetConfigFile(filepath string, resyncInterval uint, stateKey string, eventListenerType string) (*port.Config, error) {
	c := &port.Config{
		ResyncInterval:    resyncInterval,
		StateKey:          stateKey,
		EventListenerType: eventListenerType,
	}
	config, err := os.ReadFile(filepath)
	if err != nil {
		return c, &FileNotFoundError{err.Error()}
	}

	err = yaml.Unmarshal(config, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
