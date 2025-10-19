# SourceTracer - Citation and Evidence Analysis System

> **"情報源を、徹底的に追う。"**

**自動文献収集・事実検証・情報源追跡システム**

SourceTracerは、学術論文や記事執筆時に、ユーザーの主張をサポートする文献を自動収集し、意見と事実を峻別して信頼性分析を行うオープンソースシステムです。

### プロジェクトの使命

現代の情報社会において、**意見と事実の境界が曖昧**になり、根拠なき主張が氾濫しています。

**SourceTracerが解決する課題:**
- ✅ **意見と事実の峻別**: 人間が曖昧にしがちな境界を自動判定
- ✅ **批判的分析**: 主張の弱点を多角的に指摘
- ✅ **徹底的なソース追跡**: 複数検索エンジン・データベースを横断
- ✅ **学術的誠実性**: 盗作・盗説を防止し、適切な引用を促進

### なぜSourceTracerが重要か

このツールが不完全な場合、以下の深刻な問題が生じます：
- ❌ **誤情報の拡散**: 低品質なエビデンスを「検証済み」と誤認
- ❌ **盗作・盗説の助長**: 適切な引用なしにコンテンツを流用
- ❌ **情報操作**: 偏ったソース選択により特定の主張を不当に強化
- ❌ **学術的信頼の崩壊**: 誤った判定により研究の信頼性が失われる

**だからこそ、品質への妥協は許されません。**

### 開発哲学

- ✅ **TDD厳守**（t-wada流） - テスト駆動開発の徹底
- ✅ **ゼロ・トレランス** - モック・スクリプトでのごまかし禁止
- ✅ **完全な透明性** - すべての判定根拠を検証可能に
- ✅ **倫理的配慮** - バイアスを最小化し、公正な評価を実現

## 🎯 主要機能

### 1. 意見と事実の峻別
- LLMを用いた高度な意味解析
- Claim（主張）の自動抽出と分類
- Opinion / Fact / Mixed の3段階判定
- 人間が曖昧にする境界を明確化

### 2. 徹底的な情報源追跡
**多角的データ収集戦略:**
- **学術データベース**: Semantic Scholar, arXiv, PubMed, CORE
- **汎用検索**: Google Custom Search, DuckDuckGo, Baidu
- **Playwright自動化**: JavaScript重視サイトも確実にスクレイピング
- **クロスリファレンス**: 複数ソースで事実を相互検証

### 3. 批判的分析レポート
- **主張ごとの検証**: 個別の主張に対する賛成・反対エビデンス
- **弱点の指摘**: ロジックの穴、不十分な根拠を明示
- **対立意見の提示**: バランスの取れた視点を提供
- **バイアス検出**: 偏ったソース選択を警告

### 4. 信頼性スコアリング
- **多次元評価**:
  - ソースタイプ（査読論文 > プレプリント > ニュース > ブログ）
  - 被引用数・著者評価
  - 公開日時・情報の新しさ
  - ドメイン信頼度・HTTPS対応
- **ユーザー定義フィルター**: 信頼性閾値のカスタマイズ可能
- **透明なアルゴリズム**: スコア計算ロジックを完全公開

---

## 🏗️ アーキテクチャ

```
[Flutter Frontend] ←→ [Go API Gateway] ←→ [Worker Pool (Playwright)]
                            ↓
            ┌───────────────┼───────────────┐
            ↓               ↓               ↓
      [PostgreSQL]     [MongoDB]        [Redis]
      (メタデータ)    (生データ)       (キャッシュ)
```

### 技術スタック
- **Backend**: Go 1.21+, Gin, Playwright
- **Frontend**: Flutter (Web/Desktop/Mobile)
- **Database**: PostgreSQL (pgvector), MongoDB, Redis
- **LLM**: OpenAI, Claude, DeepSeek（切り替え可能）
- **Deploy**: Podman/Docker, GCP (Cloud Run, Cloud SQL)

---

## 🚀 クイックスタート

### 前提条件
- Go 1.21+
- Podman or Docker
- Flutter SDK 3.16+

### ローカル開発

1. **リポジトリクローン**
```bash
git clone https://github.com/yourusername/sourcetracer.git
cd sourcetracer
```

2. **環境変数設定**
```bash
cp .env.example .env
# .envを編集してAPIキーを設定
```

3. **コンテナ起動（Podman Compose推奨）**
```bash
podman compose up -d
# または Docker Composeの場合
docker-compose up -d
```

4. **CLIテスト**
```bash
cd backend
go run cmd/cli/main.go analyze --file sample.txt
```

5. **Webアクセス**
```
http://localhost:3000  # Flutter Frontend
http://localhost:8080  # API Server
```

---

## 📝 使用例

