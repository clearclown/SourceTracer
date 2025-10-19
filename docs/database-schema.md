# SourceTracer Database Schema Design

## Database Architecture

3層データストア戦略：

1. **PostgreSQL** - 構造化データ、メタデータ、リレーショナル
2. **MongoDB** - 非構造化データ、スクレイピング結果
3. **Redis** - キャッシュ、セッション、一時データ

---

## PostgreSQL Schema

### Extensions

```sql
-- Vector search
CREATE EXTENSION IF NOT EXISTS vector;

-- UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Full-text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

---

### 1. users (将来的な認証用)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    plan VARCHAR(20) DEFAULT 'free', -- free, pro, enterprise
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    
    -- API usage tracking
    api_key_hash VARCHAR(255),
    rate_limit_tier INTEGER DEFAULT 100,
    
    -- Preferences
    preferred_llm VARCHAR(20) DEFAULT 'claude',
    default_min_credibility DECIMAL(3,2) DEFAULT 0.50,
    
    INDEX idx_users_email (email),
    INDEX idx_users_api_key (api_key_hash)
);
```

---

### 2. analyses

分析セッションのメタデータ

```sql
CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    -- Input
    input_text TEXT NOT NULL,
    input_hash VARCHAR(64) NOT NULL, -- SHA256 for deduplication
    language VARCHAR(5) DEFAULT 'en',
    
    -- Processing
    status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
    llm_provider VARCHAR(20),
    llm_model VARCHAR(50),
    
    -- Results
    claims_count INTEGER DEFAULT 0,
    evidences_count INTEGER DEFAULT 0,
    overall_credibility DECIMAL(3,2),
    
    -- Metadata
    processing_time_ms INTEGER,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Search options (JSONB for flexibility)
    options JSONB,
    
    INDEX idx_analyses_user_id (user_id),
    INDEX idx_analyses_created_at (created_at DESC),
    INDEX idx_analyses_status (status),
    INDEX idx_analyses_input_hash (input_hash)
);
```

---

### 3. claims

抽出された主張

```sql
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    
    -- Claim content
    text TEXT NOT NULL,
    type VARCHAR(20) NOT NULL, -- opinion, fact, mixed, unclear
    confidence DECIMAL(3,2) NOT NULL,
    
    -- Position in original text
    start_position INTEGER NOT NULL,
    end_position INTEGER NOT NULL,
    context TEXT, -- surrounding text
    
    -- Analysis results
    recommendation VARCHAR(20), -- SUPPORTED, UNSUPPORTED, MIXED, INSUFFICIENT
    summary TEXT,
    
    -- Credibility
    overall_credibility DECIMAL(3,2),
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_claims_analysis_id (analysis_id),
    INDEX idx_claims_type (type),
    INDEX idx_claims_credibility (overall_credibility DESC)
);
```

---

### 4. evidences

発見されたエビデンス

```sql
CREATE TABLE evidences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    
    -- Source info
    source_name VARCHAR(100) NOT NULL, -- Semantic Scholar, Google, etc.
    source_type VARCHAR(50) NOT NULL, -- peer_reviewed_journal, preprint, news, blog
    url TEXT NOT NULL,
    
    -- Content
    title TEXT NOT NULL,
    snippet TEXT,
    authors JSONB, -- ["Smith, J.", "Doe, A."]
    
    -- Metadata
    published_date DATE,
    doi VARCHAR(255),
    isbn VARCHAR(20),
    journal_name VARCHAR(255),
    
    -- Metrics
    credibility_score DECIMAL(3,2) NOT NULL,
    relevance_score DECIMAL(3,2) NOT NULL,
    citation_count INTEGER DEFAULT 0,
    
    -- Flags
    is_counter_evidence BOOLEAN DEFAULT FALSE,
    is_paywalled BOOLEAN DEFAULT FALSE,
    
    -- Raw data reference (MongoDB)
    raw_data_id VARCHAR(24), -- MongoDB ObjectId
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_evidences_claim_id (claim_id),
    INDEX idx_evidences_source_type (source_type),
    INDEX idx_evidences_credibility (credibility_score DESC),
    INDEX idx_evidences_url (url)
);

-- URL uniqueness constraint per claim
CREATE UNIQUE INDEX idx_evidences_claim_url ON evidences(claim_id, url);
```

---

### 5. sources

利用可能な検索ソースの設定

```sql
CREATE TABLE sources (
    id VARCHAR(50) PRIMARY KEY, -- semantic_scholar, arxiv, etc.
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- academic, preprint, news, general
    
    -- Configuration
    api_endpoint TEXT,
    requires_api_key BOOLEAN DEFAULT FALSE,
    rate_limit_per_minute INTEGER,
    rate_limit_per_day INTEGER,
    
    -- Credibility
    base_credibility_weight DECIMAL(3,2) NOT NULL,
    
    -- Status
    is_enabled BOOLEAN DEFAULT TRUE,
    last_error TEXT,
    last_success_at TIMESTAMP,
    
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Seed data
INSERT INTO sources (id, name, type, base_credibility_weight, rate_limit_per_minute) VALUES
    ('semantic_scholar', 'Semantic Scholar', 'academic', 0.95, 100),
    ('arxiv', 'arXiv', 'preprint', 0.85, 0),
    ('pubmed', 'PubMed', 'academic', 0.95, 10),
    ('google', 'Google Search', 'general', 0.60, 10),
    ('duckduckgo', 'DuckDuckGo', 'general', 0.55, 0);
```

