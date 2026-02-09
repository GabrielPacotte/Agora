package services

import (
	"context"
	"errors"
	"time"

	"github.com/GabrielPacotte/Agora/internal/auth"
	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, login, displayName, password string) (*auth.TokenPair, uuid.UUID, error)
	Login(ctx context.Context, login, password string) (*auth.TokenPair, uuid.UUID, error)
	Refresh(ctx context.Context, refreshTokenRaw string) (*auth.TokenPair, uuid.UUID, error)
	Logout(ctx context.Context, refreshTokenRaw string) error
}

type authService struct {
	users      repositories.UserRepository
	sessions   repositories.SessionRepository
	jwt        auth.JWTIssuer
	opaque     auth.OpaqueTokenGenerator
	now        func() time.Time
	accessTTL  time.Duration
	refreshTTL time.Duration

	refreshNBytes int
	bcryptCost    int
}

func NewAuthService(
	users repositories.UserRepository,
	sessions repositories.SessionRepository,
	jwt auth.JWTIssuer,
	opaque auth.OpaqueTokenGenerator,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) AuthService {
	return &authService{
		users:         users,
		sessions:      sessions,
		jwt:           jwt,
		opaque:        opaque,
		now:           func() time.Time { return time.Now().UTC() },
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		refreshNBytes: 32,
		bcryptCost:    bcrypt.DefaultCost,
	}
}

func (s *authService) Register(ctx context.Context, login, displayName, password string) (*auth.TokenPair, uuid.UUID, error) {
	now := s.now()

	if password == "" {
		return nil, uuid.Nil, domain.ErrInvalidPassword
	}

	user, err := domain.NewUser(now, login, displayName)
	if err != nil {
		return nil, uuid.Nil, err
	}

	hash, err := bcryptHash(password, s.bcryptCost)
	if err != nil {
		return nil, uuid.Nil, err
	}

	if err := s.users.Create(ctx, user, hash); err != nil {
		return nil, uuid.Nil, err
	}

	pair, err := s.issueSessionAndTokens(ctx, user.ID(), now)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return pair, user.ID(), nil
}

func (s *authService) Login(ctx context.Context, login, password string) (*auth.TokenPair, uuid.UUID, error) {
	if password == "" {
		return nil, uuid.Nil, domain.ErrInvalidPassword
	}

	user, hash, err := s.users.GetByLoginWithPasswordHash(ctx, login)
	if err != nil {
		return nil, uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, uuid.Nil, domain.ErrInvalidCredentials
	}

	now := s.now()

	pair, err := s.issueSessionAndTokens(ctx, user.ID(), now)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return pair, user.ID(), nil
}

func (s *authService) Refresh(ctx context.Context, refreshTokenRaw string) (*auth.TokenPair, uuid.UUID, error) {
	now := s.now()

	if refreshTokenRaw == "" {
		return nil, uuid.Nil, domain.ErrInvalidToken
	}

	session, err := s.sessions.GetByRefreshToken(ctx, refreshTokenRaw)
	if err != nil {
		return nil, uuid.Nil, err
	}

	if session.IsRevoked() || now.After(session.ExpiresAt()) {
		return nil, uuid.Nil, domain.ErrInvalidToken
	}

	if err := s.sessions.TouchLastUsed(ctx, session.ID(), now); err != nil {
		return nil, uuid.Nil, err
	}

	if err := s.sessions.RevokeByRefreshToken(ctx, refreshTokenRaw, now); err != nil {
		return nil, uuid.Nil, err
	}

	pair, err := s.issueSessionAndTokens(ctx, session.UserID(), now)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return pair, session.UserID(), nil
}

func (s *authService) Logout(ctx context.Context, refreshTokenRaw string) error {
	now := s.now()

	if refreshTokenRaw == "" {
		return domain.ErrInvalidToken
	}

	err := s.sessions.RevokeByRefreshToken(ctx, refreshTokenRaw, now)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (s *authService) issueSessionAndTokens(ctx context.Context, userID uuid.UUID, now time.Time) (*auth.TokenPair, error) {
	access, err := s.jwt.GenerateAccessToken(userID, now, s.accessTTL)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := s.opaque.GenerateOpaqueToken(s.refreshNBytes)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := now.Add(s.refreshTTL)
	refresh, err := auth.NewRefreshToken(rawRefresh, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	pair, err := auth.NewTokenPair(access, refresh, s.accessTTL)
	if err != nil {
		return nil, err
	}

	session, err := auth.NewSession(uuid.New(), userID, now, refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	if err := s.sessions.Create(ctx, session, refresh.Raw()); err != nil {
		return nil, err
	}

	return &pair, nil
}

func bcryptHash(password string, cost int) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
