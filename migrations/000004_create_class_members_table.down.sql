-- Drop class_members table
DROP INDEX IF EXISTS idx_class_members_role;
DROP INDEX IF EXISTS idx_class_members_class_id;
DROP INDEX IF EXISTS idx_class_members_user_id;
DROP TABLE IF EXISTS class_members;
