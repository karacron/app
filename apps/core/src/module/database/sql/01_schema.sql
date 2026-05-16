CREATE TABLE IF NOT EXISTS users (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       VARCHAR(255) NOT NULL,
	surname    VARCHAR(255) NOT NULL,
	email      VARCHAR(255) NOT NULL UNIQUE,
	birthdate  DATE         NOT NULL,
	created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assistants (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	reference  VARCHAR(255) NOT NULL UNIQUE,
	name       VARCHAR(255) NOT NULL,
	is_local   BOOLEAN      NOT NULL DEFAULT FALSE,
	is_active  BOOLEAN      NOT NULL DEFAULT FALSE,
	is_primary BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assistant_personalities (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	reference       VARCHAR(100) NOT NULL UNIQUE,
	domain          VARCHAR(100),
	skill_reference VARCHAR(255),
	is_default      BOOLEAN      NOT NULL DEFAULT FALSE,
	is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
	created_at      DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at      DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assistant_personality_translations (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	personality_id  INT          NOT NULL REFERENCES assistant_personalities (id) ON DELETE CASCADE,
	locale          VARCHAR(10)  NOT NULL,
	display_name    VARCHAR(255) NOT NULL,
	description     TEXT,
	created_at      DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at      DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (personality_id, locale)
);

CREATE TABLE IF NOT EXISTS assistant_tones (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	reference   VARCHAR(100) NOT NULL UNIQUE,
	is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
	is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assistant_tone_translations (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	tone_id       INT          NOT NULL REFERENCES assistant_tones (id) ON DELETE CASCADE,
	locale        VARCHAR(10)  NOT NULL,
	display_name  VARCHAR(255) NOT NULL,
	description   TEXT,
	created_at    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (tone_id, locale)
);

CREATE TABLE IF NOT EXISTS permissions (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        VARCHAR(255) NOT NULL,
	description TEXT,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skills (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	reference   VARCHAR(255) NOT NULL UNIQUE,
	name        VARCHAR(255) NOT NULL,
	description TEXT,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS setting (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id      INT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	assistant_id INT         NOT NULL REFERENCES assistants (id) ON DELETE SET NULL,
	created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS integrations (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	reference   VARCHAR(255) NOT NULL UNIQUE,
	name        VARCHAR(255) NOT NULL,
	description TEXT,
	config      TEXT,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audits (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id     INT          REFERENCES users (id) ON DELETE SET NULL,
	event       VARCHAR(255) NOT NULL,
	description TEXT,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS installation (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	reference    VARCHAR(255) NOT NULL UNIQUE,
	version      VARCHAR(50)  NOT NULL,
	platform     VARCHAR(50)  NOT NULL,
	is_onboarded BOOLEAN      NOT NULL DEFAULT FALSE,
	first_run_at DATETIME,
	created_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS model_providers (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       VARCHAR(255) NOT NULL,
	reference  VARCHAR(255) NOT NULL UNIQUE,
	type       VARCHAR(50)  NOT NULL,
	base_url   VARCHAR(512),
	key_ref    VARCHAR(512),
	is_active  BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS models (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	provider_id    INT          NOT NULL REFERENCES model_providers (id) ON DELETE CASCADE,
	name           VARCHAR(255) NOT NULL,
	reference      VARCHAR(255) NOT NULL,
	type           VARCHAR(50)  NOT NULL,
	context_window INT,
	routing_priority INT         NOT NULL DEFAULT 100,
	complexity_tier VARCHAR(50)  NOT NULL DEFAULT 'standard',
	capabilities   TEXT,
	is_active      BOOLEAN      NOT NULL DEFAULT FALSE,
	is_default     BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assistant_configs (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	assistant_id   INT          NOT NULL REFERENCES assistants (id) ON DELETE CASCADE,
	model_id       INT          NOT NULL REFERENCES models (id),
	system_prompt  TEXT,
	temperature    REAL NOT NULL DEFAULT 0.7,
	max_tokens     INT,
	memory_enabled BOOLEAN      NOT NULL DEFAULT FALSE,
	tools_enabled  BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS channels (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        VARCHAR(255) NOT NULL,
	type        VARCHAR(100) NOT NULL UNIQUE,
	description TEXT,
	is_local    BOOLEAN      NOT NULL DEFAULT FALSE,
	is_active   BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at  DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS channel_connections (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_id   INT          NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
	user_id      INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	name         VARCHAR(255) NOT NULL,
	reference    VARCHAR(255),
	status       VARCHAR(50)  NOT NULL DEFAULT 'disconnected',
	dm_policy    VARCHAR(50)  NOT NULL DEFAULT 'pairing',
	group_policy VARCHAR(50)  NOT NULL DEFAULT 'disabled',
	key_ref      VARCHAR(512),
	config       TEXT,
	created_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contacts (
	id                    INTEGER PRIMARY KEY AUTOINCREMENT,
	channel_connection_id INT          NOT NULL REFERENCES channel_connections (id) ON DELETE CASCADE,
	external_id           VARCHAR(255) NOT NULL,
	display_name          VARCHAR(255),
	is_bot                BOOLEAN      NOT NULL DEFAULT FALSE,
	is_blocked            BOOLEAN      NOT NULL DEFAULT FALSE,
	metadata              TEXT,
	created_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (channel_connection_id, external_id)
);

CREATE TABLE IF NOT EXISTS sessions (
	id                    INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id               INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	assistant_id          INT          NOT NULL REFERENCES assistants (id),
	channel_connection_id INT          REFERENCES channel_connections (id) ON DELETE SET NULL,
	title                 VARCHAR(255),
	status                VARCHAR(50)  NOT NULL DEFAULT 'active',
	summary               TEXT,
	started_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	finished_at           DATETIME,
	created_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at            DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	session_id   INT          NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
	model_id     INT          REFERENCES models (id) ON DELETE SET NULL,
	role         VARCHAR(50)  NOT NULL,
	content      TEXT         NOT NULL,
	content_type VARCHAR(50)  NOT NULL DEFAULT 'text',
	tokens_used  INT,
	latency_ms   INT,
	is_pinned    BOOLEAN      NOT NULL DEFAULT FALSE,
	created_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at   DATETIME
);

CREATE TABLE IF NOT EXISTS memories (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id           INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	assistant_id      INT          NOT NULL REFERENCES assistants (id),
	source_session_id INT          REFERENCES sessions (id) ON DELETE SET NULL,
	content           TEXT         NOT NULL,
	category          VARCHAR(100) NOT NULL DEFAULT 'fact',
	confidence        DECIMAL(3,2) NOT NULL DEFAULT 1.0,
	is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
	created_at        DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at        DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS integration_connections (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	integration_id INT          NOT NULL REFERENCES integrations (id) ON DELETE CASCADE,
	user_id        INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	name           VARCHAR(255) NOT NULL,
	reference      VARCHAR(255),
	status         VARCHAR(50)  NOT NULL DEFAULT 'disconnected',
	key_ref        VARCHAR(512),
	scopes         TEXT,
	config         TEXT,
	last_sync_at   DATETIME,
	created_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skill_assignments (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	assistant_id INT         NOT NULL REFERENCES assistants (id) ON DELETE CASCADE,
	skill_id     INT         NOT NULL REFERENCES skills (id) ON DELETE CASCADE,
	is_enabled   BOOLEAN     NOT NULL DEFAULT FALSE,
	config       TEXT,
	created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (assistant_id, skill_id)
);

CREATE TABLE IF NOT EXISTS permission_grants (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	permission_id INT         NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
	entity_type   VARCHAR(50) NOT NULL,
	entity_id     INT         NOT NULL,
	granted_by    INT         NOT NULL REFERENCES users (id),
	granted_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at    DATETIME,
	created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS devices (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id       INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	name          VARCHAR(255) NOT NULL,
	type          VARCHAR(50)  NOT NULL,
	platform      VARCHAR(50)  NOT NULL,
	fingerprint   VARCHAR(512) NOT NULL UNIQUE,
	is_trusted    BOOLEAN      NOT NULL DEFAULT FALSE,
	last_seen_at  DATETIME,
	registered_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
	id                        INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id                   INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	session_id                INT          REFERENCES sessions (id) ON DELETE SET NULL,
	integration_connection_id INT          REFERENCES integration_connections (id) ON DELETE SET NULL,
	title                     VARCHAR(255) NOT NULL,
	description               TEXT,
	status                    VARCHAR(50)  NOT NULL DEFAULT 'pending',
	priority                  VARCHAR(50)  NOT NULL DEFAULT 'medium',
	due_at                    DATETIME,
	completed_at              DATETIME,
	created_at                DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at                DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at                DATETIME
);

CREATE TABLE IF NOT EXISTS workflows (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id        INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	name           VARCHAR(255) NOT NULL,
	description    TEXT,
	trigger_type   VARCHAR(50)  NOT NULL,
	trigger_config TEXT,
	steps          TEXT,
	is_active      BOOLEAN      NOT NULL DEFAULT FALSE,
	last_run_at    DATETIME,
	created_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at     DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workflow_runs (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	workflow_id INT         NOT NULL REFERENCES workflows (id) ON DELETE CASCADE,
	status      VARCHAR(50) NOT NULL DEFAULT 'running',
	started_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	finished_at DATETIME,
	error       TEXT,
	output      TEXT,
	created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id    INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	title      VARCHAR(255) NOT NULL,
	body       TEXT,
	type       VARCHAR(50)  NOT NULL DEFAULT 'info',
	source     VARCHAR(50)  NOT NULL DEFAULT 'system',
	is_read    BOOLEAN      NOT NULL DEFAULT FALSE,
	read_at    DATETIME,
	created_at DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS filesystem_scopes (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id      INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	assistant_id INT          NOT NULL REFERENCES assistants (id) ON DELETE CASCADE,
	path         VARCHAR(512) NOT NULL,
	access_level VARCHAR(50)  NOT NULL DEFAULT 'read',
	is_recursive BOOLEAN      NOT NULL DEFAULT FALSE,
	is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
	granted_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at   DATETIME,
	created_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS capability_policies (
	id                     INTEGER PRIMARY KEY AUTOINCREMENT,
	assistant_id           INT          NOT NULL REFERENCES assistants (id) ON DELETE CASCADE,
	capability             VARCHAR(100) NOT NULL,
	is_enabled             BOOLEAN      NOT NULL DEFAULT FALSE,
	require_confirmation   BOOLEAN      NOT NULL DEFAULT TRUE,
	confirmation_mode      VARCHAR(50)  NOT NULL DEFAULT 'always',
	max_executions_per_day INT,
	created_at             DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at             DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (assistant_id, capability)
);

CREATE TABLE IF NOT EXISTS agent_actions (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id       INT          NOT NULL REFERENCES users (id),
	assistant_id  INT          NOT NULL REFERENCES assistants (id),
	session_id    INT          REFERENCES sessions (id) ON DELETE SET NULL,
	workflow_id   INT          REFERENCES workflows (id) ON DELETE SET NULL,
	capability    VARCHAR(100) NOT NULL,
	action_type   VARCHAR(100) NOT NULL,
	target        TEXT,
	params        TEXT,
	was_confirmed BOOLEAN      NOT NULL DEFAULT FALSE,
	status        VARCHAR(50)  NOT NULL DEFAULT 'success',
	error         TEXT,
	executed_at   DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at    DATETIME  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_settings (
	id                 VARCHAR(36) PRIMARY KEY,
	installation_id    VARCHAR(36) NOT NULL UNIQUE,
	first_name         VARCHAR(120),
	last_name          VARCHAR(160),
	birth_date         DATE,
	display_name       VARCHAR(160),
	avatar_url         VARCHAR(500),
	email              VARCHAR(255),
	language           VARCHAR(10)  NOT NULL DEFAULT 'en',
	timezone           VARCHAR(80)  NOT NULL DEFAULT 'UTC',
	bio                TEXT,
	date_format        VARCHAR(30),
	time_format        VARCHAR(30),
	onboarding_complete INTEGER     NOT NULL DEFAULT 0,
	preferences        TEXT,
	created_at         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Configuración de usuario
CREATE INDEX IF NOT EXISTS idx_user_settings_installation_id  ON user_settings (installation_id);
-- Proveedores y modelos
CREATE INDEX IF NOT EXISTS idx_models_provider_id             ON models (provider_id);
-- Configuración de asistente
CREATE INDEX IF NOT EXISTS idx_assistant_configs_assistant_id ON assistant_configs (assistant_id);
CREATE INDEX IF NOT EXISTS idx_assistant_configs_model_id     ON assistant_configs (model_id);
-- Conexiones de canal
CREATE INDEX IF NOT EXISTS idx_channel_connections_channel_id ON channel_connections (channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_connections_user_id    ON channel_connections (user_id);
-- Contactos
CREATE INDEX IF NOT EXISTS idx_contacts_channel_connection_id ON contacts (channel_connection_id);
-- Sesiones
CREATE INDEX IF NOT EXISTS idx_sessions_user_id               ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_assistant_id          ON sessions (assistant_id);
CREATE INDEX IF NOT EXISTS idx_sessions_channel_connection_id ON sessions (channel_connection_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status                ON sessions (status);
-- Mensajes
CREATE INDEX IF NOT EXISTS idx_messages_session_id            ON messages (session_id);
CREATE INDEX IF NOT EXISTS idx_messages_model_id              ON messages (model_id);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at            ON messages (deleted_at) WHERE deleted_at IS NOT NULL;
-- Memoria
CREATE INDEX IF NOT EXISTS idx_memories_user_id               ON memories (user_id);
CREATE INDEX IF NOT EXISTS idx_memories_assistant_id          ON memories (assistant_id);
CREATE INDEX IF NOT EXISTS idx_memories_is_active             ON memories (is_active) WHERE is_active = TRUE;
-- Conexiones de integración
CREATE INDEX IF NOT EXISTS idx_integration_connections_integration_id ON integration_connections (integration_id);
CREATE INDEX IF NOT EXISTS idx_integration_connections_user_id        ON integration_connections (user_id);
-- Asignación de skills
CREATE INDEX IF NOT EXISTS idx_skill_assignments_assistant_id ON skill_assignments (assistant_id);
-- Tareas
CREATE INDEX IF NOT EXISTS idx_tasks_user_id                  ON tasks (user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status                   ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at               ON tasks (deleted_at) WHERE deleted_at IS NOT NULL;
-- Workflows y ejecuciones
CREATE INDEX IF NOT EXISTS idx_workflow_runs_workflow_id      ON workflow_runs (workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_runs_status           ON workflow_runs (status);
-- Notificaciones
CREATE INDEX IF NOT EXISTS idx_notifications_user_id          ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read          ON notifications (is_read) WHERE is_read = FALSE;
-- Filesystem scopes
CREATE INDEX IF NOT EXISTS idx_filesystem_scopes_user_id      ON filesystem_scopes (user_id);
CREATE INDEX IF NOT EXISTS idx_filesystem_scopes_assistant_id ON filesystem_scopes (assistant_id);
-- Auditoría
CREATE INDEX IF NOT EXISTS idx_audits_user_id                 ON audits (user_id);
-- Acciones del agente
CREATE INDEX IF NOT EXISTS idx_agent_actions_user_id          ON agent_actions (user_id);
CREATE INDEX IF NOT EXISTS idx_agent_actions_assistant_id     ON agent_actions (assistant_id);
CREATE INDEX IF NOT EXISTS idx_agent_actions_session_id       ON agent_actions (session_id);
CREATE INDEX IF NOT EXISTS idx_agent_actions_capability       ON agent_actions (capability);
CREATE INDEX IF NOT EXISTS idx_agent_actions_executed_at      ON agent_actions (executed_at);
