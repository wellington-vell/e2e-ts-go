CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'todo_status') THEN
        CREATE TYPE todo_status AS ENUM ('backlog', 'todo', 'in_progress', 'done', 'canceled');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'todo_label') THEN
        CREATE TYPE todo_label AS ENUM ('bug', 'feature', 'doc');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'todo_priority') THEN
        CREATE TYPE todo_priority AS ENUM ('low', 'medium', 'high');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    text TEXT NOT NULL,
    status todo_status NOT NULL DEFAULT 'backlog',
    label todo_label,
    priority todo_priority NOT NULL DEFAULT 'medium',
    estimated_hours DECIMAL(10, 2) DEFAULT 0,
    actual_hours DECIMAL(10, 2) DEFAULT 0,
    progress INTEGER NOT NULL DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    cost DECIMAL(12, 2) DEFAULT 0,
    due_date TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
CREATE INDEX IF NOT EXISTS idx_todos_label ON todos(label);
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
