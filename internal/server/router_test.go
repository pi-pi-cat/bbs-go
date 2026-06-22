package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"bbs-go/internal/pkg/ginx"
	webspa "bbs-go/web"

	"github.com/gin-gonic/gin"
)

func TestLegacyAdminRouteFallsBackToSPAWithoutDashboardRedirect(t *testing.T) {
	app := newRouter()

	for _, path := range []string{"/admin", "/admin/users"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s status=%d want %d", path, rec.Code, http.StatusOK)
		}
		if got := rec.Header().Get("Location"); got != "" {
			t.Fatalf("%s Location=%q want empty", path, got)
		}
	}
}

func TestRouterNoRouteSeparatesAPIStaticAndSPA(t *testing.T) {
	app := newRouter()

	tests := []struct {
		method      string
		path        string
		wantStatus  int
		contentType string
	}{
		{method: http.MethodGet, path: "/api/not-exists", wantStatus: http.StatusNotFound, contentType: "application/json"},
		{method: http.MethodGet, path: "/api/admin/not-exists", wantStatus: http.StatusNotFound, contentType: "application/json"},
		{method: http.MethodPost, path: "/res/not-exists.png", wantStatus: http.StatusNotFound},
		{method: http.MethodGet, path: "/assets/not-exists.css", wantStatus: http.StatusOK, contentType: "text/html"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Fatalf("%s %s status=%d want %d", tt.method, tt.path, rec.Code, tt.wantStatus)
		}
		if tt.contentType != "" && !strings.Contains(rec.Header().Get("Content-Type"), tt.contentType) {
			t.Fatalf("%s %s Content-Type=%q want %q", tt.method, tt.path, rec.Header().Get("Content-Type"), tt.contentType)
		}
	}
}

func TestSPAHandlerFallsBackToEmbeddedFiles(t *testing.T) {
	handler := ginx.NewSPAHandler(filepath.Join(t.TempDir(), "missing"), webspa.SPA, "build/spa", ginx.DirOptions{
		ShowList:  false,
		SPA:       true,
		IndexName: "index.html",
	})

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/topic/123", nil)

	handler(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("Content-Type=%q want text/html", rec.Header().Get("Content-Type"))
	}
	if !strings.Contains(rec.Body.String(), "__reactRouterContext") {
		t.Fatalf("embedded SPA response does not look like React Router index.html")
	}
}

func TestRouterDoesNotUseRuntimeReflectionMVCMapping(t *testing.T) {
	source, err := os.ReadFile("router.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, forbidden := range []string{
		"reflect.",
		"ginx.Context",
		"jsonResult",
		"parseMVCMethodRoute",
		"registerController(",
		"buildControllerHandler",
		"Controller{Ctx:",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("router.go still contains runtime MVC mapping artifact %q", forbidden)
		}
	}

	if regexp.MustCompile(`\w+Group\.Any\(`).MatchString(text) {
		t.Fatal("router.go should register explicit HTTP methods instead of RouterGroup.Any")
	}
}

func TestControllersLayerHasBeenMigratedToHandlers(t *testing.T) {
	if _, err := os.Stat("../controllers"); !os.IsNotExist(err) {
		t.Fatalf("internal/controllers should not exist after Gin handler migration, stat err=%v", err)
	}

	for _, path := range []string{
		"../handlers/api/user_handlers.go",
		"../handlers/api/topic_handlers.go",
		"../handlers/admin/user_handlers.go",
		"../handlers/admin/role_handlers.go",
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected per-route handler file %s: %v", path, err)
		}
	}

	for _, path := range []string{
		"../handlers/api/handlers.go",
		"../handlers/admin/handlers.go",
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("generic wrapper file %s should not exist, stat err=%v", path, err)
		}
	}
}

func TestHandlersDoNotUseControllerStyleStructReceivers(t *testing.T) {
	files, err := filepath.Glob("../handlers/*/*_handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("expected handler files")
	}

	forbidden := regexp.MustCompile(`type\s+\w+Handler\s+struct|func\s+\(\w+\s+\*\w+Handler\)|&\w+Handler\{|\*ginx\.Context|ginx\.NewContext|func\s+[A-Z]\w*(Any|Get|Post|Delete)\w*\(ctx\s+\*gin\.Context\)|ginx\.WriteJSON\(ctx,\s*[a-z]\w*\(ginx\.NewContext\(ctx\)|ginx\.WriteBadRequest|ginx\.WriteStatusJSON|ginx\.WriteJSON\(ctx,\s*web\.Json|\.JsonResult\(\)`)
	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if forbidden.Match(source) {
			t.Fatalf("%s still uses controller-style handler struct receivers", file)
		}
	}
}

