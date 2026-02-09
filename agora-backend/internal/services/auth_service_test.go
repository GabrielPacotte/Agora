package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GabrielPacotte/Agora/internal/auth"
	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type mockSessionRepo struct {
	CreateFn               func(ctx context.Context, s *auth.Session, refreshTokenRaw string) error
	GetByRefreshTokenFn    func(ctx context.Context, refreshTokenRaw string) (*auth.Session, error)
	TouchLastUsedFn        func(ctx context.Context, sessionID uuid.UUID, t time.Time) error
	RevokeByRefreshTokenFn func(ctx context.Context, refreshTokenRaw string, t time.Time) error
}

func (m *mockSessionRepo) Create(ctx context.Context, s *auth.Session, refreshTokenRaw string) error {
	if m.CreateFn == nil {
		panic("CreateFn not set")
	}
	return m.CreateFn(ctx, s, refreshTokenRaw)
}

func (m *mockSessionRepo) GetByRefreshToken(ctx context.Context, refreshTokenRaw string) (*auth.Session, error) {
	if m.GetByRefreshTokenFn == nil {
		panic("GetByRefreshTokenFn not set")
	}
	return m.GetByRefreshTokenFn(ctx, refreshTokenRaw)
}

func (m *mockSessionRepo) TouchLastUsed(ctx context.Context, sessionID uuid.UUID, t time.Time) error {
	if m.TouchLastUsedFn == nil {
		panic("TouchLastUsedFn not set")
	}
	return m.TouchLastUsedFn(ctx, sessionID, t)
}

func (m *mockSessionRepo) RevokeByRefreshToken(ctx context.Context, refreshTokenRaw string, t time.Time) error {
	if m.RevokeByRefreshTokenFn == nil {
		panic("RevokeByRefreshTokenFn not set")
	}
	return m.RevokeByRefreshTokenFn(ctx, refreshTokenRaw, t)
}

type mockJWTIssuer struct {
	GenerateFn func(userID uuid.UUID, now time.Time, ttl time.Duration) (auth.AccessToken, error)
}

func (m *mockJWTIssuer) GenerateAccessToken(userID uuid.UUID, now time.Time, ttl time.Duration) (auth.AccessToken, error) {
	if m.GenerateFn == nil {
		panic("GenerateFn not set")
	}
	return m.GenerateFn(userID, now, ttl)
}

type mockOpaque struct {
	GenerateFn func(nBytes int) (string, error)
}

func (m *mockOpaque) GenerateOpaqueToken(nBytes int) (string, error) {
	if m.GenerateFn == nil {
		panic("GenerateFn not set")
	}
	return m.GenerateFn(nBytes)
}

func mustUser(t *testing.T, at time.Time, login, display string) *domain.User {
	t.Helper()
	u, err := domain.NewUser(at, login, display)
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	return u
}

func TestAuthService_Register_EmptyPassword(t *testing.T) {
	ctx := context.Background()

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(context.Context, string) (*domain.User, string, error) {
			t.Fatalf("should not be called")
			return nil, "", nil
		},
	}

	sessions := &mockSessionRepo{
		CreateFn: func(context.Context, *auth.Session, string) error {
			t.Fatalf("Create should not be called")
			return nil
		},
		GetByRefreshTokenFn:    func(context.Context, string) (*auth.Session, error) { panic("not used") },
		TouchLastUsedFn:        func(context.Context, uuid.UUID, time.Time) error { panic("not used") },
		RevokeByRefreshTokenFn: func(context.Context, string, time.Time) error { panic("not used") },
	}

	jwt := &mockJWTIssuer{GenerateFn: func(uuid.UUID, time.Time, time.Duration) (auth.AccessToken, error) {
		t.Fatalf("jwt should not be called")
		return "", nil
	}}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { t.Fatalf("opaque should not be called"); return "", nil }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	_, _, err := svc.Register(ctx, "johndoe", "John Doe", "")
	if !errors.Is(err, domain.ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	ctx := context.Background()

	u := mustUser(t, time.Now().Add(-time.Hour).UTC(), "johndoe", "John Doe")

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(ctx context.Context, login string) (*domain.User, string, error) {
			if login != "johndoe" {
				return nil, "", domain.ErrNotFound
			}
			hash, _ := bcryptHashForTest(t, "pw")
			return u, hash, nil
		},
	}

	sessions := &mockSessionRepo{
		CreateFn: func(ctx context.Context, s *auth.Session, refreshTokenRaw string) error {
			return nil
		},
		GetByRefreshTokenFn:    func(context.Context, string) (*auth.Session, error) { panic("not used") },
		TouchLastUsedFn:        func(context.Context, uuid.UUID, time.Time) error { panic("not used") },
		RevokeByRefreshTokenFn: func(context.Context, string, time.Time) error { panic("not used") },
	}

	jwt := &mockJWTIssuer{
		GenerateFn: func(userID uuid.UUID, now time.Time, ttl time.Duration) (auth.AccessToken, error) {
			return auth.AccessToken("jwt"), nil
		},
	}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { return "refresh", nil }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	pair, uid, err := svc.Login(ctx, "johndoe", "pw")
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if uid != u.ID() {
		t.Fatalf("expected uid=%v, got %v", u.ID(), uid)
	}
	if pair == nil || pair.Access == "" || pair.Refresh.Raw() == "" {
		t.Fatalf("expected non-empty token pair")
	}
}

