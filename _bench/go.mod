module bench

go 1.18

require (
	github.com/bytedance/sonic v1.8.8
	github.com/go-faster/jx v0.0.0-replaced
	github.com/json-iterator/go v1.1.12
	github.com/mailru/easyjson v0.7.7
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/romshark/jscan v1.0.0
	github.com/sugawarayuuta/sonnet v0.0.0-20230425054915-e28ba49e3d17
)

require (
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

// replace to current repository version
replace github.com/go-faster/jx => ../

// CVE-2022-28948
replace gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0
