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
	// setup types
	wantOrigin := "*"
	wantMethods := "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	wantHeaders := "authorization, origin, content-type, accept"
	wantAllow := "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"
	wantContentType := "application/json"
	m := &internal.Metadata{
		Vela: &internal.Vela{
			Address: "http://localhost:8080",
		},
	}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodOptions, "/health", nil)

	// setup mock server
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
	gotAllow := context.Writer.Header().Get("Allow")
	gotContentType := context.Writer.Header().Get("Content-Type")

	if resp.Code != http.StatusOK {
		t.Errorf("Options returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotOrigin, wantOrigin) {
		t.Errorf("Options Access-Control-Allow-Origin is %v, want %v", gotOrigin, wantOrigin)
	}

	if !reflect.DeepEqual(gotMethods, wantMethods) {
		t.Errorf("Options Access-Control-Allow-Methods is %v, want %v", gotMethods, wantMethods)
	}

	if !reflect.DeepEqual(gotHeaders, wantHeaders) {
		t.Errorf("Options Access-Control-Allow-Headers is %v, want %v", gotHeaders, wantHeaders)
	}

	if !reflect.DeepEqual(gotAllow, wantAllow) {
		t.Errorf("Options Allow is %v, want %v", gotAllow, wantAllow)
	}

	if !reflect.DeepEqual(gotContentType, wantContentType) {
		t.Errorf("Options Content-Type is %v, want %v", gotContentType, wantContentType)
	}
}

func TestMiddleware_Options_InvalidMethod(t *testing.T) {
	// setup types
	m := &internal.Metadata{
		Vela: &internal.Vela{
			Address: "http://localhost:8080",
		},
	}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Metadata(m))
	engine.Use(Options)
	engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	gotOrigin := context.Writer.Header().Get("Access-Control-Allow-Origin")
	gotMethods := context.Writer.Header().Get("Access-Control-Allow-Methods")
	gotHeaders := context.Writer.Header().Get("Access-Control-Allow-Headers")
	gotAllow := context.Writer.Header().Get("Allow")
	gotContentType := context.Writer.Header().Get("Content-Type")

	if resp.Code != http.StatusOK {
		t.Errorf("Options returned %v, want %v", resp.Code, http.StatusOK)
	}

	if len(gotOrigin) > 0 {
		t.Errorf("Options Access-Control-Allow-Origin is %v, want \"\"", gotOrigin)
	}

	if len(gotMethods) > 0 {
		t.Errorf("Options Access-Control-Allow-Methods is %v, want \"\"", gotMethods)
	}

	if len(gotHeaders) > 0 {
		t.Errorf("Options Access-Control-Allow-Headers is %v, want \"\"", gotHeaders)
	}

	if len(gotAllow) > 0 {
		t.Errorf("Options Allow is %v, want \"\"", gotAllow)
	}

	if len(gotContentType) > 0 {
		t.Errorf("Options Content-Type is %v, want \"\"", gotContentType)
	}
}

