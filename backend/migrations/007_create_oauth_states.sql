-- Migration: Create oauth_states table
-- Description: Temporary storage for OAuth flow state validation (CSRF protection)
-- Created: 2025-09-29

CREATE TABLE oauth_states (
    state_token VARCHAR(255) PRIMARY KEY,           -- Random state parameter for CSRF
    pkce_verifier VARCHAR(255) NOT NULL,            -- PKCE code verifier for security
    redirect_uri VARCHAR(1000) NOT NULL,            -- Post-auth redirect destination
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    expires_at DATETIME NOT NULL                     -- Expiration (5 minutes from creation)
);

-- Create index for cleanup operations
CREATE INDEX idx_oauth_expires ON oauth_states(expires_at);
CREATE INDEX idx_oauth_created ON oauth_states(created_at);

-- Comments for field usage:
-- state_token: Cryptographically random string (32+ characters)
--              Used as CSRF protection in OAuth flow
--              Passed to Google and verified on callback
-- pkce_verifier: PKCE (Proof Key for Code Exchange) code verifier
--                Enhances security for OAuth flow
--                Used to validate PKCE challenge
-- redirect_uri: Where to redirect user after successful authentication
--               Must be whitelisted in application configuration
--               Default to dashboard or original requested page
-- created_at: When this state was created
-- expires_at: When this state expires (5 minutes from creation)
--             States are automatically cleaned up by background job

-- Validation rules (enforced at application level):
-- 1. state_token must be cryptographically random (32+ chars)
-- 2. pkce_verifier must meet PKCE specification requirements
-- 3. expires_at must be exactly 5 minutes from created_at
-- 4. redirect_uri must be in application's allowed redirect list