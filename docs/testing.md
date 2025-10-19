# SourceTracer Testing Strategy

> **プロジェクト使命**: 情報源を徹底的に追跡し、意見と事実を峻別する。
> **警告**: このシステムの品質は人類の情報信頼性に直結します。テストの妥協は許されません。

## Testing Philosophy

### なぜテストが絶対に必要か

1. **社会的影響**: 誤った判定は人々の意思決定に直接影響
2. **複雑性**: LLM統合、信頼性スコアリングは手動検証不可能
3. **長期保守**: テストなしでは後からバグ修正が困難
4. **信頼の証明**: テストケースが仕様書となり、品質を保証

### TDD厳守（t-wada流）

**絶対ルール:**
- 🔴 Red → 🟢 Green → 🔵 Refactor のサイクルを守る
- 実装前にテストを書く（テスト後書き禁止）
- モックでビジネスロジックを隠蔽しない
- カバレッジ詐欺（意味のないテスト）禁止
- "動いたからOK"は通用しない

---

## Testing Pyramid

```
           /\
          /  \         E2E Tests (10%)
         /____\
        /      \       Integration Tests (30%)
       /________\
      /          \     Unit Tests (60%)
     /____________\
```

**目標カバレッジ:**
- Unit Tests: 80%+
- Integration Tests: 70%+
- E2E Tests: Critical paths only

---

## 1. Unit Tests

### 1.1 LLM Provider Tests

**Location:** `backend/internal/llm/*_test.go`

#### OpenAI Provider
```go
func TestOpenAIProvider_Analyze(t *testing.T) {
    // Mock HTTP client
    // Test successful response
    // Test error handling
    // Test timeout
    // Test retry logic
}

func TestOpenAIProvider_ClassifyClaim(t *testing.T) {
    // Test opinion classification
    // Test fact classification
    // Test mixed classification
    // Test edge cases (empty text, very long text)
}
```

#### Claude Provider
```go
func TestClaudeProvider_GenerateReport(t *testing.T) {
    // Test report structure
    // Test markdown formatting
    // Test evidence integration
}
```

#### LLM Router
```go
func TestLLMRouter_SelectProvider(t *testing.T) {
    // Test provider selection logic
    // Test fallback mechanism
    // Test all providers unavailable
}
```

**Test Cases:**
- ✅ 正常なレスポンス処理
- ✅ APIエラーハンドリング（401, 429, 500）
- ✅ タイムアウト処理
- ✅ リトライロジック（Exponential backoff）
- ✅ レスポンスパース（JSON, streaming）
- ✅ トークン制限チェック
- ✅ キャッシュヒット/ミス

---

### 1.2 Text Analyzer Tests

**Location:** `backend/internal/analyzer/*_test.go`

```go
func TestClaimExtractor_Extract(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []Claim
    }{
        {
            name:  "single fact claim",
            input: "The Earth is round.",
            expected: []Claim{
                {Type: Fact, Text: "The Earth is round.", Confidence: 0.95},
            },
        },
        {
            name:  "opinion claim",
            input: "Python is the best language.",
            expected: []Claim{
                {Type: Opinion, Text: "Python is the best language.", Confidence: 0.90},
            },
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}

func TestClaimClassifier_Classify(t *testing.T) {
    // Test fact identification
    // Test opinion identification
    // Test mixed claims
    // Test unclear claims
}

func TestCredibilityScorer_Score(t *testing.T) {
    // Test academic source scoring
    // Test news source scoring
    // Test blog source scoring
    // Test weight calculation
}
```

**Test Cases:**
- ✅ 短文・長文の処理
- ✅ 複数言語対応（日本語、英語、中国語）
- ✅ 特殊文字・絵文字
- ✅ 空文字列・null処理
- ✅ Position tracking（開始・終了位置）
- ✅ Context extraction

---

### 1.3 Search Engine Tests

**Location:** `backend/internal/search/*_test.go`

```go
func TestSemanticScholarClient_Search(t *testing.T) {
    // Mock API responses
    // Test result parsing
    // Test pagination
    // Test rate limiting
}

func TestGoogleSearchClient_Search(t *testing.T) {
    // Test custom search API
    // Test result filtering
}

func TestDuckDuckGoScraper_Search(t *testing.T) {
    // Mock Playwright browser
    // Test HTML parsing
    // Test result extraction
}
```

**Test Cases:**
- ✅ API レスポンスパース
- ✅ ページネーション処理
- ✅ レート制限検出・遅延
- ✅ タイムアウト処理
- ✅ 結果フィルタリング
- ✅ 重複除去
- ✅ エラーハンドリング（404, 503）

---

### 1.4 Database Layer Tests