---

### 6. search_cache

検索結果のキャッシュ（重複検索回避）

```sql
CREATE TABLE search_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Cache key
    query_text TEXT NOT NULL,
    source_id VARCHAR(50) NOT NULL REFERENCES sources(id),
    query_hash VARCHAR(64) NOT NULL, -- SHA256 of normalized query
    
    -- Results (JSON array of evidence metadata)
    results JSONB NOT NULL,
    results_count INTEGER DEFAULT 0,
    
    -- Expiration
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    hit_count INTEGER DEFAULT 0,
    
    INDEX idx_search_cache_query_hash (query_hash),
    INDEX idx_search_cache_expires_at (expires_at)
);

-- Unique constraint on query + source
CREATE UNIQUE INDEX idx_search_cache_unique ON search_cache(query_hash, source_id);
```

---

### 7. jobs

非同期ジョブ管理

```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Job details
    type VARCHAR(50) NOT NULL, -- analyze, batch_analyze
    status VARCHAR(20) DEFAULT 'queued', -- queued, processing, completed, failed
    priority INTEGER DEFAULT 5, -- 1-10
    
    -- Input
    input_data JSONB NOT NULL,
    
    -- Progress
    progress INTEGER DEFAULT 0, -- 0-100
    current_step TEXT,
    
    -- Result
    result_id UUID REFERENCES analyses(id),
    error_message TEXT,
    
    -- Timing
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    -- Callback
    callback_url TEXT,
    
    INDEX idx_jobs_status (status),
    INDEX idx_jobs_user_id (user_id),
    INDEX idx_jobs_created_at (created_at DESC)
);
```

---

### 8. embeddings (Vector Search)

テキストのベクトル埋め込み（意味検索用）

```sql
CREATE TABLE embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Target reference
    target_type VARCHAR(50) NOT NULL, -- claim, evidence, analysis
    target_id UUID NOT NULL,
    
    -- Embedding
    model VARCHAR(50) NOT NULL, -- text-embedding-3-small, etc.
    vector vector(1536), -- OpenAI: 1536, other models may differ
    
    -- Metadata
    text_snippet TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_embeddings_target (target_type, target_id)
);

-- Vector similarity search index
CREATE INDEX idx_embeddings_vector ON embeddings USING ivfflat (vector vector_cosine_ops);
```

---

### 9. audit_logs

システム監査ログ

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Event
    event_type VARCHAR(50) NOT NULL, -- api_request, analysis_created, etc.
    resource_type VARCHAR(50), -- analysis, evidence
    resource_id UUID,
    
    -- Details
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_audit_logs_user_id (user_id),
    INDEX idx_audit_logs_event_type (event_type),
    INDEX idx_audit_logs_created_at (created_at DESC)
);
```

---

## MongoDB Collections

### 1. raw_scrapes

Playwrightで取得した生データ

```javascript
{
  _id: ObjectId("..."),
  
  // Source
  url: "https://example.com/article",
  source_engine: "duckduckgo",
  search_query: "climate change evidence",
  
  // Content
  html: "<html>...</html>",
  text_content: "Extracted text...",
  screenshot_base64: "data:image/png;base64,...", // optional
  
  // Metadata
  http_status: 200,
  response_headers: {
    "content-type": "text/html",
    // ...
  },
  
  // Timestamps
  scraped_at: ISODate("2025-10-19T10:30:00Z"),
  expires_at: ISODate("2025-10-26T10:30:00Z"), // 7 days TTL
  
  // References
  evidence_id: "uuid-from-postgresql", // optional
}

// Indexes
db.raw_scrapes.createIndex({ url: 1, scraped_at: -1 })
db.raw_scrapes.createIndex({ expires_at: 1 }, { expireAfterSeconds: 0 }) // TTL
db.raw_scrapes.createIndex({ search_query: 1 })
```

---

### 2. llm_responses

LLMレスポンスのキャッシュ

```javascript
{
  _id: ObjectId("..."),
  
  // Request
  provider: "claude",
  model: "claude-sonnet-4-20250514",
  prompt_hash: "sha256...",
  prompt: "Analyze this text...",
  
  // Response
  response: {
    text: "This is an opinion because...",
    usage: {
      prompt_tokens: 150,
      completion_tokens: 80,
      total_tokens: 230
    }
  },
  
  // Cache control
  created_at: ISODate("2025-10-19T10:30:00Z"),
  expires_at: ISODate("2025-10-20T10:30:00Z"), // 24 hours
  hit_count: 3,
  
  // Metadata
  temperature: 0.3,
  max_tokens: 4000
}

