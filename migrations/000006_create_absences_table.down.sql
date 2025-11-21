-- Drop absences table
DROP INDEX IF EXISTS idx_absences_status;
DROP INDEX IF EXISTS idx_absences_absence_date;
DROP INDEX IF EXISTS idx_absences_reporter_id;
DROP INDEX IF EXISTS idx_absences_class_id;
DROP TABLE IF EXISTS absences;