**Location:** `backend/internal/db/*_test.go`

```go
func TestAnalysisRepository_Create(t *testing.T) {
    // Use testcontainers for real PostgreSQL
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewAnalysisRepository(db)
    analysis := &Analysis{...}
    
    err := repo.Create(context.Background(), analysis)
    assert.NoError(t, err)
    assert.NotEmpty(t, analysis.ID)
}

func TestEvidenceRepository_FindByClaimID(t *testing.T) {
    // Test query with joins
    // Test ordering
    // Test pagination
}

func TestSearchCache_Get(t *testing.T) {
    // Test cache hit
    // Test cache miss
    // Test expiration
}
```

**Test Cases:**
- ✅ CRUD操作（Create, Read, Update, Delete）
- ✅ トランザクション処理
- ✅ 制約違反（UNIQUE, FK）
- ✅ NULL値処理
- ✅ JSONB フィールド操作
- ✅ ベクトル検索（pgvector）
- ✅ ページネーション
- ✅ 並行アクセス（Race condition）

---

### 1.5 Scraper Tests

**Location:** `backend/internal/scraper/*_test.go`

```go
func TestBrowserPool_Acquire(t *testing.T) {
    // Test pool initialization
    // Test concurrent access
    // Test max pool size
    // Test browser reuse
}

func TestPlaywrightScraper_Scrape(t *testing.T) {
    // Mock Playwright
    // Test JavaScript execution
    // Test element waiting
    // Test screenshot capture
}
```

**Test Cases:**
- ✅ ブラウザプール管理
- ✅ ページロード待機
- ✅ JavaScript実行
- ✅ スクリーンショット取得
- ✅ タイムアウト処理
- ✅ エラーページ検出
- ✅ リソースブロック（画像、CSS）

---

## 2. Integration Tests

**Tag:** `// +build integration`

### 2.1 API Endpoint Tests

**Location:** `backend/test/integration/api_test.go`

```go
// +build integration

func TestAnalyzeEndpoint_Integration(t *testing.T) {
    // Start test server
    server := setupTestServer(t)
    defer server.Close()
    
    // Make request
    resp := makeRequest(t, server.URL+"/api/v1/analyze", map[string]interface{}{
        "text": "Climate change is real.",
        "options": map[string]interface{}{
            "llm_provider": "mock",
            "min_credibility": 0.7,
        },
    })
    
    // Verify response
    assert.Equal(t, 200, resp.StatusCode)
    
    var result AnalyzeResponse
    json.Unmarshal(resp.Body, &result)
    assert.True(t, result.Success)
    assert.NotEmpty(t, result.Data.Claims)
}

func TestAsyncAnalyze_Integration(t *testing.T) {
    // Test job creation
    // Poll job status
    // Verify completion
}
```

**Test Cases:**
- ✅ POST /api/v1/analyze（全フロー）
- ✅ POST /api/v1/analyze/async + GET /api/v1/jobs/{id}
- ✅ POST /api/v1/classify
- ✅ POST /api/v1/evidence/search
- ✅ GET /api/v1/sources
- ✅ レート制限動作確認
- ✅ エラーレスポンス形式
- ✅ CORS設定

---

### 2.2 Database Integration Tests

**Location:** `backend/test/integration/db_test.go`

```go
func TestFullAnalysisWorkflow_Integration(t *testing.T) {
    db := setupTestDB(t)
    
    // 1. Create analysis
    analysis := createAnalysis(t, db)
    
    // 2. Create claims
    claims := createClaims(t, db, analysis.ID)
    
    // 3. Create evidences
    evidences := createEvidences(t, db, claims[0].ID)
    
    // 4. Verify relationships
    retrieved := getAnalysisWithClaims(t, db, analysis.ID)
    assert.Len(t, retrieved.Claims, len(claims))
}

func TestVectorSearch_Integration(t *testing.T) {
    // Insert embeddings
    // Perform similarity search
    // Verify results ordered by distance
}
```

**Test Cases:**
- ✅ マルチテーブルトランザクション
- ✅ カスケード削除
- ✅ ベクトル類似度検索
- ✅ 全文検索（pg_trgm）
- ✅ キャッシュ有効性
- ✅ 並行書き込み

---

### 2.3 LLM Integration Tests

**Location:** `backend/test/integration/llm_test.go`

```go
func TestRealLLMProviders_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping real LLM test in short mode")
    }
    
    providers := []string{"openai", "claude", "deepseek"}
    
    for _, provider := range providers {
        t.Run(provider, func(t *testing.T) {
            client := createLLMClient(t, provider)
            result, err := client.Analyze(context.Background(), "Test text")
            assert.NoError(t, err)
            assert.NotEmpty(t, result)
        })
    }
}
```