// Indexes
db.llm_responses.createIndex({ prompt_hash: 1, provider: 1 }, { unique: true })
db.llm_responses.createIndex({ expires_at: 1 }, { expireAfterSeconds: 0 })
```

---

### 3. analysis_reports

詳細な分析レポート（大容量テキスト）

```javascript
{
  _id: ObjectId("..."),
  analysis_id: "uuid-from-postgresql",
  
  // Full report
  markdown_report: "# Analysis Report\n\n...",
  html_report: "<h1>Analysis Report</h1>...",
  
  // Sections
  sections: [
    {
      title: "Executive Summary",
      content: "...",
      order: 1
    },
    {
      title: "Detailed Findings",
      content: "...",
      order: 2
    }
  ],
  
  // Export formats
  exports: {
    pdf_url: "https://storage.googleapis.com/...",
    csv_url: "https://storage.googleapis.com/..."
  },
  
  created_at: ISODate("2025-10-19T10:30:00Z")
}

// Indexes
db.analysis_reports.createIndex({ analysis_id: 1 }, { unique: true })
```

---

## Redis Keys Structure

### 1. LLM Response Cache

```
llm_cache:{provider}:{prompt_hash} → JSON response
TTL: 24 hours
```

### 2. Search Results Cache

```
search_cache:{source}:{query_hash} → JSON results
TTL: 1 hour
```

### 3. Rate Limiting

```
ratelimit:{user_id}:{endpoint} → counter
TTL: 1 minute (rolling window)
```

### 4. Job Queue

```
queue:factcheck_jobs → List of job IDs
job:{job_id}:status → JSON job status
job:{job_id}:progress → Integer (0-100)
```

### 5. Session Management

```
session:{session_id} → JSON user session
TTL: 7 days
```

### 6. Browser Pool

```
browser_pool:available → Set of browser instance IDs
browser:{instance_id}:metadata → JSON browser info
```

---

## Migration Scripts

### Initial Setup

```sql
-- migrations/001_initial_schema.up.sql

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create tables (copy from above)
-- ...

-- Seed data
INSERT INTO sources (...) VALUES (...);
```

### Rollback

```sql
-- migrations/001_initial_schema.down.sql

DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS embeddings CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
DROP TABLE IF EXISTS search_cache CASCADE;
DROP TABLE IF EXISTS evidences CASCADE;
DROP TABLE IF EXISTS claims CASCADE;
DROP TABLE IF EXISTS analyses CASCADE;
DROP TABLE IF EXISTS sources CASCADE;
DROP TABLE IF EXISTS users CASCADE;
```

---

## Query Examples

### Find all analyses with high credibility

```sql
SELECT 
    a.id,
    a.input_text,
    a.overall_credibility,
    COUNT(c.id) as claims_count
FROM analyses a
LEFT JOIN claims c ON c.analysis_id = a.id
WHERE a.overall_credibility > 0.8
  AND a.status = 'completed'
GROUP BY a.id
ORDER BY a.created_at DESC;
```

### Get evidences with sources for a claim

```sql
SELECT 
    e.title,
    e.url,
    e.credibility_score,
    e.source_type,
    s.name as source_name
FROM evidences e
JOIN sources s ON e.source_name = s.id
WHERE e.claim_id = 'claim-uuid'
ORDER BY e.credibility_score DESC;
```

### Semantic search for similar claims (Vector)

```sql
WITH query_embedding AS (
    SELECT vector FROM embeddings
    WHERE target_id = 'source-claim-uuid'
    LIMIT 1
)
SELECT 
    c.text,
    e.vector <=> (SELECT vector FROM query_embedding) as distance
FROM claims c
JOIN embeddings e ON e.target_id = c.id AND e.target_type = 'claim'
ORDER BY distance
LIMIT 10;
```

---

## Backup Strategy

### PostgreSQL
```bash
# Daily backup
pg_dump -U factcheck -d factcheck > backup_$(date +%Y%m%d).sql

# Restore
psql -U factcheck -d factcheck < backup_20251019.sql
```

### MongoDB
```bash
# Daily backup
mongodump --uri="mongodb://localhost:27017/factcheck_raw" --out=/backup/mongo_$(date +%Y%m%d)

# Restore
mongorestore --uri="mongodb://localhost:27017/factcheck_raw" /backup/mongo_20251019
```

### Redis
```bash
# Automated RDB snapshots (redis.conf)
save 900 1
save 300 10
save 60 10000
```

---

## Performance Optimization

1. **Indexes**: 適切なインデックス作成済み
2. **Partitioning**: `analyses` テーブルは日付でパーティション（大規模時）
3. **Connection Pooling**: pgBouncer使用推奨
4. **Materialized Views**: 頻繁なクエリ用
5. **Read Replicas**: GCP Cloud SQL でレプリカ追加

---

## Data Retention Policy

- **analyses**: 1年後にアーカイブ（Deleted users: 30日後）
- **search_cache**: 24時間
- **raw_scrapes**: 7日間（MongoDB TTL）
- **llm_responses**: 24時間
- **audit_logs**: 3年保存
