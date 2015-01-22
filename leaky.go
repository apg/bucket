package bucket

import (
	"container/list"
	"sync"
	"time"
)

type LeakyBucket struct {
	c        chan interface{}
	capacity int
	fn       func()

	mu  *sync.Mutex
	els *list.List
}

func NewLeakyBucket(capacity int, rate time.Duration) Bucket {
	c := make(chan interface{})
	bucket := &LeakyBucket{
		c:        c,
		capacity: capacity,
		mu:       new(sync.Mutex),
		els:      new(list.List),
	}

	// TODO: Allow for some sort of Pause()
	fn := func() {
		bucket.mu.Lock()
		if bucket.els.Len() > 0 {
			select {
			case c <- bucket.els.Remove(bucket.els.Front()):
			default:
			}
		}
		bucket.mu.Unlock()
		time.AfterFunc(rate, bucket.fn)
	}
	bucket.fn = fn
	bucket.fn()

	return bucket
}

// Put adds an element to the bucket, or spills it if the bucket is full.
func (b *LeakyBucket) Put(v interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.els.Len() <= b.capacity {
		b.els.PushBack(v)
		return nil
	} else {
		return ErrFull
	}
}

func (b *LeakyBucket) C() chan interface{} {
	return b.c
}
