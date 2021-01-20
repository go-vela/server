// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestToken_Compose(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	d := time.Minute * 5
	now := time.Now()
	exp := now.Add(d)

	claims := &Claims{
		IsActive: u.GetActive(),
		IsAdmin:  u.GetAdmin(),
		StandardClaims: jwt.StandardClaims{
			Subject:   u.GetName(),
			IssuedAt:  now.Unix(),
			ExpiresAt: exp.Unix(),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	want, err := tkn.SignedString([]byte(u.GetHash()))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	m := &types.Metadata{
		Vela: &types.Vela{
			AccessTokenDuration: d,
		},
	}

	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)
	context.Set("securecookie", false)

	// run test
	_, got, err := Compose(context, u)
	if err != nil {
		t.Errorf("Compose returned err: %v", err)
	}

	if !strings.EqualFold(got, want) {
		t.Errorf("Compose is %v, want %v", got, want)
	}
}

func TestToken_Parse(t *testing.T) {
	// setup types
	want := new(library.User)
	want.SetID(1)
	want.SetName("foo")
	want.SetRefreshToken("fresh")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetActive(false)
	want.SetAdmin(false)
	want.SetFavorites([]string{})

	m := &types.Metadata{
		Vela: &types.Vela{
			AccessTokenDuration: time.Minute * 5,
		},
	}

	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)

	tkn, err := CreateAccessToken(want, time.Minute*5)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(want)

	// run test
	got, err := Parse(tkn, db)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestToken_Parse_Error_NoParse(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// run test
	got, err := Parse("!@#$%^&*()", db)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestToken_Parse_Error_InvalidSignature(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	claims := &Claims{
		IsActive: u.GetActive(),
		IsAdmin:  u.GetAdmin(),
		StandardClaims: jwt.StandardClaims{
			Subject: u.GetName(),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := tkn.SignedString([]byte(u.GetHash()))
	if err != nil {
		t.Errorf("Unable to create test token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// run test
	got, err := Parse(token, db)
	if err == nil {
		t.Errorf("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestToken_Parse_AccessToken_Expired(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	tkn, err := CreateAccessToken(u, time.Minute*-1)
	if err != nil {
		t.Errorf("Unable to create token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// run test
	_, err = Parse(tkn, db)
	if err == nil {
		t.Errorf("Parse should return error due to expiration")
	}
}

func TestToken_Refresh(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	d := time.Minute * 5

	m := &types.Metadata{
		Vela: &types.Vela{
			AccessTokenDuration: d,
		},
	}

	rt, _, err := CreateRefreshToken(u, d)
	if err != nil {
		t.Errorf("unable to create refresh token")
	}

	u.SetRefreshToken(rt)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// set up context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)
	context.Set("database", db)

	// run tests
	got, err := Refresh(context, rt)
	if err != nil {
		t.Error("Refresh should not error")
	}

	if len(got) == 0 {
		t.Errorf("Refresh should have returned an access token")
	}
}

func TestToken_Refresh_Expired(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	d := time.Minute * -1

	m := &types.Metadata{
		Vela: &types.Vela{
			AccessTokenDuration: d,
		},
	}

	rt, _, err := CreateRefreshToken(u, d)
	if err != nil {
		t.Errorf("unable to create refresh token")
	}

	u.SetRefreshToken(rt)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// set up context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)
	context.Set("database", db)

	// run tests
	_, err = Refresh(context, rt)
	if err == nil {
		t.Error("Refresh with expired token should error")
	}
}

func TestToken_Refresh_TokenMissing(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")

	d := time.Minute * -1

	m := &types.Metadata{
		Vela: &types.Vela{
			AccessTokenDuration: d,
		},
	}

	rt, _, err := CreateRefreshToken(u, d)
	if err != nil {
		t.Errorf("unable to create refresh token")
	}

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()

	_ = db.CreateUser(u)

	// set up context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	context.Set("metadata", m)
	context.Set("database", db)

	// run tests
	_, err = Refresh(context, rt)
	if err == nil {
		t.Error("Refresh with token that doesn't exist in database should error")
	}
}

func TestToken_Retrieve_Refresh(t *testing.T) {
	// setup types
	want := "fresh"

	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
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
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)
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
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

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
	request, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// run test
	got, err := RetrieveRefreshToken(request)
	if err == nil {
		t.Errorf("Retrieve should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Retrieve is %v, want \"\"", got)
	}
}
