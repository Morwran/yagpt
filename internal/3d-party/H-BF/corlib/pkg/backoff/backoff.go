package backoff

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type (
	//Backoff is alias to backoff.BackOff
	Backoff = backoff.BackOff

	//BackOffContext is alias to backoff.BackOffContext
	BackOffContext = backoff.BackOffContext //nolint:revive
)

var (
	//Stop alias
	Stop = backoff.Stop

	//StopBackoff stop backoff
	StopBackoff backoff.StopBackOff

	//ZeroBackoff zero backoff
	ZeroBackoff backoff.ZeroBackOff

	//WithContext backoff with context
	WithContext = backoff.WithContext
)

//NewConstantBackOff ...
func NewConstantBackOff(d time.Duration) Backoff {
	return &backoff.ConstantBackOff{Interval: d}
}

//ExponentialBackoffBuilder exponential backoff builder
func ExponentialBackoffBuilder() exponentialBackoffBuilder { //nolint:revive
	return exponentialBackoffBuilder{
		inner: backoff.NewExponentialBackOff(),
	}
}

var (
	_ = ZeroBackoff
	_ = StopBackoff
	_ = Stop
	_ = NewConstantBackOff
	_ = ExponentialBackoffBuilder
	_ = WithContext
)

type exponentialBackoffBuilder struct {
	inner *backoff.ExponentialBackOff
}

//Build build exponential backoff
func (eb exponentialBackoffBuilder) Build() Backoff {
	ret := new(backoff.ExponentialBackOff)
	*ret = *eb.inner
	return ret
}

//WithRandomizationFactor ...
func (eb exponentialBackoffBuilder) WithRandomizationFactor(d float64) exponentialBackoffBuilder {
	if d >= 0 {
		eb.inner.RandomizationFactor = d
	}
	return eb
}

//WithInitialInterval ...
func (eb exponentialBackoffBuilder) WithInitialInterval(d time.Duration) exponentialBackoffBuilder {
	if d >= 0 {
		eb.inner.InitialInterval = d
	}
	return eb
}

//WithMultiplier ...
func (eb exponentialBackoffBuilder) WithMultiplier(d float64) exponentialBackoffBuilder {
	if d >= 0 {
		eb.inner.Multiplier = d
	}
	return eb
}

//WithMaxInterval ...
func (eb exponentialBackoffBuilder) WithMaxInterval(d time.Duration) exponentialBackoffBuilder {
	if d >= 0 {
		eb.inner.MaxInterval = d
	}
	return eb
}

//WithMaxElapsedThreshold ...
func (eb exponentialBackoffBuilder) WithMaxElapsedThreshold(d time.Duration) exponentialBackoffBuilder {
	if d >= 0 {
		eb.inner.MaxElapsedTime = d
	}
	return eb
}
