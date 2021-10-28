package jir

import "sync"

// Valid reports whether json in data is valid.
func Valid(data []byte) bool {
	return Default.Valid(data)
}

// Config customize how the API should behave.
// The API is created from Config by API.
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

	private() // make interface private
}

// Default API.
var Default = Config{}.API()

// Compat tries to be compatible with standard library behavior.
var Compat = Config{
	SortMapKeys: true,
}.API()

// Fastest marshals float with only 6 digits precision.
var Fastest = Config{
	MarshalFloatWith6Digits:       true, // will lose precession
	ObjectFieldMustBeSimpleString: true, // do not unescape object field
}.API()

type frozenConfig struct {
	configBeforeFrozen            Config
	sortMapKeys                   bool
	indentionStep                 int
	objectFieldMustBeSimpleString bool
	onlyTaggedField               bool
	disallowUnknownFields         bool
	streamPool                    *sync.Pool
	iteratorPool                  *sync.Pool
	caseSensitive                 bool
}

func (cfg *frozenConfig) private() {}


// API creates new API from config
func (cfg Config) API() API {
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
	api.configBeforeFrozen = cfg
	return api
}

func (cfg *frozenConfig) Valid(data []byte) bool {
	iter := cfg.Iterator(data)
	defer cfg.PutIterator(iter)
	iter.Skip()
	return iter.Error == nil
}
