CREATE TABLE IF NOT EXISTS "commands" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    raw VARCHAR(255) NOT NULL,
    status VARCHAR(255) NULL,
    error_msg TEXT NULL,
    is_deleted BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "command_logs" (
    command_id INTEGER UNIQUE,
    logs TEXT,
    FOREIGN KEY (command_id) REFERENCES commands (id) ON DELETE SET null
)