package permissions

import "testing"

func TestAdminPermissionRegistryMatchesMethodAndPath(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/topic/list")
	if !ok {
		t.Fatalf("expected topic list permission to be registered")
	}
	if code != PermissionTopicView.Code {
		t.Fatalf("expected %s, got %s", PermissionTopicView.Code, code)
	}
}

func TestAdminPermissionRegistryProtectsIssueStatusUpdate(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/topic/update_issue_status")
	if !ok {
		t.Fatalf("expected issue status update permission to be registered")
	}
	if code != PermissionTopicSolve.Code {
		t.Fatalf("expected %s, got %s", PermissionTopicSolve.Code, code)
	}
}

func TestAdminPermissionRegistryAllowsRoleOptionsFromUserUpdate(t *testing.T) {
	codes, ok := GetAdminPermissionCodes("GET", "/api/admin/role/roles")
	if !ok {
		t.Fatalf("expected role options permission to be registered")
	}
	expected := []string{PermissionRoleView.Code, PermissionUserUpdate.Code}
	if len(codes) != len(expected) {
		t.Fatalf("expected %#v, got %#v", expected, codes)
	}
	for i, expectedCode := range expected {
		if codes[i] != expectedCode {
			t.Fatalf("expected %#v, got %#v", expected, codes)
		}
	}
}

func TestAdminPermissionRegistryAllowsEitherUserForbiddenPermission(t *testing.T) {
	codes, ok := GetAdminPermissionCodes("POST", "/api/admin/user/forbidden")
	if !ok {
		t.Fatalf("expected user forbidden permission to be registered")
	}
	expected := []string{PermissionUserForbidden.Code, PermissionUserForbiddenForever.Code}
	if len(codes) != len(expected) {
		t.Fatalf("expected %#v, got %#v", expected, codes)
	}
	for i, expectedCode := range expected {
		if codes[i] != expectedCode {
			t.Fatalf("expected %#v, got %#v", expected, codes)
		}
	}
}

func TestAdminPermissionRegistryProtectsForbiddenWordDelete(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/forbidden-word/delete")
	if !ok {
		t.Fatalf("expected forbidden word delete permission to be registered")
	}
	if code != PermissionForbiddenWordDelete.Code {
		t.Fatalf("expected %s, got %s", PermissionForbiddenWordDelete.Code, code)
	}
}

func TestAdminPermissionRegistryProtectsUserReportAudit(t *testing.T) {
	code, ok := GetAdminPermissionCode("POST", "/api/admin/user-report/audit")
	if !ok {
		t.Fatalf("expected user report audit permission to be registered")
	}
	if code != PermissionUserReportAudit.Code {
		t.Fatalf("expected %s, got %s", PermissionUserReportAudit.Code, code)
	}
}

func TestAdminPermissionRegistryRejectsRemovedGrowthAndLinkPaths(t *testing.T) {
	paths := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/api/admin/link/1"},
		{method: "POST", path: "/api/admin/link/list"},
		{method: "POST", path: "/api/admin/link/delete"},
		{method: "GET", path: "/api/admin/badge/1"},
		{method: "POST", path: "/api/admin/badge/list"},
		{method: "POST", path: "/api/admin/badge/update_sort"},
		{method: "GET", path: "/api/admin/level-config/1"},
		{method: "POST", path: "/api/admin/level-config/list"},
		{method: "GET", path: "/api/admin/task-config/1"},
		{method: "POST", path: "/api/admin/task-config/list"},
		{method: "POST", path: "/api/admin/task-config/update_sort"},
		{method: "GET", path: "/api/admin/user-task-log/1"},
		{method: "POST", path: "/api/admin/user-task-log/list"},
		{method: "GET", path: "/api/admin/user-exp-log/1"},
		{method: "POST", path: "/api/admin/user-exp-log/list"},
		{method: "GET", path: "/api/admin/user-badge/1"},
		{method: "POST", path: "/api/admin/user-badge/list"},
	}

	for _, path := range paths {
		if code, ok := GetAdminPermissionCode(path.method, path.path); ok {
			t.Fatalf("expected %s %s to be rejected, got %s", path.method, path.path, code)
		}
	}
}

func TestAdminPermissionRegistryRejectsUnknownAdminPath(t *testing.T) {
	if code, ok := GetAdminPermissionCode("POST", "/api/admin/unknown/action"); ok {
		t.Fatalf("expected unknown admin path to be rejected, got %s", code)
	}
}

func TestAdminPermissionRegistryRejectsCommentManagementPaths(t *testing.T) {
	paths := []struct {
		method string
		path   string
	}{
		{method: "GET", path: "/api/admin/comment/1"},
		{method: "POST", path: "/api/admin/comment/list"},
		{method: "DELETE", path: "/api/admin/comment/1"},
	}

	for _, path := range paths {
		if code, ok := GetAdminPermissionCode(path.method, path.path); ok {
			t.Fatalf("expected %s %s to be rejected, got %s", path.method, path.path, code)
		}
	}
}
