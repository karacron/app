package domain

type UserStatusRow struct {
	ID                 string `db:"id"`
	OnboardingComplete int    `db:"onboarding_complete"`
}

type UserSettingsRow struct {
	ID                 string  `db:"id"                  json:"id"`
	InstallationID     string  `db:"installation_id"     json:"installationId"`
	FirstName          *string `db:"first_name"          json:"firstName"`
	LastName           *string `db:"last_name"           json:"lastName"`
	BirthDate          *string `db:"birth_date"          json:"birthDate"`
	DisplayName        *string `db:"display_name"        json:"displayName"`
	AvatarURL          *string `db:"avatar_url"          json:"avatarUrl"`
	Email              *string `db:"email"               json:"email"`
	Language           string  `db:"language"            json:"language"`
	Timezone           string  `db:"timezone"            json:"timezone"`
	Bio                *string `db:"bio"                 json:"bio"`
	DateFormat         *string `db:"date_format"         json:"dateFormat"`
	TimeFormat         *string `db:"time_format"         json:"timeFormat"`
	OnboardingComplete int     `db:"onboarding_complete" json:"onboardingComplete"`
	Preferences        *string `db:"preferences"         json:"preferences"`
	CreatedAt          string  `db:"created_at"          json:"createdAt"`
	UpdatedAt          string  `db:"updated_at"          json:"updatedAt"`
}

type EmailReq struct {
	Email string `validate:"required,email"`
}

type CreateUserReq struct {
	Name                *string                `json:"name"                validate:"omitempty,max=120"`
	FirstName           *string                `json:"firstName"           validate:"omitempty,max=120"`
	LastName            *string                `json:"lastName"            validate:"omitempty,max=160"`
	DisplayName         *string                `json:"displayName"         validate:"omitempty,max=160"`
	AvatarURL           *string                `json:"avatarUrl"           validate:"omitempty,max=500"`
	Email               *string                `json:"email"               validate:"omitempty,email,max=255"`
	BirthDate           *string                `json:"birthDate"           validate:"omitempty"`
	Description         *string                `json:"description"         validate:"omitempty,max=1000"`
	Bio                 *string                `json:"bio"                 validate:"omitempty,max=1000"`
	Language            *string                `json:"language"            validate:"omitempty,max=10"`
	Locale              *string                `json:"locale"              validate:"omitempty,max=10"`
	Timezone            *string                `json:"timezone"            validate:"omitempty,max=80"`
	DateFormat          *string                `json:"dateFormat"          validate:"omitempty,max=30"`
	TimeFormat          *string                `json:"timeFormat"          validate:"omitempty,max=30"`
	OnboardingComplete  *bool                  `json:"onboardingComplete"  validate:"omitempty"`
	OnboardingCompleted *bool                  `json:"onboardingCompleted" validate:"omitempty"`
	Preferences         map[string]interface{} `json:"preferences"         validate:"omitempty"`
}

type UpdateUserReq struct {
	// Aliases for compatibility with frontend payload
	Name                *string                `json:"name"                validate:"omitempty,max=120"`
	Description         *string                `json:"description"         validate:"omitempty,max=1000"`
	OnboardingCompleted *bool                  `json:"onboardingCompleted" validate:"omitempty"`
	// Primary fields
	FirstName           *string                `json:"firstName"          validate:"omitempty,max=120"`
	LastName            *string                `json:"lastName"           validate:"omitempty,max=160"`
	DisplayName         *string                `json:"displayName"        validate:"omitempty,max=160"`
	AvatarURL           *string                `json:"avatarUrl"          validate:"omitempty,max=500"`
	Email               *string                `json:"email"              validate:"omitempty,email,max=255"`
	BirthDate           *string                `json:"birthDate"          validate:"omitempty"`
	Bio                 *string                `json:"bio"                validate:"omitempty,max=1000"`
	Language            *string                `json:"language"           validate:"omitempty,max=10"`
	Timezone            *string                `json:"timezone"           validate:"omitempty,max=80"`
	DateFormat          *string                `json:"dateFormat"         validate:"omitempty,max=30"`
	TimeFormat          *string                `json:"timeFormat"         validate:"omitempty,max=30"`
	OnboardingComplete  *bool                  `json:"onboardingComplete" validate:"omitempty"`
	Preferences         map[string]interface{} `json:"preferences"        validate:"omitempty"`
}

