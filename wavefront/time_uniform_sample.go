package wavefront

import (
	"sync"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"istio.io/istio/pkg/log"
)

// TimeUniformSample is a Sample to store valuse only for a specific time
type TimeUniformSample struct {
	mutex    sync.Mutex
	values   Queue
	lifetime time.Duration
}

// list of all samples created to be cleaned periodically
var sampleList = make([]*TimeUniformSample, 0)

// recycle ticker that will fire a samples cleaning
var ticker *time.Ticker

// time delay for the recycle ticker
const recycleTime = time.Second * 5

// NewTimeUniformSample constructs and new TimeUniformSample with a specific lifetime for values and
// a fixed size
func NewTimeUniformSample(lifetime time.Duration, size int) metrics.Sample {
	sample := &TimeUniformSample{lifetime: lifetime, values: newQueue(size)}

	if ticker == nil {
		ticker = time.NewTicker(recycleTime)
		go func() {
			for t := range ticker.C {
				log.Infof("cleaning histograms old values - %s", t)
				for _, sample := range sampleList {
					sample.cleanOldValues()
				}
			}
		}()
	}

	sampleList = append(sampleList, sample)
	return sample
}

// Clear clears all samples.
func (s *TimeUniformSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values = Queue{}
}

// Count returns the number of samples recorded, which may exceed the
// reservoir size.
func (s *TimeUniformSample) Count() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return int64(s.values.len)
}

// Max returns the maximum value in the sample, which may not be the maximum
// value ever to be part of the sample.
func (s *TimeUniformSample) Max() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMax(s.rawValues())
}

// Mean returns the mean of the values in the sample.
func (s *TimeUniformSample) Mean() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMean(s.rawValues())
}

// Min returns the minimum value in the sample, which may not be the minimum
// value ever to be part of the sample.
func (s *TimeUniformSample) Min() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleMin(s.rawValues())
}

// Percentile returns an arbitrary percentile of values in the sample.
func (s *TimeUniformSample) Percentile(p float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentile(s.rawValues(), p)
}

// Percentiles returns a slice of arbitrary percentiles of values in the
// sample.
func (s *TimeUniformSample) Percentiles(ps []float64) []float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SamplePercentiles(s.rawValues(), ps)
}

// Size returns the size of the sample, which is at most the reservoir size.
func (s *TimeUniformSample) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.rawValues())
}

// Snapshot returns a read-only copy of the sample.
func (s *TimeUniformSample) Snapshot() metrics.Sample {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, s.values.len)
	copy(values, s.rawValues())
	return metrics.NewSampleSnapshot(int64(s.values.len), values)
}

// StdDev returns the standard deviation of the values in the sample.
func (s *TimeUniformSample) StdDev() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleStdDev(s.rawValues())
}

// Sum returns the sum of the values in the sample.
func (s *TimeUniformSample) Sum() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleSum(s.rawValues())
}

// Update samples a new value.
func (s *TimeUniformSample) Update(v int64) {
	s.cleanOldValues()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.values.IsFull() {
		s.values.Pop()
	}
	s.values.Push(&sampleValue{v: v, time: time.Now()})
}

// remove expired values
func (s *TimeUniformSample) cleanOldValues() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now()
	var needPop bool
	for more := true; more; more = needPop {
		sample, empty := s.values.Peek()
		if empty {
			needPop = false
		} else {
			needPop = (now.Sub(sample.(*sampleValue).time).Seconds() > s.lifetime.Seconds())
			if needPop {
				s.values.Pop()
			}
		}
	}
}

// Values returns a copy of the values in the sample.
func (s *TimeUniformSample) Values() []int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.rawValues()
}

// return a array with all values without any order
func (s *TimeUniformSample) rawValues() []int64 {
	values := make([]int64, s.values.len)
	idx := 0
	for _, value := range s.values.content {
		if value != nil {
			values[idx] = value.(*sampleValue).v
			idx++
		}
	}
	return values
}

// Variance returns the variance of the values in the sample.
func (s *TimeUniformSample) Variance() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return metrics.SampleVariance(s.rawValues())
}

// store values with its creation time
type sampleValue struct {
	v    int64
	time time.Time
}

// Construct a new Queue with a fiex size
func newQueue(size int) Queue {
	queue := Queue{content: make([]interface{}, size)}
	return queue
}

// Queue is FIFO queue
type Queue struct {
	content   []interface{}
	readHead  int
	writeHead int
	len       int
}

// IsFull return true if the queue is full, false otherwise
func (q *Queue) IsFull() bool {
	return q.len >= len(q.content)
}

// Push add a new element to the end of queue.
// return true if the element is added correctly, false otherwise
func (q *Queue) Push(e interface{}) bool {
	if q.len >= len(q.content) {
		return false
	}
	q.content[q.writeHead] = e
	q.writeHead = (q.writeHead + 1) % len(q.content)
	q.len++
	return true
}

// Pop return, and remove from the queue, the first queue element and true if there
// is at least one element
// return nil and false if the queue is empty
func (q *Queue) Pop() (interface{}, bool) {
	if q.len <= 0 {
		return nil, false
	}
	result := q.content[q.readHead]
	q.content[q.readHead] = nil
	q.readHead = (q.readHead + 1) % len(q.content)
	q.len--
	return result, true
}

// Peek return, without remove it from the queue, the first queue element and true if there
// is at least one element
// return nil and false if the queue is empty
func (q *Queue) Peek() (interface{}, bool) {
	if q.len <= 0 {
		return nil, true
	}
	result := q.content[q.readHead]
	return result, false
}
