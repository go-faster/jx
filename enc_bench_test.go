package jx

import (
	"io"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkEncoderBigObject(b *testing.B) {
	b.ReportAllocs()

	e := GetEncoder()
	encodeObject(e)
	b.SetBytes(int64(len(e.Bytes())))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Reset()
		encodeObject(e)
	}
}

func encodeObject(w *Encoder) {
	w.ObjStart()

	w.FieldStart("objectId")
	w.UInt64(8838243212)

	w.FieldStart("name")
	w.Str("Jane Doe")

	w.FieldStart("address")
	w.ObjStart()
	for _, field := range addressFields {
		w.FieldStart(field.key)
		w.Str(field.val)
	}

	w.FieldStart("geo")
	{
		w.ObjStart()
		w.FieldStart("latitude")
		w.Float64(-154.550817)
		w.FieldStart("longitude")
		w.Float64(-84.176159)
		w.ObjEnd()
	}
	w.ObjEnd()

	w.FieldStart("specialties")
	w.ArrStart()
	for _, s := range specialties {
		w.Str(s)
	}
	w.ArrEnd()

	for i, text := range longText {
		w.FieldStart("longText" + strconv.Itoa(i))
		w.Str(text)
	}

	for i := 0; i < 25; i++ {
		num := i * 18328
		w.FieldStart("integerField" + strconv.Itoa(i))
		w.Int64(int64(num))
	}

	w.ObjEnd()
}

type field struct{ key, val string }

var (
	addressFields = []field{
		{"address1", "123 Example St"},
		{"address2", "Apartment 5D, Suite 3"},
		{"city", "Miami"},
		{"state", "FL"},
		{"postalCode", "33133"},
		{"country", "US"},
	}
	specialties = []string{
		"Web Design",
		"Go Programming",
		"Tennis",
		"Cycling",
		"Mixed martial arts",
	}
	longText = []string{
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
		`Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?`,
		`But I must explain to you how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system, and expound the actual teachings of the great explorer of the truth, the master-builder of human happiness. No one rejects, dislikes, or avoids pleasure itself, because it is pleasure, but because those who do not know how to pursue pleasure rationally encounter consequences that are extremely painful. Nor again is there anyone who loves or pursues or desires to obtain pain of itself, because it is pain, but because occasionally circumstances occur in which toil and pain can procure him some great pleasure. To take a trivial example, which of us ever undertakes laborious physical exercise, except to obtain some advantage from it? But who has any right to find fault with a man who chooses to enjoy a pleasure that has no annoying consequences, or one who avoids a pain that produces no resultant pleasure?`,
		`At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint occaecati cupiditate non provident, similique sunt in culpa qui officia deserunt mollitia animi, id est laborum et dolorum fuga. Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat facere possimus, omnis voluptas assumenda est, omnis dolor repellendus. Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.`,
		`On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammeled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains.`,
	}
)

func encodeFloats(enc *Encoder, arr []float64) {
	enc.ArrStart()
	for _, num := range arr {
		enc.Float64(num)
	}
	enc.ArrEnd()
}

func BenchmarkEncodeFloats(b *testing.B) {
	const N = 100_000
	arr := make([]float64, N)
	for i := 0; i < N; i++ {
		arr[i] = rand.NormFloat64()
	}
	size := func() int64 {
		var enc Encoder
		encodeFloats(&enc, arr)
		return int64(len(enc.Bytes()))
	}()
	b.Logf("Size: %d bytes", size)

	b.Run("Buffered", func(b *testing.B) {
		b.SetBytes(size)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Notice: no buffer reuse.
				var enc Encoder
				encodeFloats(&enc, arr)
			}
		})
	})
	b.Run("Stream", func(b *testing.B) {
		b.SetBytes(size)
		b.RunParallel(func(pb *testing.PB) {
			enc := NewStreamingEncoder(io.Discard, -1)
			for pb.Next() {
				enc.ResetWriter(io.Discard)
				encodeFloats(enc, arr)
				_ = enc.Close()
			}
		})
	})
}
