package tests

import (
	"bytes"
	"fmt"
	"github.com/pubgo/g/logs"
	"testing"
)
import "github.com/RoaringBitmap/roaring"

func TestName(t *testing.T) {

	// example inspired by https://github.com/fzandona/goroar
	fmt.Println("==roaring==")
	rb1 := roaring.BitmapOf(10001, 1, 2, 3, 4, 5, 100, 1000)
	fmt.Println(rb1.String())
	rb1.Add(100000)
	rb1.Add(100000)
	fmt.Println(rb1.String())
	fmt.Println(rb1.Contains(100))
	fmt.Println(rb1.GetCardinality())
	fmt.Println(rb1.GetSerializedSizeInBytes())
	fmt.Println(rb1.GetSizeInBytes())
	fmt.Println(rb1.HasRunCompression())
	fmt.Println(rb1.IsEmpty())
	logs.P("Stats ", rb1.Stats())
	logs.P("Maximum ", rb1.Maximum())
	logs.P("Minimum ", rb1.Minimum())
	logs.P("Rank ", rb1.Rank(100))
	//logs.P("Rank ", rb1.RunOptimize)
	fmt.Println(rb1.Select(1))
	iter1 := rb1.Iterator()
	iter1.AdvanceIfNeeded(100)
	for iter1.HasNext() {
		fmt.Println(iter1.Next())
	}

	iter2 := rb1.ReverseIterator()
	for iter2.HasNext() {
		fmt.Println(iter2.Next())
	}

	rb2 := roaring.BitmapOf(3, 4, 1000)
	fmt.Println(rb2.String())

	rb3 := roaring.New()
	fmt.Println(rb3.String())

	fmt.Println("Cardinality: ", rb1.GetCardinality())

	fmt.Println("Contains 3? ", rb1.Contains(3))

	rb1.And(rb2)

	rb3.Add(1)
	rb3.Add(5)

	rb3.Or(rb1)

	// computes union of the three bitmaps in parallel using 4 workers
	roaring.ParOr(4, rb1, rb2, rb3)
	// computes intersection of the three bitmaps in parallel using 4 workers
	roaring.ParAnd(4, rb1, rb2, rb3)

	// prints 1, 3, 4, 5, 1000
	i := rb3.Iterator()
	for i.HasNext() {
		fmt.Println(i.Next())
	}
	fmt.Println()

	// next we include an example of serialization
	buf := new(bytes.Buffer)
	rb1.WriteTo(buf) // we omit error handling
	newrb := roaring.New()
	newrb.ReadFrom(buf)
	if rb1.Equals(newrb) {
		fmt.Println("I wrote the content to a byte stream and read it back.")
	}
	// you can iterate over bitmaps using ReverseIterator(), Iterator, ManyIterator()
}
