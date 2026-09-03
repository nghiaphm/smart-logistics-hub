package integration

import (
	"context"
	"testing"

	"my-web-app.com/smart-logistic-hub/internal/workspace/repository"
)

// TASK-097: gán is_admin cho 1 user trong workspace A → IsWorkspaceAdmin
// trả true cho A, false cho workspace B; thu hồi → false.
func TestWorkspaceMemberAdminRepo(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	// workspace A + B
	for _, code := range []string{"WS-MA-001", "WS-MA-002"} {
		if _, err := testDB.ExecContext(ctx,
			"INSERT INTO workspaces (workspace_code, name) VALUES (?, ?)", code, code); err != nil {
			t.Fatalf("insert workspace %s: %v", code, err)
		}
	}
	// user
	if _, err := testDB.ExecContext(ctx,
		"INSERT INTO users (username) VALUES (?)", "member-admin-001"); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	members := repository.NewMemberRepository(testDB)

	// gán admin cho user (workspace_id=1, user_id=1 — bảng vừa truncate)
	if err := members.Upsert(ctx, 1, 1, true, nil); err != nil {
		t.Fatalf("Upsert(is_admin=true) error = %v", err)
	}
	wu, err := members.GetByWorkspaceAndUser(ctx, 1, 1)
	if err != nil {
		t.Fatalf("GetByWorkspaceAndUser() error = %v", err)
	}
	if !wu.IsAdmin {
		t.Error("GetByWorkspaceAndUser() IsAdmin = false, want true")
	}

	adminA, err := members.IsWorkspaceAdmin(ctx, 1, 1)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin(A) error = %v", err)
	}
	if !adminA {
		t.Error("IsWorkspaceAdmin(user, workspace A) = false, want true")
	}

	adminB, err := members.IsWorkspaceAdmin(ctx, 1, 2)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin(B) error = %v", err)
	}
	if adminB {
		t.Error("IsWorkspaceAdmin(user, workspace B) = true, want false")
	}

	// thu hồi (is_admin=false) → IsWorkspaceAdmin false
	if err := members.Upsert(ctx, 1, 1, false, nil); err != nil {
		t.Fatalf("Upsert(is_admin=false) error = %v", err)
	}
	wu, err = members.GetByWorkspaceAndUser(ctx, 1, 1)
	if err != nil {
		t.Fatalf("GetByWorkspaceAndUser() after revoke error = %v", err)
	}
	if wu.IsAdmin {
		t.Error("IsAdmin still true after revoke")
	}
	adminAfter, err := members.IsWorkspaceAdmin(ctx, 1, 1)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin() after revoke error = %v", err)
	}
	if adminAfter {
		t.Error("IsWorkspaceAdmin = true after revoke, want false")
	}
}

// User không có membership trong workspace → IsWorkspaceAdmin false (không lỗi).
func TestWorkspaceMemberIsAdminNoRow(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	members := repository.NewMemberRepository(testDB)
	ok, err := members.IsWorkspaceAdmin(ctx, 999, 999)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin() error = %v", err)
	}
	if ok {
		t.Error("IsWorkspaceAdmin() = true for non-member, want false")
	}
}
