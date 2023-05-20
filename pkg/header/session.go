package header

import (
	"context"
	"net/http"
	"time"
)

const SessionKey = "ffx_session"
const SessionDuration = 2 * time.Hour

// SetSessionId 设置sessionId
func SetSessionId(ctx context.Context, sessionId string, domain string) {
	cookie := &http.Cookie{
		Name:     SessionKey,
		Value:    sessionId,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(SessionDuration), // 2小时后过期
		HttpOnly: true,                            // 禁止js访问此cookie
	}
	SetCookie(ctx, []*http.Cookie{cookie})
}

// ClearSessionId 清除sessionId
func ClearSessionId(ctx context.Context, sessionId string, domain string) {
	cookie := &http.Cookie{
		Name:     SessionKey,
		Value:    sessionId,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(-1 * time.Minute),
		HttpOnly: true, // 禁止js访问此cookie
	}
	SetCookie(ctx, []*http.Cookie{cookie})
}

// GetSessionId 获取sessionId
func GetSessionId(ctx context.Context) string {
	cookies := ReadCookies(ctx, SessionKey)
	if len(cookies) > 0 {
		return cookies[0].Value
	}
	return ""
}