**Test Cases:**
- ✅ 実際のAPI呼び出し（環境変数で制御）
- ✅ レスポンス形式検証
- ✅ エラーハンドリング
- ✅ キャッシュ機能

---

### 2.4 Playwright Integration Tests

**Location:** `backend/test/integration/scraper_test.go`

```go
func TestRealBrowserScraping_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping browser test in short mode")
    }
    
    pool := setupBrowserPool(t)
    defer pool.Close()
    
    scraper := NewPlaywrightScraper(pool)
    result, err := scraper.Scrape("https://example.com")
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result.HTML)
    assert.Contains(t, result.HTML, "<html")
}
```

**Test Cases:**
- ✅ 実際のWebページスクレイピング
- ✅ JavaScript heavy サイト
- ✅ 検索エンジンページ解析
- ✅ ブラウザプールの並行処理

---

## 3. E2E Tests

**Tool:** Playwright (TypeScript)

**Location:** `backend/test/e2e/`

### 3.1 Full Analysis Flow

```typescript
// test/e2e/analyze.spec.ts
import { test, expect } from '@playwright/test';

test('complete analysis workflow', async ({ page }) => {
  // 1. Navigate to app
  await page.goto('http://localhost:3000');
  
  // 2. Enter text
  await page.fill('#input-text', 'Climate change is accelerating.');
  
  // 3. Submit
  await page.click('#analyze-button');
  
  // 4. Wait for results
  await page.waitForSelector('.analysis-results');
  
  // 5. Verify claims displayed
  const claims = await page.locator('.claim-item');
  expect(await claims.count()).toBeGreaterThan(0);
  
  // 6. Verify evidences
  const evidences = await page.locator('.evidence-item');
  expect(await evidences.count()).toBeGreaterThan(0);
  
  // 7. Check credibility scores
  const credibility = await page.textContent('.overall-credibility');
  expect(parseFloat(credibility)).toBeGreaterThan(0);
});
```

### 3.2 CLI E2E Tests

```bash
#!/bin/bash
# test/e2e/cli_test.sh

# Build CLI
go build -o bin/factcheck cmd/cli/main.go

# Test basic analysis
output=$(echo "Python is the best language" | ./bin/factcheck check --stdin --llm mock)
echo "$output" | grep -q "opinion"
echo "✓ CLI basic test passed"

# Test file analysis
./bin/factcheck analyze --file test/fixtures/sample.txt --output json > /tmp/result.json
jq -e '.success == true' /tmp/result.json
echo "✓ CLI file test passed"

# Test error handling
./bin/factcheck analyze --file nonexistent.txt 2>&1 | grep -q "file not found"
echo "✓ CLI error handling passed"
```

**Test Cases:**
- ✅ ユーザー入力 → 分析 → 結果表示
- ✅ 非同期ジョブの進捗表示
- ✅ エラーメッセージ表示
- ✅ エクスポート機能（JSON, Markdown, PDF）
- ✅ フィルタリング（信頼性、ソースタイプ）
- ✅ ページネーション
- ✅ レスポンシブデザイン

---

## 4. Performance Tests

**Tool:** `k6` (Load Testing)

**Location:** `test/performance/`

### 4.1 API Load Test

```javascript
// test/performance/api_load.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% under 2s
    http_req_failed: ['rate<0.01'],    // Error rate < 1%
  },
};

export default function () {
  const payload = JSON.stringify({
    text: 'Climate change is real.',
    options: { llm_provider: 'mock' },
  });
  
  const res = http.post('http://localhost:8080/api/v1/analyze', payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 2s': (r) => r.timings.duration < 2000,
  });
  
  sleep(1);
}
```

**Metrics:**
- ✅ スループット（req/sec）
- ✅ レスポンスタイム（p50, p95, p99）
- ✅ エラーレート
- ✅ データベース接続プール使用率
- ✅ メモリ使用量
- ✅ CPU使用率

---

## 5. Security Tests

### 5.1 Input Validation

```go
func TestInputValidation_SQLInjection(t *testing.T) {
    maliciousInputs := []string{
        "'; DROP TABLE users; --",
        "1' OR '1'='1",
        "<script>alert('xss')</script>",
    }
    
    for _, input := range maliciousInputs {
        _, err := analyzeText(input)
        // Should not cause errors or unexpected behavior
        assert.NoError(t, err)
    }
}

func TestInputValidation_LargePayload(t *testing.T) {
    // Test 10MB text
    largeText := strings.Repeat("a", 10*1024*1024)
    _, err := analyzeText(largeText)
    assert.Error(t, err) // Should reject
}
```

