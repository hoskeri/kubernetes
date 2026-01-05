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
	"fmt"
	"net/url"
	"strings"
)

type invalidUnixURLError struct {
	u string
}

func (e invalidUnixURLError) Error() string {
	return fmt.Sprintf("invalid url %q; must be of the form https+unix://ADDRESS[/PATH][?QUERY][#FRAGMENT]", e.u)
}

// ParseHTTPSUnixURI parses https+unix urls of the form "https+unix://ADDRESS[/PATH][?QUERY][#FRAGMENT]".
func ParseHTTPSUnixURI(s string) (*url.URL, error) {
	if !strings.HasPrefix(s, "https+unix:") {
		return url.Parse(s)
	}

	if !strings.HasPrefix(s, "https+unix://") {
		return nil, invalidUnixURLError{u: s}
	}

	rest := s[len("https+unix://"):]
	if rest == "" {
		return nil, invalidUnixURLError{u: s}
	}

	// The address is terminated by the first '/', '?', or '#' (or end of string).
	endAddr := len(rest)
	for i, c := range rest {
		if c == '/' || c == '?' || c == '#' {
			endAddr = i
			break
		}
	}

	rawAddress := rest[:endAddr]
	if rawAddress == "" {
		return nil, invalidUnixURLError{u: s}
	}

	address, err := url.PathUnescape(rawAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape address %q: %v", rawAddress, err)
	}

	if !strings.HasPrefix(address, "/") {
		return nil, fmt.Errorf("address %q must be an absolute path starting with '/'", address)
	}

	// Parse path, query, and fragment by using a dummy base URL.
	dummyURLStr := "https://dummy" + rest[endAddr:]
	dummyURL, err := url.Parse(dummyURLStr)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme:      "https+unix",
		Host:        address,
		Path:        dummyURL.Path,
		RawPath:     dummyURL.RawPath,
		ForceQuery:  dummyURL.ForceQuery,
		RawQuery:    dummyURL.RawQuery,
		Fragment:    dummyURL.Fragment,
		RawFragment: dummyURL.RawFragment,
	}

	return u, nil
}
