-- Migration: Create jobs table
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type            VARCHAR(50) NOT NULL,
    payload         JSONB NOT NULL,
    priority        VARCHAR(10) NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
    status          VARCHAR(20) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'processing', 'success', 'failed', 'terminated')),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP,
    result          TEXT,
    error_message   TEXT,
    retries         INT NOT NULL DEFAULT 0
);