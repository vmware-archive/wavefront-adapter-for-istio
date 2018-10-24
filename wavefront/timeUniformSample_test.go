package wavefront

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeUniformSample(t *testing.T) {
	sample := NewTimeUniformSample(time.Second)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)

	assertEqual(t, sample.Count(), int64(5), "error count")
	assertEqual(t, sample.Mean(), float64(1), "error mean")

	time.Sleep(2 * time.Second)

	sample.Update(2)
	sample.Update(1)
	assertEqual(t, sample.Count(), int64(2), "error count")
	assertEqual(t, sample.Mean(), float64(1.5), "error mean")
}

func TestTimeUniformSampleAutoClean(t *testing.T) {
	sample := NewTimeUniformSample(time.Second)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)

	assertEqual(t, sample.Count(), int64(5), "error count")
	assertEqual(t, sample.Mean(), float64(1), "error mean")

	// wait for autoclean (after 5 seconds)
	time.Sleep(6 * time.Second)

	assertEqual(t, sample.Count(), int64(0), "error count")
	assertEqual(t, sample.Mean(), float64(0), "error mean")
}

func TestTimeUniformSampleNoAutoClean(t *testing.T) {
	sample := NewTimeUniformSample(time.Minute)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)
	sample.Update(1)

	assertEqual(t, sample.Count(), int64(5), "error count")
	assertEqual(t, sample.Mean(), float64(1), "error mean")

	// wait for autoclean (after 5 seconds)
	time.Sleep(6 * time.Second)

	sample.Update(1)
	sample.Update(1)
	assertEqual(t, sample.Count(), int64(7), "error count")
	assertEqual(t, sample.Mean(), float64(1), "error mean")
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		t.Logf("OK - %v == %v", a, b)
		return
	}
	message = fmt.Sprintf("%v -- %v != %v", message, a, b)
	t.Fatal(message)
}
