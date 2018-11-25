package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer_Start_ReturnsErrorIfTimerAlreadyStarted(t *testing.T) {
	tim := New()
	e := tim.Start("test1")
	assert.Nil(t, e)
	e = tim.Start("test2")
	assert.Nil(t, e)
}

func TestTimer_End_ReturnsErrorIfTimerAlreadyEnded(t *testing.T) {
	tim := New()
	e := tim.Start("test1")
	assert.Nil(t, e)
	e = tim.End("test1")
	assert.Nil(t, e)
	e = tim.End("test1")
	assert.NotNil(t, e)
}

func TestTimer_End_ReturnsErrorIfTimerNotStarted(t *testing.T) {
	tim := New()
	e := tim.End("test1")
	assert.NotNil(t, e)
}

func TestTimer_Elapsed_ReturnsErrorIfTimerHasNotStartedAndEnded(t *testing.T) {
	tim := New()
	_, e := tim.Elapsed("test1")
	assert.NotNil(t, e)
	e = tim.Start("test1")
	assert.Nil(t, e)
	_, e = tim.Elapsed("test1")
	assert.NotNil(t, e)
	e = tim.End("test1")
	assert.Nil(t, e)
	_, e = tim.Elapsed("test1")
	assert.Nil(t, e)
}

func TestTimer_Elapsed_ReturnsSaneValue(t *testing.T) {
	tim := New()
	_ = tim.Start("test1")
	time.Sleep(5 * time.Millisecond)
	_ = tim.End("test1")
	d, _ := tim.Elapsed("test1")
	fuzzyEqual(t, time.Duration(5), d, 2*time.Millisecond)
}

func TestTimer_Elapsed_TracksMultipleTimers(t *testing.T) {
	tim := New()
	_ = tim.Start("test1")
	_ = tim.Start("test2")
	_ = tim.Start("test3")
	time.Sleep(10 * time.Millisecond)
	_ = tim.End("test1")
	_ = tim.End("test2")
	_ = tim.End("test3")
	d1, _ := tim.Elapsed("test1")
	d2, _ := tim.Elapsed("test2")
	d3, _ := tim.Elapsed("test3")
	fuzzyEqualMulti(t, []time.Duration{10, 10, 10}, []time.Duration{d1, d2, d3}, 2*time.Millisecond)
}

// Mostly for interests sake, not for actual testing as it's somewhat variable how it prints out the timers
func TestTimer_String(t *testing.T) {
	tim := New()
	_ = tim.Start("test1")
	_ = tim.Start("test2")
	_ = tim.Start("test3")
	time.Sleep(10 * time.Millisecond)
	_ = tim.End("test1")
	_ = tim.End("test2")
	//_ = tim.End("test3")
	//fmt.Print(tim)
	//time.Sleep(150 * time.Millisecond)
	_ = tim.End("test3")
	fmt.Print(tim)
}

func fuzzyEqualMulti(t *testing.T, expected []time.Duration, actual []time.Duration, fuzz time.Duration) {
	if len(expected) != len(actual) {
		assert.Fail(t, "expected and actual are not the same length")
	}
	for i, ex := range expected {
		ax := actual[i]
		fuzzyEqual(t, ex, ax, fuzz)
		return
	}
}

func fuzzyEqual(t *testing.T, ex time.Duration, ax time.Duration, fuzz time.Duration) {
	if !(ax >= ex-fuzz && ax <= ex+fuzz) {
		assert.Failf(t, "expected and actual were not fuzzily equal", "%d %d", ex, ax)
	}
}
