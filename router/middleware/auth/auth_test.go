// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-vela/server/constants"
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

func TestToken_Retrieve_TokenHeader(t *testing.T) {
	// setup types
	want := "foobar"

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	request.Header.Set("Token", want)

	// run test
	got := RetrieveTokenHeader(request)

	if !strings.EqualFold(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestToken_Retrieve_TokenHeader_Empty(t *testing.T) {
	// setup types
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)

	// run test
	got := RetrieveTokenHeader(request)

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}

func TestToken_Retrieve_TokenHeader_And_Access(t *testing.T) {
	// setup types
	wantAccess := "foo"
	wantToken := "bar"

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wantAccess))
	request.Header.Set("Token", wantToken)

	gotAccess, err := RetrieveAccessToken(request)
	if err != nil {
		t.Errorf("Retrieve returned err: %v", err)
	}

	gotTkn := RetrieveTokenHeader(request)

	if !strings.EqualFold(gotAccess, wantAccess) {
		t.Errorf("Retrieve is %v, want %v", gotAccess, wantAccess)
	}

	if !strings.EqualFold(gotTkn, wantToken) {
		t.Errorf("Retrieve is %v, want %v", gotTkn, wantToken)
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
