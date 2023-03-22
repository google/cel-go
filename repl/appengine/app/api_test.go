package app

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func (r *evaluateResponse) getResponses() []commandResponse {
	if r == nil {
		return nil
	}

	return r.Responses
}

func generateRange(n int, toString func(int) string) []string {
	var result []string
	for i := 0; i < n; i++ {
		result = append(result, toString(i))
	}
	return result
}

func errAgrees(err error, msg string) bool {
	if err != nil && msg != "" {
		return strings.Contains(err.Error(), msg)
	}

	if (err == nil) != (msg == "") {
		return false
	}

	return true
}

func TestAPIEvaluate(t *testing.T) {
	testCases := []struct {
		service *service
		req     evaluateRequest
		resp    *evaluateResponse
		err     string
	}{
		{
			req: evaluateRequest{
				Commands: []string{"%let x = 2", "x * x * x"},
			},
			resp: &evaluateResponse{
				Responses: []commandResponse{
					{Evaluated: true},
					{
						ReplOutput: "8 : int",
						Issue:      "",
						Evaluated:  true,
					},
				},
			},
		},
		{
			req: evaluateRequest{
				Commands: []string{"%non_command", "}"},
			},
			resp: &evaluateResponse{
				Responses: []commandResponse{
					{
						ReplOutput: "",
						Issue:      "unsupported command: non_command",
						Evaluated:  true,
					},
					{
						ReplOutput: "",
						Issue:      "(1:0) no viable alternative at input '}'",
						Evaluated:  false,
					},
				},
			},
		},
		{
			service: &service{commandCountLimit: 5},
			req: evaluateRequest{
				Commands: generateRange(5, func(i int) string { return fmt.Sprintf("%%let x%d = %d", i, i) }),
			},
			resp: &evaluateResponse{
				Responses: []commandResponse{
					{Evaluated: true},
					{Evaluated: true},
					{Evaluated: true},
					{Evaluated: true},
					{Evaluated: true},
				},
			},
		},
		{
			service: &service{commandCountLimit: 4},
			req: evaluateRequest{
				Commands: generateRange(5, func(i int) string { return fmt.Sprintf("%%let x%d = %d", i, i) }),
			},
			resp: nil,
			err:  "number of commands (5) exceeded limit (4)",
		},
	}

	for _, tc := range testCases {
		s := &service{}
		if tc.service != nil {
			s = tc.service
		}
		resp, err := s.evaluate(&tc.req)
		if !errAgrees(err, tc.err) {
			t.Errorf("Expected err '%v', got '%v'", tc.err, err)
		}

		if !reflect.DeepEqual(resp.getResponses(), tc.resp.getResponses()) {
			t.Errorf("Expected evaluate resp %v, got %v", tc.resp, resp)
		}
	}
}

func TestAPIJSON(t *testing.T) {
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
