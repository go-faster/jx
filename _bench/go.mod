module bench

go 1.17

require (
	github.com/bytedance/sonic v1.0.0-rc.1.0.20211109044705-0f66ab721157
	github.com/go-faster/jx v0.0.0-replaced
	github.com/json-iterator/go v1.1.12
	github.com/mailru/easyjson v0.7.7
)

require (
	github.com/chenzhuoyu/base64x v0.0.0-20211019084208-fb5309c8db06 // indirect
	github.com/go-faster/errors v0.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
)

// replace to current repository version
replace github.com/go-faster/jx => ../