func TestMiddleware_Cors(t *testing.T) {
	tests := []struct {
		name                  string
		m                     *internal.Metadata
		origin                string
		expectedOrigin        string
		expectedCredentials   string
		expectedExposeHeaders string
	}{
		{
			name: "*",
			m: &internal.Metadata{
				Vela: &internal.Vela{
					Address:          "http://localhost:8080",
					CorsAllowOrigins: []string{},
				},
			},
			origin:                "http://localhost:8888",
			expectedOrigin:        "*",
			expectedCredentials:   "",
			expectedExposeHeaders: "link, x-total-count",
		},
		{
			name: "WebAddress is origin",
			m: &internal.Metadata{
				Vela: &internal.Vela{
					WebAddress:       "http://localhost:8888",
					CorsAllowOrigins: []string{},
				},
			},
			origin:                "http://localhost:8888",
			expectedOrigin:        "http://localhost:8888",
			expectedCredentials:   "true",
			expectedExposeHeaders: "link, x-total-count",
		},
		{
			name: "CORSAllowOrigins origin is web address",
			m: &internal.Metadata{
				Vela: &internal.Vela{
					WebAddress:       "http://localhost:8888",
					CorsAllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
				},
			},
			origin:                "http://localhost:8888",
			expectedOrigin:        "http://localhost:8888",
			expectedCredentials:   "true",
			expectedExposeHeaders: "link, x-total-count",
		},
		{
			name: "CORSAllowOrigins origin is in list",
			m: &internal.Metadata{
				Vela: &internal.Vela{
					WebAddress:       "",
					CorsAllowOrigins: []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8888"},
				},
			},
			origin:                "http://localhost:8888",
			expectedOrigin:        "http://localhost:8888",
			expectedCredentials:   "true",
			expectedExposeHeaders: "link, x-total-count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			resp := httptest.NewRecorder()
			context, engine := gin.CreateTestContext(resp)
			context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)
			context.Request.Header.Add("Origin", tt.origin)

			// inject metadata
			engine.Use(func(c *gin.Context) {
				c.Set("metadata", tt.m)
				c.Next()
			})
			engine.Use(Cors)
			engine.GET("/health", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			engine.ServeHTTP(context.Writer, context.Request)

			gotOrigin := context.Writer.Header().Get("Access-Control-Allow-Origin")
			if gotOrigin != tt.expectedOrigin {
				t.Errorf("Access-Control-Allow-Origin is %v; want %v", gotOrigin, tt.expectedOrigin)
			}

			gotCredentials := context.Writer.Header().Get("Access-Control-Allow-Credentials")
			if gotCredentials != tt.expectedCredentials {
				t.Errorf("Access-Control-Allow-Credentials is %v; want %v", gotCredentials, tt.expectedCredentials)
			}

			gotExposeHeaders := context.Writer.Header().Get("Access-Control-Expose-Headers")
			if gotExposeHeaders != tt.expectedExposeHeaders {
				t.Errorf("Access-Control-Expose-Headers is %v; want %v", gotExposeHeaders, tt.expectedExposeHeaders)
			}
		})
	}
}

func TestMiddleware_Secure(t *testing.T) {
	// setup types
	wantFrameOptions := "DENY"
	wantContentTypeOptions := "nosniff"
	wantProtection := "1; mode=block"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

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

	if resp.Code != http.StatusOK {
		t.Errorf("Secure returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(gotFrameOptions, wantFrameOptions) {
		t.Errorf("Secure X-Frame-Options is %v, want %v", gotFrameOptions, wantFrameOptions)
	}

	if !reflect.DeepEqual(gotContentTypeOptions, wantContentTypeOptions) {
		t.Errorf("Secure X-Content-Type-Options is %v, want %v", gotContentTypeOptions, wantContentTypeOptions)
	}

	if !reflect.DeepEqual(gotProtection, wantProtection) {
		t.Errorf("Secure X-XSS-Protection is %v, want %v", gotProtection, wantProtection)
	}
}

func TestMiddleware_Secure_TLS(t *testing.T) {
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
	context.Request.TLS = new(tls.ConnectionState)

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

	if !reflect.DeepEqual(gotFrameOptions, wantFrameOptions) {
		t.Errorf("Secure X-Frame-Options is %v, want %v", gotFrameOptions, wantFrameOptions)
	}

	if !reflect.DeepEqual(gotContentTypeOptions, wantContentTypeOptions) {
		t.Errorf("Secure X-Content-Type-Options is %v, want %v", gotContentTypeOptions, wantContentTypeOptions)
	}

	if !reflect.DeepEqual(gotProtection, wantProtection) {
		t.Errorf("Secure X-XSS-Protection is %v, want %v", gotProtection, wantProtection)
	}

	if !reflect.DeepEqual(gotSecurity, wantSecurity) {
		t.Errorf("Secure Strict-Transport-Security is %v, want %v", gotSecurity, wantSecurity)
	}
}
