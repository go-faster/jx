package bench

// Common names for benchmarks.
const (
	// Encode is name for encoding benchmarks.
	Encode = "Encode"
	// Decode is name for decoding benchmarks.
	Decode = "Decode"
	// Scan is name for scanning benchmarks.
	Scan = "Scan"
	// JX is name for benchmarks related to go-faster/jx package.
	JX = "jx"
	// Std is name for benchmarks related to encoding/json.
	Std = "std"
	// Sonic is name for benchmarks related to bytedance/sonic package.
	Sonic = "sonic"
	// JSONIter for json-iterator/go.
	JSONIter = "json-iterator"
	// EasyJSON for mailru/easyjson.
	EasyJSON = "easyjson"
	// FFJSON for pquerna/ffjson.
	FFJSON = "ffjson"
	// JScan for romshark/jscan.
	JScan = "jscan"
	// Baseline directly writes string to buffer, no encoding.
	Baseline = "Baseline"
)
