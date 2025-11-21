-- Drop photos table
DROP INDEX IF EXISTS idx_photos_created_at;
DROP INDEX IF EXISTS idx_photos_uploader_id;
DROP INDEX IF EXISTS idx_photos_class_id;
DROP TABLE IF EXISTS photos;
