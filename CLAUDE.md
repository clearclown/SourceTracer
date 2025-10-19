# SourceTracer Project - Claude Memory

> **プロジェクト使命**: 情報源を徹底的に追跡し、意見と事実を峻別する。

> **警告**: このシステムが中途半端な実装になると、盗作・盗説が横行し、人類に大きな損害をもたらす可能性があります。
> 品質への妥協は許されません。TDD厳守、モック・スクリプトでのごまかし禁止。

プロジェクト概要については @README.md を、利用可能なAPIエンドポイントについては @docs/api-design.md を参照してください。

---

## プロジェクトの重要性と責任

### SourceTracerのミッション

**"Trace every claim to its source"** - すべての主張を情報源まで追跡する

1. **意見と事実の峻別**: 人間が曖昧にしがちな境界を明確化
2. **批判的思考の支援**: 主張の弱点を多角的に指摘
3. **徹底的なソース追跡**: 複数の検索エンジン・データベースを横断
4. **学術的誠実性の保護**: 盗作・盗説を防止し、適切な引用を促進

### なぜこのプロジェクトが重要か

1. **情報の信頼性危機**: デジタル時代において、意見と事実の境界が曖昧になり、誤情報が拡散
2. **学術的誠実性**: 不正確な引用や盗説は学問の信頼性を根本から損なう
3. **検索エンジンの限界**: Google等の単一検索では、アカデミックな情報が見つからないことが多い
4. **社会的影響**: このツールが不完全な場合、以下のリスクが生じる：
   - **誤った自信**: 低品質なエビデンスを「検証済み」と誤認させる
   - **盗作の助長**: 適切な引用なしにコンテンツを流用する行為を正当化
   - **情報操作**: 偏ったソース選択により特定の主張を不当に強化
   - **信頼の喪失**: 誤った判定により、ツール全体への信頼が失われる

### 開発者の責任

我々は以下を厳守する：

- **品質第一**: 速度よりも正確性を優先
- **透明性**: アルゴリズムの動作を説明可能にする
- **検証可能性**: すべての判定根拠を明示
- **倫理的配慮**: バイアスを最小化し、公正な評価を行う

---

## TDD厳守ルール（t-wada流）

### 🚨 CRITICAL: TDDは妥協不可

**このプロジェクトでTDDを厳守する理由:**

1. **人命・社会に影響**: 誤った情報判定は人々の意思決定に直接影響
2. **複雑なロジック**: LLM統合、信頼性スコアリングは手動テスト不可能
3. **長期保守性**: TDDなしでは後からバグ修正が困難
4. **信頼性証明**: テストケースが仕様書となり、品質を保証

### ❌ 絶対禁止事項

- **モックでのごまかし**: 本質的なロジックをモックで隠蔽しない
- **テスト後書き**: 実装後にテストを追加する行為
- **カバレッジ詐欺**: 意味のないテストでカバレッジを稼ぐ
- **統合テストのみ**: ユニットテストを省略して統合テストだけで済ませる
- **手動確認**: "動いたからOK"で次に進む
- **TODO放置**: テストリストの項目を実装せずに残す

### ✅ TDD TODOリスト管理

#### 基本方針

- 🔴 **Red**: 失敗するテストを書く（コンパイルエラーもOK）
- 🟢 **Green**: テストを通す最小限の実装
- 🔵 **Refactor**: リファクタリング
- **小さなステップ**: 一度に1つの機能のみ実装
- **仮実装**: ベタ書き（例: `return 42`）から始める
- **三角測量**: 2つ目、3つ目のテストケースで一般化
- **明白な実装**: 自信がある場合は直接実装してもOK
- **不安駆動**: 不安なところから先にテストを書く

#### TDD実践の具体的ステップ