func TestGinRouterRegistersCompatibleAPIPaths(t *testing.T) {
	app := newRouter()
	routes := registeredRoutes(app)

	for _, want := range []string{
		http.MethodGet + " /api/topic/category_navs",
		http.MethodGet + " /api/login/wx_login_config",
		http.MethodPost + " /api/topic/accept_answer/:id",
		http.MethodPost + " /api/admin/topic/update_issue_status",
		http.MethodGet + " /api/admin/search/reindex/status",
		http.MethodPost + " /api/admin/search/reindex",
		http.MethodGet + " /api/admin/seo/sitemap/status",
		http.MethodPost + " /api/admin/seo/sitemap/generate",
		http.MethodPost + " /api/admin/role/list",
		http.MethodPost + " /api/admin/role/update_sort",
		http.MethodDelete + " /api/admin/topic/recommend",
	} {
		if _, ok := routes[want]; !ok {
			t.Fatalf("route %q is not registered", want)
		}
	}

	for _, notWant := range []string{
		http.MethodGet + " /api/search/reindex",
		http.MethodPost + " /api/search/reindex",
		http.MethodPut + " /api/admin/role/list",
		http.MethodDelete + " /api/admin/role/list",
	} {
		if _, ok := routes[notWant]; ok {
			t.Fatalf("route %q should not be registered", notWant)
		}
	}
}

func TestSimplifiedProductDoesNotRegisterRemovedCommunityGrowthAPIRoutes(t *testing.T) {
	app := newRouter()
	routes := registeredRoutes(app)

	for _, notWant := range []string{
		http.MethodGet + " /api/topic/favorite/:id",
		http.MethodPost + " /api/article/favorite/:id",
		http.MethodGet + " /api/user/favorites",
		http.MethodGet + " /api/user/msg_recent",
		http.MethodGet + " /api/user/messages",
		http.MethodGet + " /api/user/score_logs",
		http.MethodGet + " /api/user/score/rank",
		http.MethodPost + " /api/favorite/add",
		http.MethodPost + " /api/favorite/delete",
		http.MethodPost + " /api/checkin/checkin",
		http.MethodGet + " /api/checkin/checkin",
		http.MethodGet + " /api/checkin/rank",
		http.MethodGet + " /api/link/list",
		http.MethodGet + " /api/link/top_links",
		http.MethodGet + " /api/search/user",
		http.MethodPost + " /api/fans/follow",
		http.MethodPost + " /api/fans/unfollow",
		http.MethodGet + " /api/fans/is_followed",
		http.MethodGet + " /api/fans/fans",
		http.MethodGet + " /api/fans/followed",
		http.MethodGet + " /api/fans/recent/fans",
		http.MethodGet + " /api/fans/recent/follow",
		http.MethodGet + " /api/task/tasks",
		http.MethodGet + " /api/task/groups",
		http.MethodGet + " /api/badge/badges",
		http.MethodGet + " /api/admin/common/task_event_types",
		http.MethodPost + " /api/admin/favorite/list",
		http.MethodPost + " /api/admin/favorite/create",
		http.MethodPost + " /api/admin/favorite/update",
		http.MethodGet + " /api/admin/favorite/:id",
		http.MethodPost + " /api/admin/link/list",
		http.MethodPost + " /api/admin/link/create",
		http.MethodPost + " /api/admin/link/update",
		http.MethodPost + " /api/admin/link/delete",
		http.MethodPost + " /api/admin/link/update_sort",
		http.MethodGet + " /api/admin/link/:id",
		http.MethodPost + " /api/admin/user-score-log/list",
		http.MethodGet + " /api/admin/user-score-log/:id",
		http.MethodGet + " /api/admin/task-config/groups",
		http.MethodPost + " /api/admin/task-config/list",
		http.MethodPost + " /api/admin/task-config/create",
		http.MethodPost + " /api/admin/task-config/update",
		http.MethodPost + " /api/admin/task-config/delete",
		http.MethodPost + " /api/admin/task-config/update_sort",
		http.MethodGet + " /api/admin/task-config/:id",
		http.MethodGet + " /api/admin/badge/list",
		http.MethodPost + " /api/admin/badge/list",
		http.MethodPost + " /api/admin/badge/create",
		http.MethodPost + " /api/admin/badge/update",
		http.MethodPost + " /api/admin/badge/delete",
		http.MethodPost + " /api/admin/badge/update_sort",
		http.MethodGet + " /api/admin/badge/:id",
		http.MethodPost + " /api/admin/level-config/list",
		http.MethodPost + " /api/admin/level-config/save_all",
		http.MethodGet + " /api/admin/level-config/:id",
		http.MethodPost + " /api/admin/user-task-log/list",
		http.MethodGet + " /api/admin/user-task-log/:id",
		http.MethodPost + " /api/admin/user-exp-log/list",
		http.MethodGet + " /api/admin/user-exp-log/:id",
		http.MethodPost + " /api/admin/user-badge/list",
		http.MethodGet + " /api/admin/user-badge/:id",
	} {
		if _, ok := routes[notWant]; ok {
			t.Fatalf("route %q should not be registered in the simplified product", notWant)
		}
	}
}

func registeredRoutes(app *gin.Engine) map[string]struct{} {
	routes := map[string]struct{}{}
	for _, route := range app.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}
	return routes
}
