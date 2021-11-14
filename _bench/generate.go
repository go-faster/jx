package bench

//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers hello_world.go
//go:generate go run github.com/mailru/easyjson/easyjson -no_std_marshalers small.go
//go:generate go run github.com/pquerna/ffjson -w hello_world_ffjson_gen.go hello_world_ffjson.go
