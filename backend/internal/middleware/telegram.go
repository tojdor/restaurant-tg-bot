package middleware

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	botSecretHeader      = "X-Bot-Secret"
	telegramUserIDHeader = "X-Telegram-User-ID"
)

type principalContextKey struct{}

// Principal is the verified identity of the Telegram user who initiated a bot action.
type Principal struct {
	UserID         int
	TelegramUserID int64
	Role           string
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok
}

// RequireRoles allows only requests made by the trusted bot service on behalf
// of a registered Telegram user having one of the permitted roles.
func RequireRoles(pool *pgxpool.Pool, botSecret string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if botSecret == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get(botSecretHeader)), []byte(botSecret)) != 1 {
				http.Error(w, "unauthorized bot request", http.StatusUnauthorized)
				return
			}

			telegramUserID, err := strconv.ParseInt(r.Header.Get(telegramUserIDHeader), 10, 64)
			if err != nil || telegramUserID <= 0 {
				http.Error(w, "X-Telegram-User-ID header is required", http.StatusUnauthorized)
				return
			}

			var principal Principal
			principal.TelegramUserID = telegramUserID
			err = pool.QueryRow(r.Context(), `SELECT id, role FROM users WHERE telegram_user_id = $1`, telegramUserID).
				Scan(&principal.UserID, &principal.Role)
			if err != nil {
				if err == pgx.ErrNoRows {
					http.Error(w, "telegram user is not registered", http.StatusForbidden)
					return
				}
				http.Error(w, "failed to authorize user", http.StatusInternalServerError)
				return
			}

			for _, role := range roles {
				if principal.Role == role {
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalContextKey{}, principal)))
					return
				}
			}

			http.Error(w, "insufficient permissions", http.StatusForbidden)
		})
	}
}
