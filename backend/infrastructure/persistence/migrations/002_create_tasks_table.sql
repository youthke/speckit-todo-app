-- Migration: Create tasks table for DDD Task Management
-- Date: 2025-09-28
-- Description: Creates the tasks table with DDD value objects support

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(10) NOT NULL DEFAULT 'medium',
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_user_status ON tasks(user_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_user_priority ON tasks(user_id, priority);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- Add constraints for status values
CREATE TRIGGER IF NOT EXISTS check_task_status
    BEFORE INSERT ON tasks
    FOR EACH ROW
    WHEN NEW.status NOT IN ('pending', 'completed', 'archived')
    BEGIN
        SELECT RAISE(ABORT, 'Invalid task status. Must be: pending, completed, or archived');
    END;

CREATE TRIGGER IF NOT EXISTS check_task_status_update
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    WHEN NEW.status NOT IN ('pending', 'completed', 'archived')
    BEGIN
        SELECT RAISE(ABORT, 'Invalid task status. Must be: pending, completed, or archived');
    END;

-- Add constraints for priority values
CREATE TRIGGER IF NOT EXISTS check_task_priority
    BEFORE INSERT ON tasks
    FOR EACH ROW
    WHEN NEW.priority NOT IN ('low', 'medium', 'high')
    BEGIN
        SELECT RAISE(ABORT, 'Invalid task priority. Must be: low, medium, or high');
    END;

CREATE TRIGGER IF NOT EXISTS check_task_priority_update
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    WHEN NEW.priority NOT IN ('low', 'medium', 'high')
    BEGIN
        SELECT RAISE(ABORT, 'Invalid task priority. Must be: low, medium, or high');
    END;