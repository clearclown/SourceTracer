# SourceTracer API Design v1.0

## Base URL

```
Development: http://localhost:8080/api/v1
Production:  https://api.sourcetracer.org/api/v1
```

---

## Authentication

現在のバージョンでは認証なし（OSS版）。将来的にJWT実装予定。

**Production版（GCP）:**
```
Authorization: Bearer <jwt_token>
```

---

## Common Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "metadata": {
    "request_id": "req_abc123",
    "timestamp": "2025-10-19T10:30:00Z",
    "processing_time_ms": 1234
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "INVALID_INPUT",
    "message": "Text input is required",
    "details": {
      "field": "text",
      "reason": "cannot be empty"
    }
  },
  "metadata": {
    "request_id": "req_abc123",
    "timestamp": "2025-10-19T10:30:00Z"
  }
}
```

### Error Codes
- `INVALID_INPUT` - バリデーションエラー
- `LLM_ERROR` - LLMプロバイダーエラー
- `SEARCH_ERROR` - 検索APIエラー
- `DATABASE_ERROR` - DB接続エラー
- `RATE_LIMIT_EXCEEDED` - レート制限超過
- `INTERNAL_ERROR` - サーバー内部エラー

---

## Endpoints

### 1. Health Check

**GET** `/health`

システムヘルスチェック

**Response:**
```json
{
  "status": "ok",
  "services": {
    "database": "ok",
    "mongodb": "ok",
    "redis": "ok",
    "llm_providers": {
      "openai": "ok",
      "claude": "ok",
      "deepseek": "error"
    }
  },
  "version": "1.0.0",
  "uptime_seconds": 12345
}
```

---

### 2. Text Analysis

**POST** `/analyze`

テキスト分析とエビデンス収集のメインエンドポイント

**Request:**
```json
{
  "text": "Climate change is accelerating faster than predicted.",
  "options": {
    "llm_provider": "claude",  // openai, claude, deepseek (optional)
    "min_credibility": 0.7,    // 0.0-1.0 (optional, default: 0.5)
    "max_results": 10,         // 1-50 (optional, default: 10)
    "search_engines": [        // optional, default: all
      "semantic_scholar",
      "arxiv",
      "google",
      "duckduckgo"
    ],
    "include_opinions": true,  // 意見も含めるか (optional, default: true)
    "language": "ja"           // ja, en, zh (optional, auto-detect)
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "analysis_id": "anl_xyz789",
    "claims": [
      {
        "id": "clm_001",
        "text": "Climate change is accelerating",
        "type": "fact",              // opinion, fact, mixed, unclear
        "confidence": 0.92,
        "position": {
          "start": 0,
          "end": 33
        },
        "evidences": [
          {
            "id": "evd_001",
            "source": "Semantic Scholar",
            "title": "Accelerating climate change impacts",
            "authors": ["Smith, J.", "Doe, A."],
            "url": "https://doi.org/10.1234/example",
            "snippet": "Recent data shows climate change accelerating...",
            "credibility": 0.95,
            "relevance": 0.88,
            "published_date": "2024-03-15",
            "source_type": "peer_reviewed_journal",
            "citation_count": 47,
            "metadata": {
              "journal": "Nature Climate Change",
              "doi": "10.1234/example"
            }
          }
        ],
        "counter_evidences": [
          {
            "id": "evd_002",
            "source": "Blog Post",
            "title": "Climate data questionable",
            "url": "https://example.com/blog",
            "credibility": 0.25,
            "relevance": 0.60,
            "published_date": "2024-01-10",
            "source_type": "blog"
          }
        ],
        "summary": "Strong consensus supports acceleration of climate change based on peer-reviewed research.",
        "recommendation": "SUPPORTED"  // SUPPORTED, UNSUPPORTED, MIXED, INSUFFICIENT
      }
    ],
    "overall_credibility": 0.87,
    "processing_time_ms": 3420
  },
  "metadata": {
    "request_id": "req_abc123",
    "timestamp": "2025-10-19T10:30:00Z",
    "llm_used": "claude-sonnet-4",
    "sources_searched": ["semantic_scholar", "arxiv", "google"]
  }
}
```

---

### 3. Async Analysis (Background Job)

**POST** `/analyze/async`

大量テキストや複数ドキュメント分析用の非同期ジョブ

**Request:**
```json
{
  "text": "...",  // または "file_url"
  "options": { /* same as /analyze */ },
  "callback_url": "https://your-app.com/webhook"  // optional
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "job_id": "job_abc123",
    "status": "queued",  // queued, processing, completed, failed
    "estimated_time_seconds": 120
  }
}
```

---

### 4. Job Status

**GET** `/jobs/{job_id}`

非同期ジョブの進捗確認

**Response:**
```json
{
  "success": true,
  "data": {
    "job_id": "job_abc123",
    "status": "processing",
    "progress": 45,  // percentage
    "created_at": "2025-10-19T10:30:00Z",
    "updated_at": "2025-10-19T10:31:30Z",
    "result": null  // completed時に結果が入る
  }
}
```

---

### 5. Search Sources

**GET** `/sources`

利用可能な検索ソース一覧

**Response:**
```json
{
  "success": true,
  "data": {
    "sources": [
      {
        "id": "semantic_scholar",
        "name": "Semantic Scholar",
        "type": "academic",
        "credibility_weight": 0.95,
        "status": "available",
        "rate_limit": "100/min"
      },
      {
        "id": "arxiv",
        "name": "arXiv",
        "type": "preprint",
        "credibility_weight": 0.85,
        "status": "available",
        "rate_limit": "unlimited"
      },
      {
        "id": "google",
        "name": "Google Custom Search",
        "type": "general",
        "credibility_weight": 0.60,
        "status": "available",
        "rate_limit": "100/day"
      }
    ]
  }
}
```

---

### 6. Claim Classification

**POST** `/classify`

単一の主張を分類（意見/事実判定のみ）

**Request:**
```json
{
  "text": "Python is the best programming language.",
  "llm_provider": "openai"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "type": "opinion",
    "confidence": 0.98,
    "reasoning": "Contains subjective judgment ('best') without measurable criteria."
  }
}
```

---

### 7. Evidence Search

**POST** `/evidence/search`

特定のクエリでエビデンス検索のみ実行

**Request:**
```json
{
  "query": "effects of caffeine on sleep",
  "sources": ["pubmed", "semantic_scholar"],
  "max_results": 5,
  "min_credibility": 0.8
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "evidences": [ /* Evidence objects */ ],
    "total_found": 47,
    "returned": 5
  }
}
```

---

### 8. Credibility Scoring

**POST** `/credibility/score`

URLまたはドメインの信頼性スコア計算

**Request:**
```json
{
  "url": "https://www.nature.com/articles/example",
  "context": "Climate change research"  // optional
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "url": "https://www.nature.com/articles/example",
    "credibility_score": 0.96,
    "factors": {
      "domain_reputation": 0.98,
      "source_type": "peer_reviewed_journal",
      "citation_count": 123,
      "recency_score": 0.85,
      "author_reputation": 0.92
    },
    "recommendation": "highly_credible"
  }
}
```

---

### 9. Analysis History

**GET** `/history`

ユーザーの分析履歴（将来的な機能）

**Query Parameters:**
- `limit` (default: 20)
- `offset` (default: 0)
- `sort` (default: created_at_desc)

**Response:**
```json
{
  "success": true,
  "data": {
    "analyses": [
      {
        "id": "anl_xyz789",
        "text_preview": "Climate change is accelerating...",
        "created_at": "2025-10-19T10:30:00Z",
        "claims_count": 3,
        "overall_credibility": 0.87
      }
    ],
    "total": 150,
    "limit": 20,
    "offset": 0
  }
}
```

---

### 10. Export Report

**GET** `/analyze/{analysis_id}/export`

分析結果をエクスポート

**Query Parameters:**
- `format`: `json`, `markdown`, `pdf`, `csv`

**Response:**
- `json`: JSON形式
- `markdown`: Markdown形式の詳細レポート
- `pdf`: PDFファイル（バイナリ）
- `csv`: 引用リスト CSV

---

## Rate Limiting

**Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1634567890
```