**Test Cases:**
- ✅ SQLインジェクション防止
- ✅ XSS防止
- ✅ Path traversal防止
- ✅ 過大リクエスト拒否
- ✅ レート制限機能
- ✅ 環境変数漏洩チェック

---

## 6. Test Data & Fixtures

### 6.1 Sample Texts

**Location:** `test/fixtures/`

```
test/fixtures/
├── opinions.txt          # 意見の例
├── facts.txt             # 事実の例
├── mixed.txt             # 混合の例
├── long_article.txt      # 長文テスト用
├── multilingual.txt      # 多言語テスト
└── edge_cases.txt        # エッジケース
```

### 6.2 Mock Data

```go
// test/mocks/llm_mock.go
type MockLLMProvider struct{}

func (m *MockLLMProvider) Analyze(ctx context.Context, text string) (*Analysis, error) {
    return &Analysis{
        Claims: []Claim{
            {Type: Opinion, Text: text, Confidence: 0.9},
        },
    }, nil
}

// test/mocks/search_mock.go
type MockSearchEngine struct{}

func (m *MockSearchEngine) Search(query string) ([]Evidence, error) {
    return []Evidence{
        {Title: "Mock Evidence", URL: "https://example.com", Credibility: 0.8},
    }, nil
}
```

---

## 7. CI/CD Testing Pipeline

**GitHub Actions:** `.github/workflows/test.yml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: pgvector/pgvector:pg16
        env:
          POSTGRES_PASSWORD: test
      mongodb:
        image: mongo:7
      redis:
        image: redis:7
    steps:
      - uses: actions/checkout@v3
      - name: Run integration tests
        run: go test -v -tags=integration ./...

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Start services
        run: docker-compose up -d
      - name: Run E2E tests
        run: ./test/e2e/run_all.sh

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
```

---

## 8. Test Commands

### ローカル実行

```bash
# All unit tests
go test ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test ./internal/llm/...

# Verbose mode
go test -v ./...

# Race detection
go test -race ./...

# Integration tests
go test -tags=integration ./...

# Short mode (skip slow tests)
go test -short ./...

# Run specific test
go test -run TestAnalyzer_Extract ./internal/analyzer/

# Benchmark
go test -bench=. ./...
```

### CI/CD実行

```bash
# Pre-commit checks
make test-all

# Coverage report
make coverage

# Lint
make lint

# E2E
make e2e-test
```

---

## 9. Makefile Targets

```makefile
# Makefile
.PHONY: test test-unit test-integration test-e2e coverage lint

test: test-unit test-integration

test-unit:
	go test -v -race -cover ./...

test-integration:
	docker-compose -f docker-compose.test.yml up -d
	go test -v -tags=integration ./...
	docker-compose -f docker-compose.test.yml down

test-e2e:
	./test/e2e/run_all.sh

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run
	go vet ./...

test-all: lint test test-e2e
```

---

## 10. Testing Best Practices

### DO's ✅
- **Table-Driven Tests** を使用
- テスト名は説明的に（`TestAnalyzer_ExtractClaims_WithEmptyText`）
- `testify/assert` でアサーション
- モックは `interface` で実装
- CI/CDでテスト自動化
- カバレッジ80%以上を目標
- テストデータをコミット（`test/fixtures/`）

### DON'Ts ❌
- 本番APIキーをテストで使用しない
- テスト間で状態共有しない
- `time.Sleep()` に依存しない（flaky test）
- ハードコードされたポート番号（環境変数使用）
- テストデータをDBに永続化しない

---

## 11. Coverage Goals

| Component | Target Coverage |
|-----------|----------------|
| LLM Providers | 85%+ |
| Text Analyzer | 90%+ |
| Search Engines | 80%+ |
| Database Layer | 85%+ |
| API Handlers | 75%+ |
| Scraper | 70%+ |
| **Overall** | **80%+** |

---

## 12. Test Monitoring

### Metrics to Track
- ✅ Test execution time
- ✅ Flaky test rate
- ✅ Code coverage trend
- ✅ Failed test frequency
- ✅ CI/CD pipeline duration

### Tools
- **CodeCov**: カバレッジ追跡
- **SonarQube**: コード品質
- **GitHub Actions**: CI/CD
- **k6**: パフォーマンステスト

---

## 13. Future Test Enhancements

- [ ] Mutation Testing (go-mutesting)
- [ ] Chaos Engineering (テスト環境で障害注入)
- [ ] A/B Testing フレームワーク
- [ ] Visual Regression Tests (Percy, Chromatic)
- [ ] Accessibility Tests (axe-core)
- [ ] Contract Testing (Pact)
