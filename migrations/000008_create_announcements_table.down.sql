-- Drop announcements table
DROP INDEX IF EXISTS idx_announcements_publish_at;
DROP INDEX IF EXISTS idx_announcements_author_id;
DROP INDEX IF EXISTS idx_announcements_class_id;
DROP TABLE IF EXISTS announcements;
