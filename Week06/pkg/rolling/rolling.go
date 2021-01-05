package rolling

import (
	"sync"
	"time"
)

/*
* Bucket duration is 1 second, reference implementation: https://github.com/afex/hystrix-go
* It should be sufficient for most use cases.
 */
//Number is the struct to store the statics within a window
type Number struct {
	WindowSize uint8
	Buckets    map[int64]*float64
	Mutex      *sync.RWMutex
}

//NewNumber is the constructor for Number
func NewNumber(windowSize uint8) *Number {
	r := &Number{
		WindowSize: windowSize,
		Buckets:    make(map[int64]*float64),
		Mutex:      &sync.RWMutex{},
	}
	return r
}

func (r *Number) getCurrentBucket() *float64 {
	now := time.Now().Unix()
	var bucket *float64
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		ptrNum := float64(0)
		bucket = &ptrNum
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r *Number) removeOldBuckets() {
	oldestBucket := time.Now().Unix() - int64(r.WindowSize)

	for timestamp := range r.Buckets {
		if timestamp <= oldestBucket {
			delete(r.Buckets, timestamp)
		}
	}
}

//Increment is
func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	*b += i
	r.removeOldBuckets()
}

//Sum is
func (r *Number) Sum(now time.Time) float64 {
	sum := float64(0)

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-int64(r.WindowSize) {
			sum += *bucket
		}
	}

	return sum
}

//Avg is
func (r *Number) Avg(now time.Time) float64 {
	return r.Sum(now) / float64(r.WindowSize)
}

// Max returns the maximum value seen in the last WindowSize seconds.
func (r *Number) Max(now time.Time) float64 {
	var max float64

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-int64(r.WindowSize) {
			if *bucket > max {
				max = *bucket
			}
		}
	}

	return max
}

// UpdateMax updates the maximum value in the current bucket.
func (r *Number) UpdateMax(val float64) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	if val > *b {
		*b = val
	}
	r.removeOldBuckets()
}
