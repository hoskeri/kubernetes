/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package net

import (
	"testing"
)

func TestParseHTTPSUnixURI(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedHost string
		expectedPath string
		expectedRawQ string
		expectedFrag string
		expectErr    bool
	}{
		{
			name:         "valid socket and path",
			input:        "https+unix://%2Ftmp%2Fkube.sock/api/v1/pods",
			expectedHost: "/tmp/kube.sock",
			expectedPath: "/api/v1/pods",
		},
		{
			name:         "valid socket no path",
			input:        "https+unix://%2Ftmp%2Fkube.sock",
			expectedHost: "/tmp/kube.sock",
		},
		{
			name:         "standard https url remains untouched",
			input:        "https://example.com/api/v1",
			expectedHost: "example.com",
			expectedPath: "/api/v1",
		},
		{
			name:      "invalid triple slash",
			input:     "https+unix:///tmp/kube.sock",
			expectErr: true,
		},
		{
			name:      "missing host and double slash",
			input:     "https+unix://",
			expectErr: true,
		},
		{
			name:      "non-absolute socket path",
			input:     "https+unix://tmp%2Fkube.sock",
			expectErr: true,
		},
		{
			name:      "invalid scheme prefix",
			input:     "https+unix:/tmp/kube.sock",
			expectErr: true,
		},
		{
			name:         "with query parameters",
			input:        "https+unix://%2Ftmp%2Fkube.sock/api/v1/pods?watch=true&timeout=5m",
			expectedHost: "/tmp/kube.sock",
			expectedPath: "/api/v1/pods",
			expectedRawQ: "watch=true&timeout=5m",
		},
		{
			name:         "with fragment",
			input:        "https+unix://%2Ftmp%2Fkube.sock/api/v1/pods#section1",
			expectedHost: "/tmp/kube.sock",
			expectedPath: "/api/v1/pods",
			expectedFrag: "section1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ParseHTTPSUnixURI(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if u.Host != tc.expectedHost {
				t.Errorf("expected host %q, got %q", tc.expectedHost, u.Host)
			}
			if u.Path != tc.expectedPath {
				t.Errorf("expected path %q, got %q", tc.expectedPath, u.Path)
			}
			if u.RawQuery != tc.expectedRawQ {
				t.Errorf("expected RawQuery %q, got %q", tc.expectedRawQ, u.RawQuery)
			}
			if u.Fragment != tc.expectedFrag {
				t.Errorf("expected Fragment %q, got %q", tc.expectedFrag, u.Fragment)
			}
		})
	}
}
