-- Migration: Create authentication_sessions table
-- Description: Session management for OAuth and traditional authentication
-- Created: 2025-09-29

CREATE TABLE authentication_sessions (
    id VARCHAR(255) PRIMARY KEY,                     -- Session identifier (UUID or JWT ID)
    user_id INTEGER NOT NULL,                        -- Reference to users table
    session_token TEXT NOT NULL UNIQUE,              -- JWT session token
    refresh_token TEXT,                              -- OAuth refresh token (encrypted)
    access_token TEXT,                               -- OAuth access token (encrypted)
    token_expires_at DATETIME,                       -- OAuth token expiration time
    session_expires_at DATETIME NOT NULL,           -- Session expiration (24 hours)
    last_activity DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Last user activity
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,    -- Session creation
    user_agent TEXT,                                 -- Browser/client info
    ip_address VARCHAR(45),                          -- Client IP (IPv4/IPv6)

    -- Foreign key constraint
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_session_token ON authentication_sessions(session_token);
CREATE INDEX idx_user_sessions ON authentication_sessions(user_id);
CREATE INDEX idx_session_expires ON authentication_sessions(session_expires_at);
CREATE INDEX idx_last_activity ON authentication_sessions(last_activity);

-- Index for cleanup jobs
CREATE INDEX idx_token_expires ON authentication_sessions(token_expires_at);

-- Comments for field usage:
-- id: Unique session identifier, can be UUID or JWT 'jti' claim
-- session_token: The actual JWT token stored as string
-- refresh_token: OAuth refresh token, encrypted before storage
-- access_token: OAuth access token, encrypted before storage
-- token_expires_at: When OAuth tokens expire (for refresh logic)
-- session_expires_at: When the entire session expires (24h default)
-- last_activity: Updated on each API call for session extension
-- user_agent: For security auditing and session management
-- ip_address: For security auditing and suspicious activity detection