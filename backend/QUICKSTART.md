# SourceTracer Backend - Quickstart Guide

## 🎯 Current Status: v0.2.0-alpha

実装済み機能：
- ✅ ドメインモデル (Claim, Evidence)
- ✅ ルールベース分類器 (Opinion/Fact判定)
- ✅ **Claude API統合** (LLMベース分類)
- ✅ **Semantic Scholar統合** (学術論文検索)
- ✅ **Analyzer** (複数クレーム抽出・分類)
- ✅ 設定管理 (環境変数読み込み)
- ✅ LLMプロバイダーインターフェース
- ✅ CLIツール (v0.2)

## 🚀 クイックスタート

### 1. ビルド

```bash
cd backend
go build -o bin/sourcetracer cmd/cli/main.go
```

### 2. 実行例

#### デフォルトサンプルテキストで実行
```bash
./bin/sourcetracer
```

出力例:
```
🔍 SourceTracer CLI v0.2.0-alpha
📝 Environment: development
🤖 Default LLM: claude

=============================================================
📄 Original Text: Climate change is accelerating faster than predicted.
-------------------------------------------------------------
📊 Analysis ID: anl_1760860015098070181
⏱️  Processing Time: 0ms
📈 Overall Credibility: 0.70

📋 Claims Found: 1

Claim #1:
  Text: Climate change is accelerating faster than predicted.
  Type: fact
  Confidence: 0.70
  ✅ FACT - Verify with reliable sources

=============================================================
```

#### カスタムテキストで実行
```bash
./bin/sourcetracer -text "Python is the best language. Go is faster. Rust is safer."
```

出力例:
```
📋 Claims Found: 3

Claim #1:
  Text: Python is the best language.
  Type: opinion
  Confidence: 0.90
  ✅ OPINION - Consider providing evidence

Claim #2:
  Text: Go is faster.
  Type: fact
  Confidence: 0.70
  ✅ FACT - Verify with reliable sources

Claim #3:
  Text: Rust is safer.
  Type: fact
  Confidence: 0.70
  ✅ FACT - Verify with reliable sources
```

### 3. テスト実行

```bash
# すべてのテスト実行
go test ./...

# カバレッジ付き
go test ./... -cover

# 詳細出力
go test ./... -v
```

現在のカバレッジ:
- `classifier`: 100.0% ✅
- `domain`: 87.5%
- `analyzer`: 66.1%
- `search`: 66.7%
- `config`: 54.5%
- `llm`: 34.0% (統合テストは別途)

## 📝 現在の分類アルゴリズム

v0.1では**ルールベース分類器**を使用：

### Opinion キーワード
- best, worst
- should, must
- better, worse
- excellent, terrible
- amazing, awful
- prefer

### ロジック
1. テキストにOpinionキーワードが含まれる → **Opinion** (confidence: 0.9)
2. それ以外 → **Fact** (confidence: 0.7)
3. 空文字列 → **Unclear** (confidence: 0.0)

**注意**: これは暫定実装です。将来的にLLMベースの高精度分類に置き換えます。

## 🧪 テスト駆動開発 (TDD)

このプロジェクトは**t-wada流TDD**を厳守しています：

### テスト例: Classifier

```go
// 🔴 Red: 失敗するテストを書く
func TestClassifier_ClassifyOpinion(t *testing.T) {
    classifier := NewClassifier()
    result := classifier.Classify("Python is the best")

    assert.Equal(t, Opinion, result.Type)
    assert.Greater(t, result.Confidence, 0.8)
}

// 🟢 Green: 最小限の実装
func (c *Classifier) Classify(text string) *Result {
    if strings.Contains(text, "best") {
        return &Result{Type: Opinion, Confidence: 0.9}
    }
    return &Result{Type: Fact, Confidence: 0.7}
}

// 🔵 Refactor: リファクタリング
const OpinionKeywords = []string{"best", "worst", "should"}
```

## 🛠️ 開発コマンド

```bash
# モジュール整理
go mod tidy

# フォーマット
gofmt -w .

# Linter (要インストール)
golangci-lint run

# 特定パッケージのテスト
go test ./internal/classifier/ -v

# ベンチマーク
go test -bench=. ./internal/classifier/
```

## 📂 ディレクトリ構造

```
backend/
├── cmd/
│   └── cli/              # CLIツール
├── internal/
│   ├── classifier/       # 分類器 (100%カバレッジ)
│   ├── config/           # 設定管理
│   ├── domain/           # ドメインモデル
│   └── llm/              # LLMプロバイダー
└── bin/                  # ビルド成果物
```

## ✅ v0.2 新機能

### 1. Analyzer (複数クレーム抽出)
- 文単位でのクレーム抽出
- 各クレームの個別分類
- 全体的な信頼度スコア計算

### 2. Claude API統合
- プロンプトエンジニアリング
- JSON応答パース
- エラーハンドリング

### 3. Semantic Scholar統合
- 学術論文検索API
- 引用数ベースの信頼度計算
- パース・構造化

## 🔜 次のステップ (v0.3)

- [ ] APIサーバー (Gin framework)
- [ ] PostgreSQL統合
- [ ] リアルタイムLLM分類
- [ ] 複数ソースからのエビデンス集約
- [ ] E2E統合テスト

## 📖 参考

- [README.md](../README.md) - プロジェクト全体概要
- [CLAUDE.md](../CLAUDE.md) - 開発ガイドライン
- [docs/api-design.md](../docs/api-design.md) - API設計
- [docs/testing.md](../docs/testing.md) - テスト戦略

---

**SourceTracer v0.2.0-alpha** - 情報源を、徹底的に追う。すべての主張にエビデンスを。
