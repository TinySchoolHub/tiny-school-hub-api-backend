-- Create absences table
CREATE TABLE IF NOT EXISTS absences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_name VARCHAR(255) NOT NULL,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    absence_date DATE NOT NULL,
    reported_by VARCHAR(20) NOT NULL CHECK (reported_by IN ('TEACHER', 'PARENT')),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'ACKED')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for efficient lookups
CREATE INDEX idx_absences_class_id ON absences(class_id);
CREATE INDEX idx_absences_reporter_id ON absences(reporter_id);
CREATE INDEX idx_absences_absence_date ON absences(absence_date DESC);
CREATE INDEX idx_absences_status ON absences(status);
