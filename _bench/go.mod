module bench

go 1.24.0

require (
	github.com/bytedance/sonic v1.14.1
	github.com/go-faster/errors v0.7.1
	github.com/go-faster/jx v1.1.0
	github.com/json-iterator/go v1.1.12
	github.com/mailru/easyjson v0.9.0
	github.com/minio/simdjson-go v0.4.5
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/romshark/jscan v1.0.0
	github.com/stretchr/testify v1.11.1
	github.com/sugawarayuuta/sonnet v0.0.0-20231004000330-239c7b6e4ce8
	github.com/valyala/fastjson v1.6.4
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.21.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// replace to current repository version
replace github.com/go-faster/jx => ../

// CVE-2022-28948
replace gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0
