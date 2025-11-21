-- Create class_members table
CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    role_in_class VARCHAR(20) NOT NULL CHECK (role_in_class IN ('TEACHER', 'PARENT')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, class_id)
);

-- Create indexes for efficient lookups
CREATE INDEX idx_class_members_user_id ON class_members(user_id);
CREATE INDEX idx_class_members_class_id ON class_members(class_id);
CREATE INDEX idx_class_members_role ON class_members(role_in_class);