### CLI
```bash
# ファイル分析
sourcetracer analyze --file article.md --llm claude

# 標準入力から意見検出
echo "地球は平面である" | sourcetracer detect-opinion --stdin

# ソース追跡（徹底モード）
sourcetracer trace "Climate change is accelerating" \
  --engines semantic_scholar,arxiv,google,duckduckgo,baidu \
  --min-credibility 0.8

# JSON出力
sourcetracer analyze --file paper.txt --output json > result.json

# 批判的分析
sourcetracer critique --file draft.md --include-counter-evidence
```

### API (v0.3+)

**サーバー起動:**
```bash
cd backend
go run cmd/api/main.go
```

**APIリクエスト例:**
```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Python is the best language. Climate change is real.",
    "options": {
      "include_evidences": false,
      "max_claims": 10
    }
  }'
```

**レスポンス例:**
```json
{
  "success": true,
  "data": {
    "analysis_id": "anl_1760860015098070181",
    "claims": [
      {
        "id": "clm_abc123",
        "text": "Python is the best language.",
        "type": "opinion",
        "confidence": 0.9
      },
      {
        "id": "clm_def456",
        "text": "Climate change is real.",
        "type": "fact",
        "confidence": 0.7
      }
    ],
    "overall_credibility": 0.8,
    "processing_time_ms": 5
  }
}
```

---

## 🗂️ プロジェクト構造

```
sourcetracer/
├── backend/          # Go APIサーバー、Worker、CLI
├── frontend/         # Flutter Webアプリ
├── docs/             # 設計ドキュメント
├── docker-compose.yml
└── podman-compose.yml
```

詳細は `docs/architecture.md` を参照。

---

## 🧪 テスト

```bash
# ユニットテスト
cd backend
go test ./...

# 統合テスト
go test -tags=integration ./...

# E2Eテスト（要Docker）
./scripts/e2e-test.sh
```

---

## 📊 データベース

- **PostgreSQL**: ユーザー、検索履歴、引用管理
- **MongoDB**: スクレイピング生データ、HTMLスナップショット
- **Redis**: LLMレスポンスキャッシュ、セッション管理

スキーマ詳細: `docs/database-schema.md`

---

## 🔐 セキュリティ

- APIキーは `.env` で管理（絶対にコミットしない）
- HTTPS必須（本番環境）
- レート制限：100 req/min（未認証）、1000 req/min（認証済み）
- ユーザー入力の厳密なバリデーション

---

## 🛣️ ロードマップ

- [x] 基本アーキテクチャ設計
- [x] **v0.1: CLI + 基本分析機能** ✅
  - [x] ドメインモデル (Claim, Evidence)
  - [x] ルールベース分類器
  - [x] 設定管理
  - [x] CLIツール
- [x] **v0.2: LLM統合 + Evidence検索** ✅
  - [x] Claude API統合
  - [x] Semantic Scholar検索
  - [x] Analyzer (複数クレーム抽出)
- [x] **v0.3: Web API** ✅
  - [x] Gin REST API server
  - [x] `/api/v1/analyze` endpoint
  - [x] CORS対応
  - [x] エラーハンドリング
- [x] **v0.4-v0.5: データベース統合** ✅
  - [x] PostgreSQL統合
  - [x] 検索履歴保存
  - [x] Repository pattern実装
- [x] **v0.6: 検索プロバイダー拡張** ✅
  - [x] arXiv検索統合
  - [x] Google Custom Search統合
  - [x] マルチプロバイダーアグリゲーター
- [x] **v0.7: バックエンド完全統合** ✅
  - [x] SearchAggregator統合
  - [x] PostgreSQL自動接続
  - [x] 環境変数ベース設定
- [x] **v0.8: LLM + デプロイメント** ✅ (現在のバージョン)
  - [x] DeepSeek LLMプロバイダー
  - [x] Podman Compose設定
  - [x] マルチステージDockerfile
  - [x] 本番環境対応
- [ ] v0.9: Flutter Web UI
  - [ ] テキスト入力フォーム
  - [ ] 分析結果表示
  - [ ] 履歴ビュー
- [ ] v1.0: OSS公開
  - [ ] ドキュメント完成
  - [ ] CI/CD完全自動化
  - [ ] パフォーマンス最適化

---

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. フィーチャーブランチ作成 (`git checkout -b feature/amazing-feature`)
3. コミット (`git commit -m 'Add amazing feature'`)
4. プッシュ (`git push origin feature/amazing-feature`)
5. Pull Request作成

詳細: `CONTRIBUTING.md`

---

## 📄 ライセンス

MIT License - 詳細は `LICENSE` ファイルを参照

---

## 👨‍💻 作者

- **あなたの名前** - [GitHub](https://github.com/yourusername)

---

## 🙏 謝辞

- Semantic Scholar API
- Playwright Team
- Anthropic Claude
- OpenAI

---

## 📞 サポート

- Issue: https://github.com/yourusername/sourcetracer/issues
- Email: your-email@example.com
- Discord: [コミュニティリンク]

---

**SourceTracer** - 情報源を、徹底的に追う。すべての主張にエビデンスを。
