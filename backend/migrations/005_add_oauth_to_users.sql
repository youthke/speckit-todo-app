-- Migration: Add OAuth fields to users table
-- Description: Extends existing users table with Google OAuth authentication support
-- Created: 2025-09-29

-- Add OAuth-related columns to existing users table
ALTER TABLE users ADD COLUMN google_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN oauth_provider VARCHAR(50);
ALTER TABLE users ADD COLUMN oauth_created_at DATETIME;

-- Modify password_hash to be nullable for OAuth-only users
-- Note: This assumes password_hash exists and needs to be made nullable
-- ALTER TABLE users MODIFY COLUMN password_hash VARCHAR(255) NULL;

-- Create indexes for efficient OAuth lookups
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_oauth_provider ON users(oauth_provider);

-- Add constraints to ensure data integrity
-- Google ID must be unique when present
-- OAuth provider must be 'google' when google_id is present
-- Either password_hash OR google_id must be present (enforced at application level)

-- Sample validation queries (for reference, not executed):
-- SELECT * FROM users WHERE google_id IS NOT NULL AND oauth_provider != 'google'; -- Should return empty
-- SELECT * FROM users WHERE google_id IS NULL AND password_hash IS NULL; -- Should return empty