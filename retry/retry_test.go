package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSimpleDelayStrategy(t *testing.T) {
	delay := time.Millisecond * 10
	retryStrat := SimpleDelayStrategy(delay)

	const expectedIterations = 15

	var lastTry time.Time
	var count int
	var times [][]time.Time
	retryStrat.Try(func() (success bool) {
		count++
		now := time.Now()
		if lastTry.Unix() > 0 {
			times = append(times, []time.Time{lastTry, now})
		}
		lastTry = now
		return count == expectedIterations
	})

	durations := lo.Map[[]time.Time, time.Duration](times, func(item []time.Time, index int) time.Duration {
		return item[1].Sub(item[0])
	})

	assert.Equal(t, expectedIterations, count)
	assert.Len(t, durations, expectedIterations-1)
	for _, dur := range durations {
		assert.InEpsilon(t, delay.Microseconds(), dur.Microseconds(), 0.20)
	}
}

func TestTry(t *testing.T) {
	for _, val := range []bool{true, false} {
		t.Run(fmt.Sprintf("when function returns %v", val), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var called bool

			fakeResult := gofakeit.SentenceSimple()
			fakeFunc := func() (string, bool) {
				called = true
				return fakeResult, val
			}

			mockStrategy := NewMockStrategy(ctrl)
			mockStrategy.EXPECT().
				Try(gomock.Any()).
				DoAndReturn(func(fn func() bool) {
					resultOfFn := fn()
					assert.Equal(t, val, resultOfFn)
				})

			result := Try(mockStrategy, fakeFunc)
			assert.Equal(t, fakeResult, result)
			assert.True(t, called)
		})
	}
}

func TestTryE(t *testing.T) {
	for _, val := range []bool{true, false} {
		t.Run(fmt.Sprintf("when function returns %v", val), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var called bool

			fakeResult := gofakeit.SentenceSimple()
			fakeError := errors.New(gofakeit.SentenceSimple())
			fakeFunc := func() (string, error, bool) {
				called = true
				return fakeResult, fakeError, val
			}

			mockStrategy := NewMockStrategy(ctrl)
			mockStrategy.EXPECT().
				Try(gomock.Any()).
				DoAndReturn(func(fn func() bool) {
					resultOfFn := fn()
					assert.Equal(t, val, resultOfFn)
				})

			result, err := TryE(mockStrategy, fakeFunc)
			assert.Equal(t, fakeResult, result)
			assert.Equal(t, fakeError, err)
			assert.True(t, called)
		})
	}
}
