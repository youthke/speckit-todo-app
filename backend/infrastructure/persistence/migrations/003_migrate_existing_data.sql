-- Migration: Migrate existing data to DDD structure
-- Date: 2025-09-28
-- Description: Migrates existing tasks to new DDD schema with User context

-- First, create a default user for existing tasks
INSERT OR IGNORE INTO users (
    id,
    email,
    first_name,
    last_name,
    timezone,
    preferences,
    created_at,
    updated_at
) VALUES (
    1,
    'default@todo-app.local',
    'Default',
    'User',
    'UTC',
    '{"default_task_priority": "medium", "email_notifications": true, "theme_preference": "auto"}',
    datetime('now'),
    datetime('now')
);

-- Migrate existing tasks if the old structure exists
-- This assumes there might be an existing 'tasks' table with different structure
-- If the table doesn't exist or has different columns, this will be ignored

-- Check if old tasks table exists and has different structure
-- If so, migrate the data
INSERT OR IGNORE INTO tasks (
    title,
    description,
    status,
    priority,
    user_id,
    created_at,
    updated_at
)
SELECT
    COALESCE(Title, 'Untitled Task') as title,
    COALESCE(Description, '') as description,
    CASE
        WHEN Completed = 1 THEN 'completed'
        ELSE 'pending'
    END as status,
    'medium' as priority, -- Default priority for migrated tasks
    1 as user_id, -- Assign to default user
    COALESCE(CreatedAt, datetime('now')) as created_at,
    COALESCE(UpdatedAt, datetime('now')) as updated_at
FROM old_tasks
WHERE EXISTS (
    SELECT name FROM sqlite_master
    WHERE type='table' AND name='old_tasks'
);

-- If there's an existing tasks table with the current schema, ensure it has proper defaults
UPDATE tasks
SET user_id = 1
WHERE user_id IS NULL OR user_id = 0;

UPDATE tasks
SET priority = 'medium'
WHERE priority IS NULL OR priority = '';

UPDATE tasks
SET status = CASE
    WHEN status IS NULL OR status = '' THEN 'pending'
    WHEN status = 'true' OR status = '1' THEN 'completed'
    WHEN status = 'false' OR status = '0' THEN 'pending'
    ELSE status
END;

-- Clean up any invalid data
DELETE FROM tasks WHERE title IS NULL OR title = '';
DELETE FROM tasks WHERE user_id NOT IN (SELECT id FROM users);