func TestAuthService_Login_BadPassword(t *testing.T) {
	ctx := context.Background()
	u := mustUser(t, time.Now().Add(-time.Hour).UTC(), "johndoe", "John Doe")

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(ctx context.Context, login string) (*domain.User, string, error) {
			hash, _ := bcryptHashForTest(t, "pw")
			return u, hash, nil
		},
	}

	sessions := &mockSessionRepo{
		CreateFn: func(context.Context, *auth.Session, string) error {
			t.Fatalf("sessions.Create should not be called")
			return nil
		},
		GetByRefreshTokenFn:    func(context.Context, string) (*auth.Session, error) { panic("not used") },
		TouchLastUsedFn:        func(context.Context, uuid.UUID, time.Time) error { panic("not used") },
		RevokeByRefreshTokenFn: func(context.Context, string, time.Time) error { panic("not used") },
	}

	jwt := &mockJWTIssuer{GenerateFn: func(uuid.UUID, time.Time, time.Duration) (auth.AccessToken, error) {
		t.Fatalf("jwt should not be called")
		return "", nil
	}}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { t.Fatalf("opaque should not be called"); return "", nil }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	_, _, err := svc.Login(ctx, "johndoe", "wrong")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Refresh_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	sessionID := uuid.New()

	now := time.Now().UTC()
	expiresAt := now.Add(30 * time.Minute)

	sess, err := auth.NewSession(sessionID, userID, now.Add(-time.Hour), expiresAt)
	if err != nil {
		t.Fatalf("NewSession: %v", err)
	}

	var touched bool
	var revoked bool
	var created bool

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(context.Context, string) (*domain.User, string, error) { panic("not used") },
	}

	sessions := &mockSessionRepo{
		GetByRefreshTokenFn: func(ctx context.Context, raw string) (*auth.Session, error) {
			if raw != "refresh_raw" {
				return nil, domain.ErrNotFound
			}
			return sess, nil
		},
		TouchLastUsedFn: func(ctx context.Context, sid uuid.UUID, t time.Time) error {
			touched = true
			if sid != sessionID {
				return auth.ErrInvalidSession
			}
			return nil
		},
		RevokeByRefreshTokenFn: func(ctx context.Context, raw string, t time.Time) error {
			revoked = true
			return nil
		},
		CreateFn: func(ctx context.Context, s *auth.Session, raw string) error {
			created = true
			if s.UserID() != userID {
				t.Fatalf("expected userID")
			}
			if raw == "" {
				t.Fatalf("expected refresh raw")
			}
			return nil
		},
	}

	jwt := &mockJWTIssuer{GenerateFn: func(uuid.UUID, time.Time, time.Duration) (auth.AccessToken, error) {
		return auth.AccessToken("jwt2"), nil
	}}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { return "refresh_new", nil }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	pair, gotUID, err := svc.Refresh(ctx, "refresh_raw")
	if err != nil {
		t.Fatalf("Refresh error: %v", err)
	}
	if gotUID != userID {
		t.Fatalf("uid mismatch")
	}
	if pair == nil || pair.Access == "" || pair.Refresh.Raw() == "" {
		t.Fatalf("expected token pair")
	}
	if !touched || !revoked || !created {
		t.Fatalf("expected touched=%v revoked=%v created=%v", touched, revoked, created)
	}
}

func TestAuthService_Refresh_EmptyToken(t *testing.T) {
	ctx := context.Background()

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(context.Context, string) (*domain.User, string, error) { panic("not used") },
	}
	sessions := &mockSessionRepo{
		CreateFn: func(context.Context, *auth.Session, string) error { panic("not used") },
		GetByRefreshTokenFn: func(context.Context, string) (*auth.Session, error) {
			t.Fatalf("should not be called")
			return nil, nil
		},
		TouchLastUsedFn:        func(context.Context, uuid.UUID, time.Time) error { panic("not used") },
		RevokeByRefreshTokenFn: func(context.Context, string, time.Time) error { panic("not used") },
	}
	jwt := &mockJWTIssuer{GenerateFn: func(uuid.UUID, time.Time, time.Duration) (auth.AccessToken, error) { panic("not used") }}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { panic("not used") }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	_, _, err := svc.Refresh(ctx, "")
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestAuthService_Logout_NotFoundIsOK(t *testing.T) {
	ctx := context.Background()

	users := &mockUserRepository{
		GetByLoginWithPasswordHashFn: func(context.Context, string) (*domain.User, string, error) { panic("not used") },
	}

	sessions := &mockSessionRepo{
		RevokeByRefreshTokenFn: func(ctx context.Context, raw string, t time.Time) error {
			return domain.ErrNotFound
		},
		CreateFn:            func(context.Context, *auth.Session, string) error { panic("not used") },
		GetByRefreshTokenFn: func(context.Context, string) (*auth.Session, error) { panic("not used") },
		TouchLastUsedFn:     func(context.Context, uuid.UUID, time.Time) error { panic("not used") },
	}

	jwt := &mockJWTIssuer{GenerateFn: func(uuid.UUID, time.Time, time.Duration) (auth.AccessToken, error) { panic("not used") }}
	opaque := &mockOpaque{GenerateFn: func(int) (string, error) { panic("not used") }}

	svc := services.NewAuthService(users, sessions, jwt, opaque, 15*time.Minute, 7*24*time.Hour)

	if err := svc.Logout(ctx, "whatever"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func bcryptHashForTest(t *testing.T, password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
