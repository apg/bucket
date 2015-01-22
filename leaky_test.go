package bucket

import (
	"fmt"
	"testing"
	"time"
)

func TestLeakyCapacity(t *testing.T) {
	bucket := NewLeakyBucket(10, time.Millisecond)
	if err := bucket.Put(1); err != nil {
		t.Error("Expected to be able to put an item in a bucket with capacity of 1")
	}

	select {
	case <-bucket.C():
	case <-time.After(20 * time.Millisecond):
		t.Error("Time out of 2 seconds before a delivery happened.")
	}

	select {
	case <-bucket.C():
		t.Error("Received something from bucket even though it should be empty")
	case <-time.After(2 * time.Millisecond):
	}
}

func ExampleLeakyBucket() {
	var errors int
	bucket := NewLeakyBucket(2, 10*time.Millisecond)

	// consumer, at a steady pace
	go func() {
		for v := range bucket.C() {
			fmt.Printf("%v\n", v)
		}
	}()

	// produce a 100 values at somewhat random intervals.
	for i := 0; i < 10; i++ {
		if err := bucket.Put(i); err != nil {
			fmt.Printf("Couldn't Put %v, because bucket is full\n", i)
			errors++
		}

		time.Sleep(3 * time.Millisecond)
	}

	fmt.Printf("Couldn't Put %d values\n", errors)
	// Output: Couldn't Put 3, because bucket is full
	// 0
	// Couldn't Put 5, because bucket is full
	// Couldn't Put 6, because bucket is full
	// 1
	// Couldn't Put 8, because bucket is full
	// Couldn't Put 9, because bucket is full
	// 2
	// Couldn't Put 5 values
}
