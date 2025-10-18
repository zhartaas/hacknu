-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on created_at for faster sorting
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Insert some sample data for development
INSERT INTO users (id, email, name, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'john.doe@example.com', 'John Doe', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440001', 'jane.smith@example.com', 'Jane Smith', NOW(), NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'bob.wilson@example.com', 'Bob Wilson', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
