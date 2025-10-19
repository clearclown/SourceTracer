-- SourceTracer Database Schema v1
-- PostgreSQL 14+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (for future authentication)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Analysis sessions
CREATE TABLE IF NOT EXISTS analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    original_text TEXT NOT NULL,
    overall_credibility DECIMAL(3,2) CHECK (overall_credibility >= 0 AND overall_credibility <= 1),
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Claims extracted from analyses
CREATE TABLE IF NOT EXISTS claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    claim_type VARCHAR(20) NOT NULL CHECK (claim_type IN ('opinion', 'fact', 'mixed', 'unclear')),
    confidence DECIMAL(3,2) CHECK (confidence >= 0 AND confidence <= 1),
    position_start INTEGER,
    position_end INTEGER,
    summary TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Evidence supporting or contradicting claims
CREATE TABLE IF NOT EXISTS evidences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    source VARCHAR(100) NOT NULL,
    title TEXT NOT NULL,
    authors TEXT[], -- Array of author names
    url TEXT,
    snippet TEXT,
    credibility DECIMAL(3,2) CHECK (credibility >= 0 AND credibility <= 1),
    relevance DECIMAL(3,2) CHECK (relevance >= 0 AND relevance <= 1),
    published_date DATE,
    source_type VARCHAR(50) NOT NULL,
    citation_count INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_analyses_user_id ON analyses(user_id);
CREATE INDEX idx_analyses_created_at ON analyses(created_at DESC);
CREATE INDEX idx_claims_analysis_id ON claims(analysis_id);
CREATE INDEX idx_claims_claim_type ON claims(claim_type);
CREATE INDEX idx_evidences_claim_id ON evidences(claim_id);
CREATE INDEX idx_evidences_source_type ON evidences(source_type);

-- Comments for documentation
COMMENT ON TABLE analyses IS 'Text analysis sessions with overall metrics';
COMMENT ON TABLE claims IS 'Individual claims extracted from analyzed text';
COMMENT ON TABLE evidences IS 'Supporting or contradicting evidence for claims';
COMMENT ON COLUMN claims.claim_type IS 'Type: opinion, fact, mixed, or unclear';
COMMENT ON COLUMN evidences.source_type IS 'Type: peer_reviewed_journal, preprint, news, blog, etc.';
