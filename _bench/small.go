/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bench

import (
	"github.com/go-faster/jx"
)

// "Small" case from sonic testdata:
// https://github.com/bytedance/sonic/blob/0f66ab721157029e48bbee72a06b3cf30bd536b1/testdata/small.go

//easyjson:json
type Small struct {
	BookId  int           `json:"id"`
	BookIds []int         `json:"ids"`
	Title   string        `json:"title"`
	Titles  []string      `json:"titles"`
	Price   float64       `json:"price"`
	Prices  []float64     `json:"prices"`
	Hot     bool          `json:"hot"`
	Hots    []bool        `json:"hots"`
	Author  SmallAuthor   `json:"author"`
	Authors []SmallAuthor `json:"authors"`
	Weights []int         `json:"weights"`
}

func (s Small) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("id")
	e.Int(s.BookId)

	e.FieldStart("ids")
	e.ArrStart()
	for _, v := range s.BookIds {
		e.Int(v)
	}
	e.ArrEnd()

	e.FieldStart("title")
	e.Str(s.Title)

	e.FieldStart("titles")
	e.ArrStart()
	for _, v := range s.Titles {
		e.Str(v)
	}
	e.ArrEnd()

	e.FieldStart("price")
	e.Float64(s.Price)

	e.FieldStart("prices")
	e.ArrStart()
	for _, v := range s.Prices {
		e.Float64(v)
	}
	e.ArrEnd()

	e.FieldStart("hot")
	e.Bool(s.Hot)

	e.FieldStart("hots")
	e.ArrStart()
	for _, v := range s.Hots {
		e.Bool(v)
	}
	e.ArrEnd()

	e.FieldStart("author")
	s.Author.Encode(e)

	e.FieldStart("authors")
	e.ArrStart()
	for _, v := range s.Authors {
		v.Encode(e)
	}
	e.ArrEnd()

	e.FieldStart("weights")
	if s.Weights == nil {
		e.Null()
	} else {
		e.ArrStart()
		for _, v := range s.Weights {
			e.Int(v)
		}
		e.ArrEnd()
	}

	e.ObjEnd()
}

//easyjson:json
type SmallAuthor struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Male bool   `json:"male"`
}

func (a SmallAuthor) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("name")
	e.Str(a.Name)
	e.FieldStart("age")
	e.Int(a.Age)
	e.FieldStart("male")
	e.Bool(a.Male)
	e.ObjEnd()
}