```
1. TODOリストに項目を追加
   例: "[ ] Claim型をOpinion/Fact/Mixedに分類する"

2. 🔴 失敗するテストを書く
   func TestClassifier_Opinion(t *testing.T) {
       result := Classify("Python is the best")
       assert.Equal(t, Opinion, result.Type)
   }
   
   → コンパイルエラー: Classify関数が存在しない

3. 🟢 最小限の実装（仮実装）
   func Classify(text string) ClaimType {
       return Opinion  // ベタ書き
   }
   
   → テスト通過

4. 🔴 次のテストを追加（三角測量）
   func TestClassifier_Fact(t *testing.T) {
       result := Classify("Earth is round")
       assert.Equal(t, Fact, result.Type)
   }
   
   → テスト失敗: Opinionしか返さない

5. 🟢 一般化
   func Classify(text string) ClaimType {
       if contains(text, "best", "worst", "should") {
           return Opinion
       }
       return Fact
   }
   
   → テスト通過

6. 🔵 リファクタリング
   - マジックワードをconstに
   - ロジックを関数分割
   - 変数名を明確化

7. ✅ TODOリストを更新
   [x] Claim型をOpinion/Fact/Mixedに分類する
   [ ] Confidenceスコアを計算する ← 次のタスク
```

#### TODOリストのフォーマット

```markdown
### 実装TODOリスト

#### Phase 1: コア機能
- [ ] Claim抽出
  - [ ] 文を分割する
  - [ ] 各文をClaimとして認識
  - [ ] Position情報を記録
- [ ] Claim分類
  - [x] Opinion判定
  - [x] Fact判定
  - [ ] Mixed判定
  - [ ] Unclear判定
  - [ ] Confidenceスコア計算

#### Phase 2: LLM統合
- [ ] プロバイダーインターフェース定義
- [ ] OpenAI実装
- [ ] Claude実装
- [ ] DeepSeek実装
- [ ] フォールバック機能

#### 実装中に気づいたTODO
- [ ] エラーメッセージの国際化
- [ ] ログレベルの設定可能化
```

#### コミットルール

```bash
# 🔴 テストを書いたら
git commit -m "test: add failing test for claim classification"

# 🟢 テストを通したら
git commit -m "feat: implement claim classification to pass test"

# 🔵 リファクタリングしたら
git commit -m "refactor: extract opinion keywords to constant"

# 複数ステップを一度に進めない！
# 各ステップごとにコミットする
```

### TDD実践例: Claim分類器

```go
// ステップ1: 🔴 失敗するテストを書く
func TestClassifier_ClassifyOpinion(t *testing.T) {
    classifier := NewClassifier()
    result := classifier.Classify("Python is the best language")
    
    assert.Equal(t, Opinion, result.Type)
    assert.Greater(t, result.Confidence, 0.8)
}
// → コンパイルエラー: Classifierが存在しない

// ステップ2: 🟢 最小限の実装
type Classifier struct{}

func NewClassifier() *Classifier {
    return &Classifier{}
}

func (c *Classifier) Classify(text string) ClaimType {
    return ClaimType{Type: Opinion, Confidence: 0.9} // 仮実装
}
// → テスト通過

// ステップ3: 🔴 次のテストケース
func TestClassifier_ClassifyFact(t *testing.T) {
    classifier := NewClassifier()
    result := classifier.Classify("The Earth orbits the Sun")
    
    assert.Equal(t, Fact, result.Type)
}
// → テスト失敗

// ステップ4: 🟢 三角測量で一般化
func (c *Classifier) Classify(text string) ClaimType {
    opinionKeywords := []string{"best", "worst", "should", "must"}
    
    for _, keyword := range opinionKeywords {
        if strings.Contains(strings.ToLower(text), keyword) {
            return ClaimType{Type: Opinion, Confidence: 0.9}
        }
    }
    
    return ClaimType{Type: Fact, Confidence: 0.7}
}
// → すべてのテスト通過

// ステップ5: 🔵 リファクタリング
const (
    OpinionKeywords = []string{"best", "worst", "should", "must", "better"}
)

func (c *Classifier) Classify(text string) ClaimType {
    normalized := strings.ToLower(text)
    
    for _, keyword := range OpinionKeywords {
        if strings.Contains(normalized, keyword) {
            confidence := c.calculateConfidence(text, keyword)
            return ClaimType{Type: Opinion, Confidence: confidence}
        }
    }
    
    return ClaimType{Type: Fact, Confidence: 0.7}
}

func (c *Classifier) calculateConfidence(text, keyword string) float64 {
    // より洗練された計算（後でテストを追加）
    return 0.9
}
```

