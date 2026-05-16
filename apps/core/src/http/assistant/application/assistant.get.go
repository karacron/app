package application

import (
	dbmod "github.com/authuser-dev/karacron/apps/core/src/module/database"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type assistantPersonality struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	DisplayName string  `json:"displayName"`
	Description *string `json:"description"`
	Domain      *string `json:"domain"`
	IsDefault   bool    `json:"isDefault"`
}

type assistantTone struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	DisplayName string  `json:"displayName"`
	Description *string `json:"description"`
	IsDefault   bool    `json:"isDefault"`
}

type assistantPersonalityRow struct {
	Reference   string  `db:"reference"`
	Domain      *string `db:"domain"`
	IsDefault   int     `db:"is_default"`
	DisplayName string  `db:"display_name"`
	Description *string `db:"description"`
}

type assistantToneRow struct {
	Reference   string  `db:"reference"`
	IsDefault   int     `db:"is_default"`
	DisplayName string  `db:"display_name"`
	Description *string `db:"description"`
}

type AssistantGetUseCase struct {
	DB  *sqlx.DB
	Log *zap.Logger
}

func NewAssistantGetUseCase(db *sqlx.DB, logger *zap.Logger) *AssistantGetUseCase {
	return &AssistantGetUseCase{DB: db, Log: logger}
}

func (u *AssistantGetUseCase) GetPersonalities(c fiber.Ctx) error {
	locale := u.resolveLocale(c)

	rows := []assistantPersonalityRow{}
	err := u.DB.Select(
		&rows,
		`SELECT
			p.reference,
			p.domain,
			p.is_default,
			COALESCE(tl.display_name, te.display_name, p.reference) AS display_name,
			COALESCE(tl.description, te.description) AS description
		FROM assistant_personalities p
		LEFT JOIN assistant_personality_translations tl
			ON tl.personality_id = p.id AND tl.locale = ?
		LEFT JOIN assistant_personality_translations te
			ON te.personality_id = p.id AND te.locale = 'en'
		WHERE p.is_active = 1
		ORDER BY p.is_default DESC, p.id ASC`,
		locale,
	)
	if err != nil {
		u.Log.Error("GetPersonalities query", zap.Error(err), zap.String("locale", locale))
		return fiber.ErrInternalServerError
	}

	out := make([]assistantPersonality, 0, len(rows))
	for _, row := range rows {
		out = append(out, assistantPersonality{
			ID:          row.Reference,
			Slug:        row.Reference,
			DisplayName: row.DisplayName,
			Description: row.Description,
			Domain:      row.Domain,
			IsDefault:   row.IsDefault == 1,
		})
	}

	return c.JSON(out)
}

func (u *AssistantGetUseCase) GetTones(c fiber.Ctx) error {
	locale := u.resolveLocale(c)

	rows := []assistantToneRow{}
	err := u.DB.Select(
		&rows,
		`SELECT
			t.reference,
			t.is_default,
			COALESCE(tl.display_name, te.display_name, t.reference) AS display_name,
			COALESCE(tl.description, te.description) AS description
		FROM assistant_tones t
		LEFT JOIN assistant_tone_translations tl
			ON tl.tone_id = t.id AND tl.locale = ?
		LEFT JOIN assistant_tone_translations te
			ON te.tone_id = t.id AND te.locale = 'en'
		WHERE t.is_active = 1
		ORDER BY t.is_default DESC, t.id ASC`,
		locale,
	)
	if err != nil {
		u.Log.Error("GetTones query", zap.Error(err), zap.String("locale", locale))
		return fiber.ErrInternalServerError
	}

	out := make([]assistantTone, 0, len(rows))
	for _, row := range rows {
		out = append(out, assistantTone{
			ID:          row.Reference,
			Slug:        row.Reference,
			DisplayName: row.DisplayName,
			Description: row.Description,
			IsDefault:   row.IsDefault == 1,
		})
	}

	return c.JSON(out)
}

func (u *AssistantGetUseCase) resolveLocale(c fiber.Ctx) string {
	if locale := normalizeLocale(c.Query("locale")); locale != "" {
		return locale
	}
	if locale := normalizeLocale(c.Query("language")); locale != "" {
		return locale
	}

	if locale := normalizeLocale(c.Get("X-Locale")); locale != "" {
		return locale
	}

	acceptLanguage := c.Get("Accept-Language")
	if acceptLanguage != "" {
		for _, part := range strings.Split(acceptLanguage, ",") {
			tag := strings.TrimSpace(strings.Split(part, ";")[0])
			if locale := normalizeLocale(tag); locale != "" {
				return locale
			}
		}
	}

	installID, err := dbmod.EnsureInstallation(u.DB)
	if err == nil {
		var lang string
		if err := u.DB.Get(&lang, `SELECT language FROM user_settings WHERE installation_id = ? LIMIT 1`, installID); err == nil {
			if locale := normalizeLocale(lang); locale != "" {
				return locale
			}
		}
	}

	return "en"
}

func normalizeLocale(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}

	if idx := strings.Index(raw, "-"); idx > 0 {
		raw = raw[:idx]
	}
	if idx := strings.Index(raw, "_"); idx > 0 {
		raw = raw[:idx]
	}

	if len(raw) < 2 || len(raw) > 5 {
		return ""
	}

	return raw
}

