package jir

import (
	"encoding/json"
	"testing"
)

type TestResp struct {
	Code uint64
}

func Benchmark_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		iter := ParseBytes(Default, input)
		if err := iter.ObjectBytes(func(i *Iterator, key []byte) error {
			switch string(key) {
			case "code":
				v, err := iter.Uint64()
				result.Code = v
				return err
			default:
				return iter.Skip()
			}
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_std_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		_ = json.Unmarshal(input, &result)
	}
}