---

## モック使用の厳格なガイドライン

### モックを使って良い場合

✅ **外部依存のみ**:
- HTTPクライアント（API呼び出し）
- データベース接続
- ファイルシステム
- 時刻取得（`time.Now()`）
- ランダム生成（`rand.Int()`）

```go
// OK: 外部API呼び出しのモック
type MockLLMClient struct {
    Response string
    Error    error
}

func (m *MockLLMClient) Call(prompt string) (string, error) {
    return m.Response, m.Error
}
```

### モックを使ってはいけない場合

❌ **ビジネスロジック**:
- 分類アルゴリズム
- スコア計算
- テキスト解析
- バリデーション

```go
// NG: ロジックをモックで隠蔽
type MockClassifier struct {
    ReturnType ClaimType  // これはダメ
}

// OK: 実装をテスト
func TestRealClassifier(t *testing.T) {
    classifier := NewClassifier()  // 実物を使う
    result := classifier.Classify("test")
    // 実際のロジックをテスト
}
```

### 統合テストとの使い分け

- **ユニットテスト**: ビジネスロジックは実装、外部依存はモック
- **統合テスト**: 可能な限り実物を使用（testcontainersなど）

---

## プロジェクト構成

- **言語**: Go 1.21+ (backend), Dart/Flutter (frontend)
- **アーキテクチャ**: マイクロサービス風モノリス
- **デプロイ**: Podman/Docker Compose

---

## コーディング規約

### Go

#### スタイル
- `gofmt` と `goimports` を常に実行
- `golangci-lint` でリント（設定: `.golangci.yml`）
- 2スペースインデント（Goデフォルト）
- 行の長さ: 最大120文字

#### 命名規則
- **パッケージ名**: 小文字、アンダースコアなし（`llmprovider` ではなく `llm`）
- **インターフェース**: 動詞 + er（`Analyzer`, `Searcher`, `Provider`）
- **構造体**: PascalCase、説明的な名前
- **変数**: camelCase、省略形は避ける（`cfg` ではなく `config`）
- **定数**: UPPER_SNAKE_CASE または PascalCase（コンテキスト依存）

#### エラーハンドリング
```go
// Good: ラップして文脈を追加
if err != nil {
    return fmt.Errorf("failed to analyze text: %w", err)
}

// Bad: エラーを握りつぶす
if err != nil {
    log.Println(err)
}
```

#### 構造体の初期化
```go
// Good: フィールド名を明示
user := &User{
    ID:   123,
    Name: "Alice",
}

// Bad: 位置依存
user := &User{123, "Alice"}
```

### Flutter/Dart

- `dart format` を常に実行
- `flutter analyze` でリント
- 2スペースインデント
- 状態管理: Riverpod使用
- ファイル名: snake_case（`fact_check_page.dart`）

---

## ディレクトリ構造

```
backend/
├── cmd/              # エントリポイント（main.go）
│   ├── api/         # APIサーバー
│   ├── worker/      # バックグラウンドワーカー
│   └── cli/         # CLIツール
├── internal/        # プライベートコード
│   ├── config/      # 環境変数・設定
│   ├── llm/         # LLMプロバイダー実装
│   ├── search/      # 検索エンジンクライアント
│   ├── scraper/     # Playwrightスクレイパー
│   ├── analyzer/    # テキスト分析ロジック
│   ├── db/          # データベースアクセス
│   └── api/         # HTTPハンドラー
└── pkg/             # 公開可能なライブラリ
```

---

## 頻繁に使用するコマンド

### 開発環境

```bash
# ローカル起動（全サービス）
podman-compose up -d

# ログ確認
podman-compose logs -f api

# データベースリセット
podman-compose down -v && podman-compose up -d

# 環境変数再読み込み
podman-compose restart api
```

### Go開発

```bash
# APIサーバー起動（ホットリロード）
cd backend
air  # または go run cmd/api/main.go

# テスト実行
go test ./...

# カバレッジ
go test -cover ./...

# Linter
golangci-lint run

# ビルド
go build -o bin/factcheck cmd/cli/main.go
```

### Flutter開発

