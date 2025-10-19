package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/sourcetracer/internal/analyzer"
	"github.com/yourusername/sourcetracer/internal/classifier"
	"github.com/yourusername/sourcetracer/internal/config"
	"github.com/yourusername/sourcetracer/internal/domain"
)

const version = "v0.2.0-alpha"

func main() {
	// コマンドラインフラグ
	text := flag.String("text", "", "Text to analyze")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("SourceTracer CLI %s\n", version)
		os.Exit(0)
	}

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔍 SourceTracer CLI %s\n", version)
	fmt.Printf("📝 Environment: %s\n", cfg.AppEnv)
	fmt.Printf("🤖 Default LLM: %s\n\n", cfg.DefaultLLMProvider)

	// テキストが指定されていない場合はサンプルテキストを使用
	if *text == "" {
		*text = "Climate change is accelerating faster than predicted."
		fmt.Printf("No text provided, using sample: %q\n\n", *text)
	}

	// Analyze
	c := classifier.NewClassifier()
	a := analyzer.NewAnalyzer(c, nil) // No search clients for now

	ctx := context.Background()
	opts := analyzer.AnalyzeOptions{
		IncludeEvidences: false,
		MaxClaims:        10,
	}

	result, err := a.Analyze(ctx, *text, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing text: %v\n", err)
		os.Exit(1)
	}

	// 結果表示
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("📄 Original Text: %s\n", *text)
	fmt.Println(strings.Repeat("-", 61))
	fmt.Printf("📊 Analysis ID: %s\n", result.AnalysisID)
	fmt.Printf("⏱️  Processing Time: %dms\n", result.ProcessingTimeMS)
	fmt.Printf("📈 Overall Credibility: %.2f\n\n", result.OverallCredibility)

	// 各クレームの表示
	fmt.Printf("📋 Claims Found: %d\n\n", len(result.Claims))
	for i, claim := range result.Claims {
		fmt.Printf("Claim #%d:\n", i+1)
		fmt.Printf("  Text: %s\n", claim.Text)
		fmt.Printf("  Type: %s\n", claim.Type.String())
		fmt.Printf("  Confidence: %.2f\n", claim.Confidence)

		switch claim.Type {
		case domain.Opinion:
			fmt.Println("  ✅ OPINION - Consider providing evidence")
		case domain.Fact:
			fmt.Println("  ✅ FACT - Verify with reliable sources")
		case domain.Mixed:
			fmt.Println("  ⚠️  MIXED - Contains opinion and fact")
		case domain.Unclear:
			fmt.Println("  ❌ UNCLEAR - Needs more context")
		}
		fmt.Println()
	}

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("\n🎯 Next steps:")
	fmt.Println("   - Search for supporting evidence (use --search flag)")
	fmt.Println("   - Verify credibility of sources")
	fmt.Println("   - Check for counter-evidence")
}
