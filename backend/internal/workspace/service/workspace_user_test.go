package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	userentity "my-web-app.com/smart-logistic-hub/internal/user/entity"
	wksentity "my-web-app.com/smart-logistic-hub/internal/workspace/entity"
	"my-web-app.com/smart-logistic-hub/internal/workspace/service"
)

type mockMemberRepo struct {
	upsertFn        func(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error
	getByFn         func(ctx context.Context, workspaceID, userID int64) (*wksentity.WorkspaceUser, error)
	isAdminFn       func(ctx context.Context, userID, workspaceID int64) (bool, error)
}

func (m *mockMemberRepo) Upsert(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, workspaceID, userID, isAdmin, actorID)
	}
	return nil
}

func (m *mockMemberRepo) GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID int64) (*wksentity.WorkspaceUser, error) {
	if m.getByFn != nil {
		return m.getByFn(ctx, workspaceID, userID)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockMemberRepo) IsWorkspaceAdmin(ctx context.Context, userID, workspaceID int64) (bool, error) {
	if m.isAdminFn != nil {
		return m.isAdminFn(ctx, userID, workspaceID)
	}
	return false, nil
}

type mockWorkspaceLookup struct {
	getByIDFn func(ctx context.Context, id int64) (*wksentity.Workspace, error)
}

func (m *mockWorkspaceLookup) GetByID(ctx context.Context, id int64) (*wksentity.Workspace, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

type mockUserLookup struct {
	getByIDFn     func(ctx context.Context, id int64) (*userentity.User, error)
	getBySubFn    func(ctx context.Context, sub string) (*userentity.User, error)
}

func (m *mockUserLookup) GetByID(ctx context.Context, id int64) (*userentity.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockUserLookup) GetByKeycloakSub(ctx context.Context, sub string) (*userentity.User, error) {
	if m.getBySubFn != nil {
		return m.getBySubFn(ctx, sub)
	}
	return nil, apierrors.ErrNotFound
}

func newTestMembership(member *mockMemberRepo, ws *mockWorkspaceLookup, users *mockUserLookup) *service.MembershipService {
	return service.NewMembershipService(member, ws, users)
}

func TestSetIsAdminGrantsAdmin(t *testing.T) {
	ctx := context.Background()
	member := &mockMemberRepo{
		getByFn: func(ctx context.Context, workspaceID, userID int64) (*wksentity.WorkspaceUser, error) {
			return &wksentity.WorkspaceUser{ID: 1, WorkspaceID: workspaceID, UserID: userID, IsAdmin: true, IsActive: true}, nil
		},
	}
	ws := &mockWorkspaceLookup{getByIDFn: func(ctx context.Context, id int64) (*wksentity.Workspace, error) {
		return &wksentity.Workspace{ID: id}, nil
	}}
	users := &mockUserLookup{
		getByIDFn: func(ctx context.Context, id int64) (*userentity.User, error) { return &userentity.User{ID: id}, nil },
		getBySubFn: func(ctx context.Context, sub string) (*userentity.User, error) { return &userentity.User{ID: 42}, nil },
	}
	svc := newTestMembership(member, ws, users)

	var gotWSID, gotUserID int64
	var gotIsAdmin bool
	var gotActor *int64
	member.upsertFn = func(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error {
		gotWSID, gotUserID, gotIsAdmin, gotActor = workspaceID, userID, isAdmin, actorID
		return nil
	}

	wu, err := svc.SetIsAdmin(ctx, 3, 7, true, "actor-sub")
	if err != nil {
		t.Fatalf("SetIsAdmin() error = %v", err)
	}
	if gotWSID != 3 || gotUserID != 7 || !gotIsAdmin {
		t.Errorf("upsert got (%d, %d, %v), want (3, 7, true)", gotWSID, gotUserID, gotIsAdmin)
	}
	if gotActor == nil || *gotActor != 42 {
		t.Errorf("actorID = %v, want 42", gotActor)
	}
	if !wu.IsAdmin {
		t.Error("response IsAdmin = false, want true")
	}
}

func TestSetIsAdminRevokesAdmin(t *testing.T) {
	ctx := context.Background()
	ws := &mockWorkspaceLookup{getByIDFn: func(ctx context.Context, id int64) (*wksentity.Workspace, error) {
		return &wksentity.Workspace{ID: id}, nil
	}}
	users := &mockUserLookup{getByIDFn: func(ctx context.Context, id int64) (*userentity.User, error) {
		return &userentity.User{ID: id}, nil
	}}

	var gotIsAdmin bool
	member := &mockMemberRepo{
		upsertFn: func(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error {
			gotIsAdmin = isAdmin
			return nil
		},
		getByFn: func(ctx context.Context, workspaceID, userID int64) (*wksentity.WorkspaceUser, error) {
			return &wksentity.WorkspaceUser{ID: 1, WorkspaceID: workspaceID, UserID: userID, IsAdmin: gotIsAdmin, IsActive: true}, nil
		},
	}
	svc := newTestMembership(member, ws, users)
	wu, err := svc.SetIsAdmin(ctx, 3, 7, false, "")
	if err != nil {
		t.Fatalf("SetIsAdmin() error = %v", err)
	}
	if gotIsAdmin {
		t.Error("upsert isAdmin = true, want false (revoke)")
	}
	if wu.IsAdmin {
		t.Error("response IsAdmin = true, want false")
	}
}

func TestSetIsAdminWorkspaceNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newTestMembership(&mockMemberRepo{}, &mockWorkspaceLookup{}, &mockUserLookup{})
	_, err := svc.SetIsAdmin(ctx, 99, 7, true, "")
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("SetIsAdmin() error = %v, want ErrNotFound", err)
	}
}

func TestSetIsAdminUserNotFound(t *testing.T) {
	ctx := context.Background()
	ws := &mockWorkspaceLookup{getByIDFn: func(ctx context.Context, id int64) (*wksentity.Workspace, error) {
		return &wksentity.Workspace{ID: id}, nil
	}}
	svc := newTestMembership(&mockMemberRepo{}, ws, &mockUserLookup{})
	_, err := svc.SetIsAdmin(ctx, 3, 99, true, "")
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("SetIsAdmin() error = %v, want ErrNotFound", err)
	}
}

func TestSetIsAdminActorUnresolvedStillUpserts(t *testing.T) {
	ctx := context.Background()
	ws := &mockWorkspaceLookup{getByIDFn: func(ctx context.Context, id int64) (*wksentity.Workspace, error) {
		return &wksentity.Workspace{ID: id}, nil
	}}
	users := &mockUserLookup{getByIDFn: func(ctx context.Context, id int64) (*userentity.User, error) {
		return &userentity.User{ID: id}, nil
	}}
	var called bool
	member := &mockMemberRepo{
		upsertFn: func(ctx context.Context, workspaceID, userID int64, isAdmin bool, actorID *int64) error {
			called = true
			if actorID != nil {
				t.Errorf("actorID = %v, want nil (actor sub not in users)", actorID)
			}
			return nil
		},
		getByFn: func(ctx context.Context, workspaceID, userID int64) (*wksentity.WorkspaceUser, error) {
			return &wksentity.WorkspaceUser{ID: 1, WorkspaceID: workspaceID, UserID: userID, IsActive: true}, nil
		},
	}
	svc := newTestMembership(member, ws, users)
	if _, err := svc.SetIsAdmin(ctx, 3, 7, true, "unknown-sub"); err != nil {
		t.Fatalf("SetIsAdmin() error = %v", err)
	}
	if !called {
		t.Error("SetIsAdmin() did not upsert")
	}
}

func TestIsWorkspaceAdminReturnsResult(t *testing.T) {
	ctx := context.Background()
	member := &mockMemberRepo{
		isAdminFn: func(ctx context.Context, userID, workspaceID int64) (bool, error) {
			return workspaceID == 10 && userID == 5, nil
		},
	}
	svc := newTestMembership(member, &mockWorkspaceLookup{}, &mockUserLookup{})

	got, err := svc.IsWorkspaceAdmin(ctx, 5, 10)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin() error = %v", err)
	}
	if !got {
		t.Error("IsWorkspaceAdmin(5, 10) = false, want true")
	}

	got, err = svc.IsWorkspaceAdmin(ctx, 5, 11)
	if err != nil {
		t.Fatalf("IsWorkspaceAdmin() error = %v", err)
	}
	if got {
		t.Error("IsWorkspaceAdmin(5, 11) = true, want false (wrong workspace)")
	}
}

func TestIsWorkspaceAdminPropagatesRepoError(t *testing.T) {
	ctx := context.Background()
	repoErr := errors.New("db down")
	member := &mockMemberRepo{
		isAdminFn: func(ctx context.Context, userID, workspaceID int64) (bool, error) {
			return false, repoErr
		},
	}
	svc := newTestMembership(member, &mockWorkspaceLookup{}, &mockUserLookup{})
	_, err := svc.IsWorkspaceAdmin(ctx, 1, 2)
	if !errors.Is(err, repoErr) {
		t.Errorf("IsWorkspaceAdmin() error = %v, want %v", err, repoErr)
	}
}