```bash
# Web開発サーバー
cd frontend/flutter_app
flutter run -d chrome

# ビルド
flutter build web

# テスト
flutter test
```

### データベース

```bash
# PostgreSQL接続
podman exec -it factcheck-postgres psql -U factcheck

# MongoDB接続
podman exec -it factcheck-mongo mongosh

# マイグレーション実行
go run cmd/migrate/main.go up
```

---

## 環境変数管理

- **絶対にコミットしない**: `.env` は `.gitignore` に含める
- **テンプレート使用**: `.env.example` を必ず最新に保つ
- **ローカルオーバーライド**: `.env.local` で個人設定（これも `.gitignore`）

必須環境変数:
```bash
OPENAI_API_KEY=
ANTHROPIC_API_KEY=
DEEPSEEK_API_KEY=
DATABASE_URL=
MONGODB_URL=
REDIS_URL=
```

---

## テスト戦略

### ユニットテスト
- すべての公開関数・メソッドに対してテスト
- テストファイル名: `*_test.go`
- カバレッジ目標: 80%以上

### 統合テスト
- タグ付け: `// +build integration`
- 実際のDBを使用（Docker testcontainers）

### E2Eテスト
- `./scripts/e2e-test.sh` で自動実行
- Playwrightでブラウザ操作

---

## Git ワークフロー

### ブランチ戦略
- `main`: 本番リリース可能
- `develop`: 開発統合ブランチ
- `feature/*`: 新機能（例: `feature/llm-provider-deepseek`）
- `fix/*`: バグ修正
- `docs/*`: ドキュメント更新

### コミットメッセージ
Conventional Commits形式:
```
feat: add DeepSeek LLM provider
fix: resolve race condition in browser pool
docs: update API design documentation
test: add unit tests for claim classifier
```

### PR前チェックリスト
- [ ] `gofmt` / `dart format` 実行
- [ ] Linter警告ゼロ
- [ ] テスト通過
- [ ] `.env.example` 更新（新しい環境変数を追加した場合）
- [ ] ドキュメント更新

---

## セキュリティ注意事項

- **APIキー**: 絶対にハードコードしない
- **ログ**: APIキーやパスワードをログ出力しない
- **入力検証**: すべてのユーザー入力をサニタイズ
- **SQL**: 常にプリペアドステートメント使用
- **依存関係**: 定期的に `go mod tidy` と脆弱性スキャン

---

## LLMプロバイダー使用ガイドライン

### 選択基準
- **OpenAI GPT-4**: 高精度が必要な場合（コスト高）
- **Claude**: 長文分析・倫理的判断
- **DeepSeek**: コスト削減・大量処理

### フォールバック順序
1. ユーザー指定プロバイダー
2. デフォルト（Claude）
3. 利用可能な他のプロバイダー

### レート制限対策
- Redis でキャッシュ（24時間）
- 同一リクエストは再送信しない
- Exponential backoff 実装

---

## トラブルシューティング

### Playwright エラー
```bash
# ブラウザインストール
podman exec -it factcheck-worker playwright install chromium
```

### PostgreSQL 接続エラー
```bash
# コンテナ状態確認
podman ps | grep postgres
# ログ確認
podman logs factcheck-postgres
```

### Redis 接続エラー
```bash
# Redis CLI接続テスト
podman exec -it factcheck-redis redis-cli ping
```

---

## パフォーマンス最適化

- **DB接続**: コネクションプール使用（最大20接続）
- **Playwright**: ブラウザインスタンス再利用（最大5並列）
- **LLMキャッシュ**: 同一クエリは24時間キャッシュ
- **画像**: 不要な場合はPlaywrightでブロック

---

## デプロイ

### OSS版（Podman）
```bash
podman-compose -f docker-compose.prod.yml up -d
```

### GCP本番環境
```bash
# Cloud Run デプロイ
gcloud run deploy factcheck-api --source .
```

詳細: `docs/deployment.md`

---

## 参考リンク

- [Go プロジェクトレイアウト](https://github.com/golang-standards/project-layout)
- [Effective Go](https://go.dev/doc/effective_go)
- [Semantic Scholar API](https://api.semanticscholar.org/api-docs/)
- [Playwright Go](https://playwright.dev/docs/intro)
