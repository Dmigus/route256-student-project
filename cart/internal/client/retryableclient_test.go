package client

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_retryRoundTripper_RoundTripOnAttempts(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	policy := NewRetryPolicyMock(mc)
	rt := NewHttpRoundTripper(mc)

	req := &http.Request{}
	resp := &http.Response{}
	policy.ShouldBeRetriedMock.Set(func(attempts int, req *http.Request, resp *http.Response, respErr error) (b1 bool) {
		if attempts < 3 {
			return true
		}
		return false
	})
	rt.RoundTripMock.Return(resp, nil)

	rrt := retryRoundTripper{next: rt, policy: policy}
	responseFromRoundTrip, err := rrt.RoundTrip(req)
	require.NoError(t, err, "unexpected err")
	assert.True(t, resp == responseFromRoundTrip, "response from retryRT differ then source")
	numberOfCallsRT := len(rt.RoundTripMock.Calls())
	assert.Equal(t, 3, numberOfCallsRT, "roundTripper called an unexpected number of times")
}

func Test_retryRoundTripper_RoundTripOnContext(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	policy := NewRetryPolicyMock(mc)
	rt := NewHttpRoundTripper(mc)

	ctx, cancel := context.WithCancel(context.Background())
	req := &http.Request{}
	req = req.WithContext(ctx)
	resp := &http.Response{}
	policy.ShouldBeRetriedMock.Set(func(attempts int, req *http.Request, resp *http.Response, respErr error) (b1 bool) {
		if attempts >= 10 {
			cancel()
		}
		return true
	})
	rt.RoundTripMock.Return(resp, nil)

	rrt := retryRoundTripper{next: rt, policy: policy}
	_, err := rrt.RoundTrip(req)
	require.ErrorIs(t, err, context.Canceled, "unexpected err")
	numberOfCallsRT := len(rt.RoundTripMock.Calls())
	assert.Equal(t, 10, numberOfCallsRT, "roundTripper called an unexpected number of times")
}

func Test_contextWasDone(t *testing.T) {
	t.Parallel()
	activeContext := context.Background()
	doneContext, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"active context",
			args{activeContext},
			false,
		},
		{
			"cancelled context",
			args{doneContext},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := contextWasDone(tt.args.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
