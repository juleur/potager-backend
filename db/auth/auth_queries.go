package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"npp_backend/entity/private"
	"npp_backend/l10n/translate"
	"npp_backend/pkg/password"
	"npp_backend/pkg/token"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (ap *AuthPsql) CreateNewUser(ctx context.Context, username string, email string, pwd string) (*private.User, *private.Error) {
	hashedPassword, err := password.HashPassword(pwd)
	if err != nil {
		return nil, err
	}
	refreshToken := token.HexKeyGenerator(32)
	user := private.User{}
	user.VerifiedAt = &sql.NullString{}
	if err := ap.db.QueryRow(ctx, `
		WITH created_user AS (
			INSERT INTO users(username, email, password, refresh_token) VALUES($1,$2,$3,$4)
			RETURNING id, username, refresh_token, verified_at
		), create_user_permission AS (
			INSERT INTO users_permission(user_id) SELECT id FROM created_user
		)
		SELECT id, username, refresh_token, verified_at FROM created_user
	`, username, email, hashedPassword, refreshToken).Scan(&user.ID, &user.Username, &user.RefreshToken, user.VerifiedAt); err != nil {
		if strings.Contains(err.Error(), "users_username_unique_idx") {
			return nil, &private.Error{
				Location:   "db.auth.CreateNewUser",
				Line:       37,
				Err:        fmt.Errorf("%w : email %s déjà pris", err, username),
				TranslKey:  translate.KeyUsernameUniqueness,
				ErrorCode:  1,
				StatusCode: http.StatusConflict,
			}
		} else if strings.Contains(err.Error(), "users_email_unique_idx") {
			return nil, &private.Error{
				Location:   "db.auth.CreateNewUser",
				Line:       46,
				Err:        fmt.Errorf("%w : username %s déjà pris", err, email),
				TranslKey:  translate.KeyEmailUniqueness,
				ErrorCode:  1,
				StatusCode: http.StatusConflict,
			}
		}
		return nil, &private.Error{
			Location:   "db.auth.CreateNewUser",
			Line:       55,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return &user, nil
}

func (ap *AuthPsql) GetUser(ctx context.Context, identifiant string) (*private.User, *private.Error) {
	user, err := getUserByIdentifiant(ap.db, ctx, identifiant)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ap *AuthPsql) FailedAttempt(ctx context.Context, userID int) *private.Error {
	if _, err := ap.db.Exec(ctx, `
		UPDATE users_permission
		SET failed_login_attempts = CASE
			WHEN failed_login_attempts = 0 THEN 1
			WHEN failed_login_attempts = 1 THEN 2
			WHEN failed_login_attempts = 2 THEN 3
			ELSE 3
		END,
		locked_until = CASE
			WHEN locked_until IS NULL AND failed_login_attempts = 2 THEN CURRENT_TIMESTAMP + '6 hours'::INTERVAL
			WHEN locked_until IS NOT NULL THEN locked_until
			ELSE NULL
		END
		WHERE user_id = $1;
	`, userID); err != nil {
		return &private.Error{
			Location:   "db.auth.FailedAttempt",
			Line:       91,
			Err:        fmt.Errorf("%w: user n°%d", err, userID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (ap *AuthPsql) IsRefreshTokenExists(ctx context.Context, refreshToken string) (*private.User, *private.Error) {
	user := private.User{}
	user.VerifiedAt = &sql.NullString{}
	user.UserPermission.LockedUntil = &sql.NullString{}
	if err := ap.db.QueryRow(ctx, `
		SELECT u.id, u.username, u.verified_at, up.locked_until
		FROM users u JOIN users_permission up ON up.user_id = u.id
		WHERE refresh_token=$1
		LIMIT 1;
	`, refreshToken).Scan(&user.ID, &user.Username, user.VerifiedAt, user.UserPermission.LockedUntil); err != nil {
		if err == pgx.ErrNoRows {
			return nil, &private.Error{
				Location:   "db.auth.IsRefreshTokenExists",
				Line:       114,
				Err:        errors.New("refreshToken introuvable"),
				TranslKey:  translate.KeyInvalidToken,
				ErrorCode:  41,
				StatusCode: http.StatusUnauthorized,
			}
		}
		return nil, &private.Error{
			Location:   "db.auth.IsRefreshTokenExists",
			Line:       123,
			Err:        err,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  41,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return &user, nil
}

func (ap *AuthPsql) UpdateLastLoggedInAt(ctx context.Context, userID int) *private.Error {
	if _, err := ap.db.Exec(ctx, `
		UPDATE users u
		SET last_logged_in_at = $1
		FROM users_permission up
		WHERE up.user_id = u.id AND u.id = $2 AND up.status = 2 AND up.failed_login_attempts < 2 AND up.locked_until IS NULL;
	`, time.Now(), userID); err != nil {
		return &private.Error{
			Location:   "db.auth.UpdateLastLoggedInAt",
			Line:       142,
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (ap *AuthPsql) UpdateRefreshToken(ctx context.Context, newRefreshToken string, userID int) *private.Error {
	if _, err := ap.db.Exec(ctx, `
		UPDATE users
		SET refresh_token = $1, last_logged_in_at = CURRENT_TIMESTAMP
		WHERE id = $2
		AND
		(SELECT EXISTS(SELECT * FROM users_permission up WHERE up.user_id = $2 AND up.status = 2 AND up.failed_login_attempts <= 2 AND up.locked_until IS NULL));
	`, newRefreshToken, userID); err != nil {
		fmt.Println(err)
		return &private.Error{
			Location:   "db.auth.UpdateRefreshToken",
			Line:       160,
			Err:        fmt.Errorf("%w: user n°%d", err, userID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func getUserByIdentifiant(db *pgxpool.Pool, ctx context.Context, identifiant string) (*private.User, *private.Error) {
	user := private.User{}
	user.VerifiedAt = &sql.NullString{}
	user.UserPermission.LockedUntil = &sql.NullString{}
	err := db.QueryRow(ctx, `
		SELECT
			u.id, u.username, u.password, u.refresh_token, u.verified_at, up.locked_until,
			CASE
				WHEN up.status = 0 THEN 'SUSPENDU'
				WHEN up.failed_login_attempts = 3 AND up.locked_until > CURRENT_TIMESTAMP THEN 'LOCKED'
				ELSE 'OK'
			END AS user_login
		FROM users u
		JOIN users_permission up ON up.user_id = u.id
		WHERE u.username = $1 OR u.email = $1
		LIMIT 1;
	`, identifiant).Scan(&user.ID, &user.Username, &user.Password, &user.RefreshToken, user.VerifiedAt, user.UserPermission.LockedUntil, &user.State)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &private.Error{
				Location:   "db.auth.getUserByIdentifiant",
				Line:       191,
				Err:        fmt.Errorf("User n°%d introuvable", user.ID),
				TranslKey:  translate.KeyUserNotFound,
				ErrorCode:  1,
				StatusCode: http.StatusNotFound,
			}
		}
		return nil, &private.Error{
			Location:   "db.auth.getUserByIdentifiant",
			Line:       200,
			Err:        fmt.Errorf("%w: user n°%d", err, user.ID),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}

	if user.State == "OK" {
		if err := resetAttempt(db, ctx, user.ID); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func resetAttempt(db *pgxpool.Pool, ctx context.Context, userId int) *private.Error {
	if _, err := db.Exec(ctx, `
		UPDATE users_permission
		SET failed_login_attempts = CASE
			WHEN locked_until IS NOT NULL AND failed_login_attempts = 3 THEN 0
			ELSE failed_login_attempts
		END,
		locked_until = CASE
			WHEN locked_until IS NOT NULL AND failed_login_attempts = 3 THEN CAST(NULL AS TIMESTAMP)
			ELSE locked_until
		END
		WHERE user_id = $1;
	`, userId); err != nil {
		return &private.Error{
			Location:   "db.auth.resetAttempt",
			Line:       231,
			Err:        fmt.Errorf("%w: user n°%d", err, userId),
			TranslKey:  translate.KeyInternalServerError,
			ErrorCode:  1,
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}
