-- Migration: Google OAuth Signup Support
-- Description: Adds Google OAuth authentication with separate google_identities table
-- Feature: 007-google
-- Created: 2025-10-01

-- Up Migration
-- Add auth_method column to users table
ALTER TABLE users ADD COLUMN auth_method VARCHAR(50) NOT NULL DEFAULT 'password';

-- Create google_identities table
CREATE TABLE IF NOT EXISTS google_identities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    google_user_id VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_google_user_id ON google_identities(google_user_id);
CREATE INDEX idx_google_email ON google_identities(email);

-- Down Migration (commented for reference)
-- DROP INDEX IF EXISTS idx_google_email;
-- DROP INDEX IF EXISTS idx_google_user_id;
-- DROP TABLE IF EXISTS google_identities;
-- ALTER TABLE users DROP COLUMN auth_method;
