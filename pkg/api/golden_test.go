package api_test

// HTTP-level golden tests: drive both routers via httptest against a real
// seeded MySQL and snapshot the JSON responses. These freeze the API wire
// format during the refactor — any diff here is either a regression or an
// intentional API change to be reviewed together with the frontends.
//
// Run:
//   make db-up && make test-integration
// Update snapshots after an intentional API change:
//   UPDATE_GOLDEN=1 make test-integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	"github.com/homenoc/dsbd-backend/pkg/testutil"
)

const (
	seedMasterEmail = "master@example.com"
	seedPassword    = "password"
	// USER_TOKEN is normally generated client-side; a fixed value keeps runs deterministic.
	fixedUserToken = "golden-test-user-token"
)

var rfc3339Like = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)

// volatile keys whose values change between runs and are replaced before comparison.
var tokenKeys = map[string]bool{
	"access_token": true,
	"user_token":   true,
	"tmp_token":    true,
	"mail_token":   true,
	"token":        true,
	"debug":        true, // contains client IP / login user
}

type response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

func TestGoldenAPI(t *testing.T) {
	testutil.SetupIntegration(t)
	gin.SetMode(gin.TestMode)

	userRouter := api.NewUserRouter()
	adminRouter := api.NewAdminRouter()

	// --- login flows (tokens extracted from raw responses) ---
	adminToken := adminLogin(t, adminRouter)
	userAccessToken := userLogin(t, userRouter)

	userHeaders := map[string]string{
		"USER_TOKEN":   fixedUserToken,
		"ACCESS_TOKEN": userAccessToken,
	}
	adminHeaders := map[string]string{"ACCESS_TOKEN": adminToken}

	type testCase struct {
		name    string
		router  *gin.Engine
		method  string
		path    string
		headers map[string]string
		body    string
	}

	// NOTE: cases run in order and later cases depend on earlier mutations
	// (service add flips group add_allow, connection add flips service add_allow).
	cases := []testCase{
		// --- user API: reads & error shapes ---
		{name: "user_health", router: userRouter, method: "GET", path: "/health"},
		{name: "user_catalog", router: userRouter, method: "GET", path: "/api/v1/catalog", headers: userHeaders},
		{name: "user_info", router: userRouter, method: "GET", path: "/api/v1/info", headers: userHeaders},
		{name: "user_service_add_allow", router: userRouter, method: "GET", path: "/api/v1/service/add_allow", headers: userHeaders},
		{name: "user_info_unauthorized", router: userRouter, method: "GET", path: "/api/v1/info",
			headers: map[string]string{"USER_TOKEN": "bogus", "ACCESS_TOKEN": "bogus"}},
		{name: "user_login_wrong_pass", router: userRouter, method: "POST", path: "/api/v1/login",
			headers: map[string]string{"USER_TOKEN": fixedUserToken, "Email": seedMasterEmail, "HASH_PASS": "WRONG"}},

		// --- admin API: reads & error shapes ---
		{name: "admin_login_wrong_creds", router: adminRouter, method: "POST", path: "/api/v1/login",
			headers: map[string]string{"USER": "admin", "PASS": "wrong"}},
		{name: "admin_unauthorized", router: adminRouter, method: "GET", path: "/api/v1/user",
			headers: map[string]string{"ACCESS_TOKEN": "bogus"}},
		{name: "admin_catalog", router: adminRouter, method: "GET", path: "/api/v1/catalog", headers: adminHeaders},
		{name: "admin_user_list", router: adminRouter, method: "GET", path: "/api/v1/user", headers: adminHeaders},
		{name: "admin_user_detail", router: adminRouter, method: "GET", path: "/api/v1/user/1", headers: adminHeaders},
		{name: "admin_group_list", router: adminRouter, method: "GET", path: "/api/v1/group", headers: adminHeaders},
		{name: "admin_group_detail", router: adminRouter, method: "GET", path: "/api/v1/group/1", headers: adminHeaders},
		{name: "admin_service_list", router: adminRouter, method: "GET", path: "/api/v1/service", headers: adminHeaders},
		{name: "admin_service_detail", router: adminRouter, method: "GET", path: "/api/v1/service/1", headers: adminHeaders},
		{name: "admin_connection_list", router: adminRouter, method: "GET", path: "/api/v1/connection", headers: adminHeaders},
		{name: "admin_connection_detail", router: adminRouter, method: "GET", path: "/api/v1/connection/1", headers: adminHeaders},
		{name: "admin_noc_list", router: adminRouter, method: "GET", path: "/api/v1/noc", headers: adminHeaders},
		{name: "admin_noc_detail", router: adminRouter, method: "GET", path: "/api/v1/noc/1", headers: adminHeaders},
		{name: "admin_router_list", router: adminRouter, method: "GET", path: "/api/v1/router", headers: adminHeaders},
		{name: "admin_router_detail", router: adminRouter, method: "GET", path: "/api/v1/router/1", headers: adminHeaders},
		{name: "admin_gateway_list", router: adminRouter, method: "GET", path: "/api/v1/gateway", headers: adminHeaders},
		{name: "admin_gateway_ip_list", router: adminRouter, method: "GET", path: "/api/v1/gateway_ip", headers: adminHeaders},
		{name: "admin_notice_list", router: adminRouter, method: "GET", path: "/api/v1/notice", headers: adminHeaders},
		{name: "admin_support_list", router: adminRouter, method: "GET", path: "/api/v1/support", headers: adminHeaders},
		{name: "admin_noc_add_bad_body", router: adminRouter, method: "POST", path: "/api/v1/noc",
			headers: adminHeaders, body: `{invalid json`},

		// --- mutations: user happy path ---
		{name: "user_connection_add", router: userRouter, method: "POST", path: "/api/v1/service/1/connection",
			headers: userHeaders,
			body: `{"connection_type":"EIP","preferred_ap":"東日本","ntt":"etc",` +
				`"address":"東京都千代田区","ipv4_route":"IPv4 Full Route","ipv6_route":"IPv6 Full Route",` +
				`"term_ip":"203.0.113.100","monitor":true}`},
		{name: "user_service_add_transit", router: userRouter, method: "POST", path: "/api/v1/service",
			headers: userHeaders,
			body: `{"service_type":"IP3B","org":"テスト組織","org_en":"Test Org",` +
				`"postcode":"100-0001","address":"東京都千代田区","address_en":"Chiyoda, Tokyo",` +
				`"abuse":"abuse@example.com","avg_upstream":10,"max_upstream":100,` +
				`"avg_downstream":10,"max_downstream":100,"max_bandwidth_as":"100Mbps",` +
				`"start_date":"2030-01-01","asn":65001,"bgp_comment":"golden test"}`},
		{name: "user_service_add_denied_add_allow", router: userRouter, method: "POST", path: "/api/v1/service",
			headers: userHeaders,
			body: `{"service_type":"IP3B","org":"x","org_en":"x","postcode":"1","address":"x","address_en":"x",` +
				`"abuse":"a@example.com","avg_upstream":10,"max_upstream":100,"avg_downstream":10,` +
				`"max_downstream":100,"start_date":"2030-01-01","asn":65001}`},

		// --- mutations: admin approval flow ---
		{name: "admin_service_approve", router: adminRouter, method: "PUT", path: "/api/v1/service/2",
			headers: adminHeaders, body: `{"pass":true}`},
		{name: "admin_service_detail_after_approve", router: adminRouter, method: "GET", path: "/api/v1/service/2",
			headers: adminHeaders},
		{name: "admin_connection_open", router: adminRouter, method: "PUT", path: "/api/v1/connection/2",
			headers: adminHeaders, body: `{"open":true}`},
		{name: "admin_connection_detail_after_open", router: adminRouter, method: "GET", path: "/api/v1/connection/2",
			headers: adminHeaders},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, rawBody := doRequest(tc.router, tc.method, tc.path, tc.headers, tc.body)
			got := response{Status: status, Body: normalize(parseJSON(t, rawBody))}
			compareGolden(t, tc.name, got)
		})
	}
}

