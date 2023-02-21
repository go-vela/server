// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestToken_Retrieve_Refresh(t *testing.T) {
	// setup types
	want := "fresh"

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	request.AddCookie(&http.Cookie{
		Name:  constants.RefreshTokenName,
		Value: want,
	})

	// run test
	got, err := RetrieveRefreshToken(request)
	if err != nil {
		t.Errorf("Retrieve returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestToken_Retrieve_Access(t *testing.T) {
	// setup types
	want := "foobar"

	header := fmt.Sprintf("Bearer %s", want)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", header)

	// run test
	got, err := RetrieveAccessToken(request)
	if err != nil {
		t.Errorf("Retrieve returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestToken_Retrieve_Access_Error(t *testing.T) {
	// setup types
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)

	// run test
	got, err := RetrieveAccessToken(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}

func TestToken_Retrieve_Refresh_Error(t *testing.T) {
	// setup types
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)

	// run test
	got, err := RetrieveRefreshToken(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}
