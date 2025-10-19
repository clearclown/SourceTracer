# SourceTracer Backend - Quickstart Guide

## 🎯 Current Status: v0.1.0-alpha

実装済み機能：
- ✅ ドメインモデル (Claim, Evidence)
- ✅ ルールベース分類器 (Opinion/Fact判定)
- ✅ 設定管理 (環境変数読み込み)
- ✅ LLMプロバイダーインターフェース
- ✅ CLIツール

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
🔍 SourceTracer CLI v0.1.0-alpha
📝 Environment: development
🤖 Default LLM: claude

=============================================================
📄 Text: Climate change is accelerating faster than predicted.
-------------------------------------------------------------
📊 Type: fact
📈 Confidence: 0.70
💡 Reasoning: No opinion keywords found
=============================================================

✅ This appears to be a FACT.
   Verify with reliable sources for accuracy.
```

#### カスタムテキストで実行
```bash
./bin/sourcetracer -text "Python is the best programming language"
```

出力例:
```
=============================================================
📄 Text: Python is the best programming language
-------------------------------------------------------------
📊 Type: opinion
📈 Confidence: 0.90
💡 Reasoning: Contains opinion keyword: best
=============================================================

✅ This appears to be an OPINION.
   Consider providing evidence or rephrasing as fact.
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
- `llm`: 85.7%
- `config`: 54.5%

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

## 🔜 次のステップ (v0.2)

- [ ] LLMベース分類器実装 (OpenAI/Claude)
- [ ] Evidence検索機能 (Semantic Scholar統合)
- [ ] APIサーバー (Gin)
- [ ] PostgreSQLスキーマ
- [ ] E2Eテスト

## 📖 参考

- [README.md](../README.md) - プロジェクト全体概要
- [CLAUDE.md](../CLAUDE.md) - 開発ガイドライン
- [docs/api-design.md](../docs/api-design.md) - API設計
- [docs/testing.md](../docs/testing.md) - テスト戦略

---

**SourceTracer v0.1.0-alpha** - 情報源を、徹底的に追う。