// adminLogin returns an admin access token via POST /api/v1/login (USER/PASS headers).
func adminLogin(t *testing.T, router *gin.Engine) string {
	t.Helper()
	status, body := doRequest(router, "POST", "/api/v1/login",
		map[string]string{"USER": "admin", "PASS": "admin"}, "")
	if status != http.StatusOK {
		t.Fatalf("admin login failed: status=%d body=%s", status, body)
	}
	var res struct {
		Token []struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil || len(res.Token) == 0 {
		t.Fatalf("admin login: unexpected response: %s", body)
	}
	return res.Token[0].AccessToken
}

// userLogin performs the two-step user login:
// GET /login (tmp token) → POST /login with HASH_PASS = sha256(storedPassHash + tmpToken).
func userLogin(t *testing.T, router *gin.Engine) string {
	t.Helper()
	status, body := doRequest(router, "GET", "/api/v1/login",
		map[string]string{"USER_TOKEN": fixedUserToken}, "")
	if status != http.StatusOK {
		t.Fatalf("user login init failed: status=%d body=%s", status, body)
	}
	var initRes struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal([]byte(body), &initRes); err != nil || initRes.Token == "" {
		t.Fatalf("user login init: unexpected response: %s", body)
	}

	storedPass := strings.ToLower(hash.Generate(seedPassword))
	hashPass := hash.Generate(storedPass + initRes.Token)

	status, body = doRequest(router, "POST", "/api/v1/login", map[string]string{
		"USER_TOKEN": fixedUserToken,
		"Email":      seedMasterEmail,
		"HASH_PASS":  hashPass,
	}, "")
	if status != http.StatusOK {
		t.Fatalf("user login failed: status=%d body=%s", status, body)
	}
	var res struct {
		Token []struct {
			AccessToken string `json:"access_token"`
		} `json:"token"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil || len(res.Token) == 0 {
		t.Fatalf("user login: unexpected response: %s", body)
	}
	return res.Token[0].AccessToken
}

func doRequest(router *gin.Engine, method, path string, headers map[string]string, body string) (int, string) {
	var reader *bytes.Reader
	if body != "" {
		reader = bytes.NewReader([]byte(body))
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func parseJSON(t *testing.T, raw string) any {
	t.Helper()
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		// non-JSON body (should not happen); keep raw for visibility
		return raw
	}
	return v
}

// normalize walks the decoded JSON and replaces volatile values:
// RFC3339-looking timestamps → "<TIME>", token/debug fields → "<VOLATILE>".
func normalize(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			if tokenKeys[strings.ToLower(k)] {
				if s, ok := item.(string); ok && s != "" {
					out[k] = "<VOLATILE>"
					continue
				}
			}
			out[k] = normalize(item)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = normalize(item)
		}
		return out
	case string:
		if rfc3339Like.MatchString(val) {
			return "<TIME>"
		}
		return val
	default:
		return v
	}
}

func compareGolden(t *testing.T, name string, got response) {
	t.Helper()
	goldenPath := filepath.Join("testdata", "golden", name+".json")

	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	gotJSON = append(gotJSON, '\n')

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("failed to create golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, gotJSON, 0o644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		t.Logf("updated golden: %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file missing (run UPDATE_GOLDEN=1 to create): %v", err)
	}
	if !bytes.Equal(want, gotJSON) {
		t.Errorf("response differs from golden %s\n--- want ---\n%s\n--- got ---\n%s",
			goldenPath, truncate(string(want)), truncate(string(gotJSON)))
	}
}

func truncate(s string) string {
	const max = 4000
	if len(s) <= max {
		return s
	}
	return s[:max] + fmt.Sprintf("\n... (%d bytes truncated)", len(s)-max)
}
