package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/brix101/linkr/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	nanoidAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nanoidSize     = 8
)

func (a *api) getLink(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	code := r.PathValue("code")

	if code == "" {
		http.Error(w, "code is required", http.StatusBadRequest)
		return
	}

	link, err := a.DB.GetLinkByCode(ctx, code)
	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(link.Url, "http://") && !strings.HasPrefix(link.Url, "https://") {
		http.Error(w, "invalid redirect URL", http.StatusBadRequest)
		return
	}

	if link.ID == 0 {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		ip := getClientIP(r)
		userAgent := r.UserAgent()
		referrer := r.Referer()

		if err = a.DB.CreateClick(ctx, db.CreateClickParams{
			LinkID: link.ID,
			IpAddress: pgtype.Text{
				String: ip,
				Valid:  ip != "",
			},
			UserAgent: pgtype.Text{
				String: userAgent,
				Valid:  userAgent != "",
			},
			Referrer: pgtype.Text{
				String: referrer,
				Valid:  referrer != "",
			},
		}); err != nil {
			slog.Error("Creating clink error", slog.Any("error", err.Error()))
		}
	}()

	http.Redirect(w, r, link.Url, http.StatusPermanentRedirect)
}

type CreateLinkRequest struct {
	URL string `json:"url"`
}

func (a *api) createLink(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "invalid url format", http.StatusBadRequest)
		return
	}

	link, err := a.DB.GetLinkByURL(ctx, req.URL)
	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if link.Code != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(link)
		return
	}

	code, err := gonanoid.Generate(nanoidAlphabet, nanoidSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nLink, err := a.DB.CreateLink(ctx, db.CreateLinkParams{
		Code: code,
		Url:  req.URL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nLink)
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
