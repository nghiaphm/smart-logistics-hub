package integration

import (
	"context"
	"testing"
)

// TASK-096: bảng workspace_users — chèn 1 dòng is_admin=true (kèm FK
// workspace/user/created_by) rồi đọc lại đúng.
func TestWorkspaceUsersInsertReadIsAdmin(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	wsRes, err := testDB.ExecContext(ctx,
		"INSERT INTO workspaces (workspace_code, name) VALUES (?, ?)",
		"WS-MEM-001", "Workspace membership test")
	if err != nil {
		t.Fatalf("insert workspace: %v", err)
	}
	workspaceID, _ := wsRes.LastInsertId()
	if workspaceID == 0 {
		t.Fatal("insert workspace did not set ID")
	}

	userRes, err := testDB.ExecContext(ctx,
		"INSERT INTO users (username, full_name, role) VALUES (?, ?, ?)",
		"ws-admin-001", "Workspace Admin", "admin")
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	userID, _ := userRes.LastInsertId()
	if userID == 0 {
		t.Fatal("insert user did not set ID")
	}

	rowRes, err := testDB.ExecContext(ctx,
		"INSERT INTO workspace_users (workspace_id, user_id, is_admin, created_by, updated_by) VALUES (?, ?, ?, ?, ?)",
		workspaceID, userID, true, userID, userID)
	if err != nil {
		t.Fatalf("insert workspace_users: %v", err)
	}
	rowID, _ := rowRes.LastInsertId()
	if rowID == 0 {
		t.Fatal("insert workspace_users did not set ID")
	}

	var (
		gotWorkspaceID int64
		gotUserID      int64
		gotCreatedBy   int64
		gotIsAdmin     int64
		gotIsActive    int64
	)
	err = testDB.QueryRowContext(ctx,
		"SELECT workspace_id, user_id, created_by, is_admin, is_active FROM workspace_users WHERE id = ?",
		rowID).Scan(&gotWorkspaceID, &gotUserID, &gotCreatedBy, &gotIsAdmin, &gotIsActive)
	if err != nil {
		t.Fatalf("read workspace_users: %v", err)
	}
	if gotWorkspaceID != workspaceID {
		t.Errorf("workspace_id = %d, want %d", gotWorkspaceID, workspaceID)
	}
	if gotUserID != userID {
		t.Errorf("user_id = %d, want %d", gotUserID, userID)
	}
	if gotCreatedBy != userID {
		t.Errorf("created_by = %d, want %d", gotCreatedBy, userID)
	}
	if gotIsAdmin != 1 {
		t.Errorf("is_admin = %d, want 1", gotIsAdmin)
	}
	if gotIsActive != 1 {
		t.Errorf("is_active = %d, want 1 (default)", gotIsActive)
	}
}

// TASK-096: UNIQUE (workspace_id, user_id) — chèn trùng user vào cùng
// workspace phải bị chặn.
func TestWorkspaceUsersUniqueWorkspaceUser(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()

	if _, err := testDB.ExecContext(ctx,
		"INSERT INTO workspaces (workspace_code, name) VALUES (?, ?)",
		"WS-MEM-002", "Unique test"); err != nil {
		t.Fatalf("insert workspace: %v", err)
	}
	if _, err := testDB.ExecContext(ctx,
		"INSERT INTO users (username) VALUES (?)", "ws-user-002"); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	stmt := `INSERT INTO workspace_users (workspace_id, user_id) VALUES (1, 1)`
	if _, err := testDB.ExecContext(ctx, stmt); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if _, err := testDB.ExecContext(ctx, stmt); err == nil {
		t.Error("second insert with same (workspace_id, user_id) should fail")
	}
}
