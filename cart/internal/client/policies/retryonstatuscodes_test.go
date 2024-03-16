//go:build unit
// +build unit

package policies

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestRetryOnStatusCodes_ShouldBeRetried(t *testing.T) {
	t.Parallel()
	policy := NewRetryOnStatusCodes([]int{420, 429}, 3)
	type args struct {
		attemptNum int
		in1        *http.Request
		resp       *http.Response
		respErr    error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"successResp",
			args{
				1,
				nil,
				&http.Response{StatusCode: 200},
				nil,
			},
			false,
		},
		{
			"badCode",
			args{
				1,
				nil,
				&http.Response{StatusCode: 429},
				nil,
			},
			true,
		},
		{
			"badCodeWithMaxAttempts",
			args{
				3,
				nil,
				&http.Response{StatusCode: 429},
				nil,
			},
			false,
		},
		{
			"returnedErr",
			args{
				1,
				nil,
				&http.Response{StatusCode: 429},
				fmt.Errorf("someerror"),
			},
			false,
		},
		{
			"responseISNil",
			args{
				1,
				nil,
				nil,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := policy.ShouldBeRetried(tt.args.attemptNum, tt.args.in1, tt.args.resp, tt.args.respErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
