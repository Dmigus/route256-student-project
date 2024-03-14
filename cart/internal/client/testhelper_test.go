package client

import (
	"github.com/gojuno/minimock/v3"
	"testing"
)

type testHelper struct {
	policyMock *mRetryPolicyMockShouldBeRetried
	nextRTMock *mHttpRoundTripperRoundTrip
	rTripper   *retryRoundTripper
}

func newTestHelper(t *testing.T) testHelper {
	mc := minimock.NewController(t)
	helper := testHelper{}
	policy := NewRetryPolicyMock(mc)
	nextRT := NewHttpRoundTripper(mc)

	helper.policyMock = &(policy.ShouldBeRetriedMock)
	helper.nextRTMock = &(nextRT.RoundTripMock)
	helper.rTripper = &retryRoundTripper{next: nextRT, policy: policy}
	return helper
}
