package app

import (
	"reflect"
	"regexp"
	"testing"
)

func responseEqual(lhs, rhs *evaluateResponse) bool {
	if lhs == nil {
		if rhs == nil {
			return true
		}
		return false
	}

	return reflect.DeepEqual(lhs.Responses, rhs.Responses)
}

func TestApiEvaluate(t *testing.T) {
	testCases := []struct {
		req  evaluateRequest
		resp evaluateResponse
		err  error
	}{
		{
			req: evaluateRequest{
				Commands: []string{"%let x = 2", "x * x * x"},
			},
			resp: evaluateResponse{
				Responses: []commandResponse{
					{Evaluated: true},
					{
						ReplOutput: "8 : int",
						Issue:      "",
						Evaluated:  true,
					},
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		s := &service{}
		resp, err := s.evaluate(&tc.req)
		if err != tc.err {
			t.Errorf("Expected err %v, got %v", tc.err, err)
		}

		if !responseEqual(resp, &tc.resp) {
			t.Errorf("Expected evaluate resp %v, got %v", tc.resp, resp)
		}
	}
}

func TestApiJson(t *testing.T) {
	testCases := []struct {
		req    string
		respRe *regexp.Regexp
		err    error
	}{
		{
			req:    `{"commands": ["%let x = 2", "x * x * x"]}`,
			respRe: regexp.MustCompile(`{"responses":[{"replOutput":"","issue":"","evaluated":true},{"replOutput":"8 : int","issue":"","evaluated":true}],"evalTime":\d+}`),
			err:    nil,
		},
	}

	for _, tc := range testCases {
		s := service{}
		resp, err := s.evaluateJSON([]byte(tc.req))
		if err != tc.err {
			t.Errorf("evaluate returned wanted %v, got %v", tc.err, err)
		}
		if tc.respRe.MatchString(string(resp)) {
			t.Errorf("evaluate response wanted %s, got %s", tc.respRe, string(resp))
		}
	}
}
