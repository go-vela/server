// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/go-vela/server/internal"
)

func TestMiddleware_NoCache(t *testing.T) {
	// setup types
	wantCacheControl := "no-cache, no-store, max-age=0, must-revalidate, value"
	wantExpires := "Thu, 01 Jan 1970 00:00:00 GMT"
	wantLastModified := time.Now().UTC().Format(http.TimeFormat)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(NoCache)
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotCacheControl := context.Writer.Header().Get("Cache-Control")
	gotExpires := context.Writer.Header().Get("Expires")
	gotLastModified := context.Writer.Header().Get("Last-Modified")

	if resp.Code != http.StatusOK {
		t.Errorf("NoCache returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotCacheControl, wantCacheControl) {
		t.Errorf("NoCache Cache-Control is %v, want %v", gotCacheControl, wantCacheControl)
	}

	if !reflect.DeepEqual(gotExpires, wantExpires) {
		t.Errorf("NoCache Expires is %v, want %v", gotExpires, wantExpires)
	}

	if !reflect.DeepEqual(gotLastModified, wantLastModified) {
		t.Errorf("NoCache Last-Modified is %v, want %v", gotLastModified, wantLastModified)
	}
}

func TestMiddleware_Options(t *testing.T) {
	// setup tests
	tests := []struct {
		name            string
		webAddress      string
		requestMethod   string
		wantStatusCode  int
		wantOrigin      string
		wantMethods     string
		wantHeaders     string
		wantAllow       string
		wantContentType string
		wantCredentials string
	}{
		{
			name:            "without web address",
			webAddress:      "",
			requestMethod:   http.MethodOptions,
			wantStatusCode:  http.StatusOK,
			wantOrigin:      "*",
			wantMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			wantHeaders:     "authorization, origin, content-type, accept",
			wantAllow:       "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS",
			wantContentType: "application/json",
			wantCredentials: "",
		},
		{
			name:            "with web address",
			webAddress:      "http://localhost:8888",
			requestMethod:   http.MethodOptions,
			wantStatusCode:  http.StatusOK,
			wantOrigin:      "http://localhost:8888",
			wantMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			wantHeaders:     "authorization, origin, content-type, accept",
			wantAllow:       "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS",
			wantContentType: "application/json",
			wantCredentials: "true",
		},
		{
			name:            "not OPTIONS request",
			webAddress:      "http://localhost:8888",
			requestMethod:   http.MethodGet,
			wantStatusCode:  http.StatusNotFound,
			wantOrigin:      "",
			wantMethods:     "",
			wantHeaders:     "",
			wantAllow:       "",
			wantContentType: "text/plain",
			wantCredentials: "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup context
			gin.SetMode(gin.TestMode)

			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(test.requestMethod, "/health", nil)

			// setup mock server
			m := &internal.Metadata{
				Vela: &internal.Vela{
					Address:    "http://localhost:8080",
					WebAddress: test.webAddress,
				},
			}
			engine.Use(Metadata(m))
			engine.Use(Options)
			engine.OPTIONS("/health", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			gotOrigin := context.Writer.Header().Get("Access-Control-Allow-Origin")
			gotMethods := context.Writer.Header().Get("Access-Control-Allow-Methods")
			gotHeaders := context.Writer.Header().Get("Access-Control-Allow-Headers")
			gotCredentialsHeaders := context.Writer.Header().Get("Access-Control-Allow-Credentials")
			gotAllow := context.Writer.Header().Get("Allow")
			gotContentType := context.Writer.Header().Get("Content-Type")

			if resp.Code != test.wantStatusCode {
				t.Errorf("Options returned %v, want %v", resp.Code, http.StatusOK)
			}

			if gotOrigin != test.wantOrigin {
				t.Errorf("Options Access-Control-Allow-Origin is %v, want %v", gotOrigin, test.wantOrigin)
			}

			if gotMethods != test.wantMethods {
				t.Errorf("Options Access-Control-Allow-Methods is %v, want %v", gotMethods, test.wantMethods)
			}

			if gotHeaders != test.wantHeaders {
				t.Errorf("Options Access-Control-Allow-Headers is %v, want %v", gotHeaders, test.wantHeaders)
			}

			if gotCredentialsHeaders != test.wantCredentials {
				t.Errorf("Options Access-Control-Allow-Credentials is %v, want %v", gotCredentialsHeaders, test.wantCredentials)
			}

			if gotAllow != test.wantAllow {
				t.Errorf("Options Allow is %v, want %v", gotAllow, test.wantAllow)
			}

			if gotContentType != test.wantContentType {
				t.Errorf("Options Content-Type is %v, want %v", gotContentType, test.wantContentType)
			}
		})
	}
}

func TestMiddleware_Cors(t *testing.T) {
	// setup tests
	tests := []struct {
		name              string
		webAddress        string
		wantOrigin        string
		wantExposeHeaders string
		wantCredentials   string
	}{
		{
			name:              "without web address",
			webAddress:        "",
			wantOrigin:        "*",
			wantExposeHeaders: "link, x-total-count",
			wantCredentials:   "",
		},
		{
			name:              "with web address",
			webAddress:        "http://localhost:8888",
			wantOrigin:        "http://localhost:8888",
			wantExposeHeaders: "link, x-total-count",
			wantCredentials:   "true",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			m := &internal.Metadata{
				Vela: &internal.Vela{
					Address:    "http://localhost:8080",
					WebAddress: test.webAddress,
				},
			}

			// setup context
			gin.SetMode(gin.TestMode)

			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

			// setup mock server
			engine.Use(Metadata(m))
			engine.Use(Cors)
			engine.GET("/health", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			gotOrigin := context.Writer.Header().Get("Access-Control-Allow-Origin")
			gotExposeHeaders := context.Writer.Header().Get("Access-Control-Expose-Headers")
			gotCredentialsHeaders := context.Writer.Header().Get("Access-Control-Allow-Credentials")

			if resp.Code != http.StatusOK {
				t.Errorf("CORS returned %v, want %v", resp.Code, http.StatusOK)
			}

			if gotOrigin != test.wantOrigin {
				t.Errorf("CORS Access-Control-Allow-Origin is %v, want %v", gotOrigin, test.wantOrigin)
			}

			if gotExposeHeaders != test.wantExposeHeaders {
				t.Errorf("CORS Access-Control-Expose-Headers is %v, want %v", gotExposeHeaders, test.wantExposeHeaders)
			}

			if gotCredentialsHeaders != test.wantCredentials {
				t.Errorf("CORS Access-Control-Allow-Credentials is %v, want %v", gotCredentialsHeaders, test.wantCredentials)
			}
		})
	}
}

func TestMiddleware_Secure(t *testing.T) {
	tests := []struct {
		name   string
		useTLS bool
	}{
		{
			name:   "with tls",
			useTLS: true,
		},
		{
			name:   "without tls",
			useTLS: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup types
			wantFrameOptions := "DENY"
			wantContentTypeOptions := "nosniff"
			wantProtection := "1; mode=block"
			wantSecurity := "max-age=63072000; includeSubDomains; preload"

			// setup context
			gin.SetMode(gin.TestMode)

			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)
			if test.useTLS {
				context.Request.TLS = new(tls.ConnectionState)
			}

			// setup mock server
			engine.Use(Secure)
			engine.GET("/health", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// run test
			engine.ServeHTTP(context.Writer, context.Request)

			gotFrameOptions := context.Writer.Header().Get("X-Frame-Options")
			gotContentTypeOptions := context.Writer.Header().Get("X-Content-Type-Options")
			gotProtection := context.Writer.Header().Get("X-XSS-Protection")
			gotSecurity := context.Writer.Header().Get("Strict-Transport-Security")

			if resp.Code != http.StatusOK {
				t.Errorf("Secure returned %v, want %v", resp.Code, http.StatusOK)
			}

			if gotFrameOptions != wantFrameOptions {
				t.Errorf("Secure X-Frame-Options is %v, want %v", gotFrameOptions, wantFrameOptions)
			}

			if gotContentTypeOptions != wantContentTypeOptions {
				t.Errorf("Secure X-Content-Type-Options is %v, want %v", gotContentTypeOptions, wantContentTypeOptions)
			}

			if gotProtection != wantProtection {
				t.Errorf("Secure X-XSS-Protection is %v, want %v", gotProtection, wantProtection)
			}

			if gotSecurity != wantSecurity {
				t.Errorf("Secure Strict-Transport-Security is %v, want %v", gotSecurity, wantSecurity)
			}
		})
	}
}

func TestMiddleware_RequestID(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(RequestID)
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotRequestID := context.Writer.Header().Get("X-Request-ID")
	_, err := uuid.Parse(gotRequestID)

	if resp.Code != http.StatusOK {
		t.Errorf("RequestID returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("RequestID X-Request-ID is not a valid UUID: %v", err)
	}
}
