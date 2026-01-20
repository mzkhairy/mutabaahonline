-- Enable Extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. TABLE: INSTITUTIONS
CREATE TABLE IF NOT EXISTS institutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT idx_institutions_code_unique UNIQUE (code)
);

-- 2. TABLE: USERS
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID REFERENCES institutions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    username VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255), 
    password_hash VARCHAR(255) NOT NULL,
    serial_number VARCHAR(50),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    must_change_password BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_username_per_institution UNIQUE(institution_id, username),
    CONSTRAINT unique_serial_per_institution UNIQUE(institution_id, serial_number) 
);

-- 3. TABLE: ACADEMIC_YEARS
CREATE TABLE IF NOT EXISTS academic_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID REFERENCES institutions(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. TABLE: CLASSES
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID REFERENCES institutions(id) ON DELETE CASCADE NOT NULL,
    academic_year_id UUID REFERENCES academic_years(id) ON DELETE CASCADE NOT NULL,
    teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    level VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. TABLE: CLASS_STUDENTS
CREATE TABLE IF NOT EXISTS class_students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE NOT NULL,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(class_id, student_id)
);

-- 6. TABLE: MUTABAAH_TEMPLATES
CREATE TABLE IF NOT EXISTS mutabaah_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id UUID REFERENCES institutions(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(100) NOT NULL,
    structure JSONB NOT NULL DEFAULT '[]', 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. TABLE: SESSIONS
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE NOT NULL,
    template_id UUID REFERENCES mutabaah_templates(id) NOT NULL,
    date DATE NOT NULL,
    name VARCHAR(150) NOT NULL,
    class_activity_data JSONB DEFAULT '{}',
    is_student_input_allowed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. TABLE: MUTABAAH_ENTRIES
CREATE TABLE IF NOT EXISTS mutabaah_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE NOT NULL,
    student_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    attendance VARCHAR(20) DEFAULT 'HADIR',
    status VARCHAR(50),
    student_activity_data JSONB DEFAULT '{}',
    note TEXT,
    is_locked BOOLEAN DEFAULT FALSE,
    last_updated_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, student_id)
);