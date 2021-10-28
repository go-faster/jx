package json

import (
	"sync"

	"github.com/modern-go/concurrent"
)

// Valid reports whether json in data is valid.
func Valid(data []byte) bool {
	return ConfigDefault.Valid(data)
}

// Config customize how the API should behave.
// The API is created from Config by Froze.
type Config struct {
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

// API the public interface of this package.
// Primary Marshal and Unmarshal.
type API interface {
	IteratorPool
	StreamPool
	Valid(data []byte) bool
}

// ConfigDefault the default API
var ConfigDefault = Config{}.Froze()

// ConfigCompat tries to be 100% compatible with standard library behavior
var ConfigCompat = Config{
	SortMapKeys: true,
}.Froze()

// ConfigFastest marshals float with only 6 digits precision
var ConfigFastest = Config{
	MarshalFloatWith6Digits:       true, // will lose precession
	ObjectFieldMustBeSimpleString: true, // do not unescape object field
}.Froze()

type frozenConfig struct {
	configBeforeFrozen            Config
	sortMapKeys                   bool
	indentionStep                 int
	objectFieldMustBeSimpleString bool
	onlyTaggedField               bool
	disallowUnknownFields         bool
	decoderCache                  *concurrent.Map
	encoderCache                  *concurrent.Map
	streamPool                    *sync.Pool
	iteratorPool                  *sync.Pool
	caseSensitive                 bool
}

func (cfg *frozenConfig) initCache() {
	cfg.decoderCache = concurrent.NewMap()
	cfg.encoderCache = concurrent.NewMap()
}

// Froze forge API from config
func (cfg Config) Froze() API {
	api := &frozenConfig{
		sortMapKeys:                   cfg.SortMapKeys,
		indentionStep:                 cfg.IndentionStep,
		objectFieldMustBeSimpleString: cfg.ObjectFieldMustBeSimpleString,
		onlyTaggedField:               cfg.OnlyTaggedField,
		disallowUnknownFields:         cfg.DisallowUnknownFields,
		caseSensitive:                 cfg.CaseSensitive,
	}
	api.streamPool = &sync.Pool{
		New: func() interface{} {
			return NewStream(api, nil, 512)
		},
	}
	api.iteratorPool = &sync.Pool{
		New: func() interface{} {
			return NewIterator(api)
		},
	}
	api.initCache()
	api.configBeforeFrozen = cfg
	return api
}

func (cfg *frozenConfig) Valid(data []byte) bool {
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	iter.Skip()
	return iter.Error == nil
}