**Limits:**
- 未認証: 100 requests/hour
- 認証済み（OSS自己ホスト）: 1000 requests/hour
- Production（有料）: 10000 requests/hour

---

## Pagination

リスト系エンドポイントは以下のクエリパラメータをサポート：

```
?limit=20&offset=0&sort=created_at_desc
```

**Response:**
```json
{
  "data": [ ... ],
  "pagination": {
    "total": 150,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

---

## WebSocket (Future Feature)

**Endpoint:** `ws://localhost:8080/ws`

リアルタイム分析進捗通知用

```json
// Client -> Server
{
  "action": "subscribe",
  "job_id": "job_abc123"
}

// Server -> Client (Progress Update)
{
  "type": "progress",
  "job_id": "job_abc123",
  "progress": 45,
  "message": "Searching academic databases..."
}

// Server -> Client (Completion)
{
  "type": "completed",
  "job_id": "job_abc123",
  "result": { /* analysis result */ }
}
```

---

## OpenAPI Specification

完全なOpenAPI 3.0仕様は `/swagger.json` で取得可能

Swagger UI: `http://localhost:8080/swagger/`

---

## Examples

### cURL Examples

```bash
# Basic Analysis
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Vaccines cause autism.",
    "options": {
      "llm_provider": "claude",
      "min_credibility": 0.8
    }
  }'

# Async Analysis
curl -X POST http://localhost:8080/api/v1/analyze/async \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Long research paper...",
    "options": {
      "max_results": 50
    }
  }'

# Check Job Status
curl http://localhost:8080/api/v1/jobs/job_abc123
```

### Go Client Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type AnalyzeRequest struct {
    Text    string         `json:"text"`
    Options AnalyzeOptions `json:"options"`
}

type AnalyzeOptions struct {
    LLMProvider    string  `json:"llm_provider"`
    MinCredibility float64 `json:"min_credibility"`
}

func analyzeText(text string) error {
    req := AnalyzeRequest{
        Text: text,
        Options: AnalyzeOptions{
            LLMProvider:    "claude",
            MinCredibility: 0.7,
        },
    }
    
    body, _ := json.Marshal(req)
    resp, err := http.Post(
        "http://localhost:8080/api/v1/analyze",
        "application/json",
        bytes.NewBuffer(body),
    )
    // ... handle response
    return nil
}
```

---

## Versioning

API versioning via URL path: `/api/v1/`, `/api/v2/`

**Deprecation Policy:**
- 6ヶ月前に告知
- 古いバージョンは最低1年サポート

---

## Error Handling Best Practices

1. **すべてのエラーレスポンスに `request_id` を含める**
2. **HTTP Status Codeを適切に使用**:
   - `200 OK` - 成功
   - `400 Bad Request` - バリデーションエラー
   - `429 Too Many Requests` - レート制限
   - `500 Internal Server Error` - サーバーエラー
   - `503 Service Unavailable` - LLMプロバイダーダウン
3. **詳細なエラーメッセージを提供**（本番環境では機密情報を隠す）

---

## Future Enhancements

- [ ] GraphQL endpoint
- [ ] Batch analysis
- [ ] File upload support (PDF, DOCX)
- [ ] Collaborative annotations
- [ ] Custom credibility weights
- [ ] Multi-language UI
