package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/sourcetracer/internal/classifier"
	"github.com/yourusername/sourcetracer/internal/config"
	"github.com/yourusername/sourcetracer/internal/domain"
)

const version = "v0.1.0-alpha"

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

	// Classify
	c := classifier.NewClassifier()
	result := c.Classify(*text)

	// 結果表示
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("📄 Text: %s\n", *text)
	fmt.Println(strings.Repeat("-", 61))
	fmt.Printf("📊 Type: %s\n", result.Type.String())
	fmt.Printf("📈 Confidence: %.2f\n", result.Confidence)
	fmt.Printf("💡 Reasoning: %s\n", result.Reasoning)
	fmt.Println("=" + strings.Repeat("=", 60))

	// タイプ別のメッセージ
	switch result.Type {
	case domain.Opinion:
		fmt.Println("\n✅ This appears to be an OPINION.")
		fmt.Println("   Consider providing evidence or rephrasing as fact.")
	case domain.Fact:
		fmt.Println("\n✅ This appears to be a FACT.")
		fmt.Println("   Verify with reliable sources for accuracy.")
	case domain.Mixed:
		fmt.Println("\n⚠️  This contains both opinion and fact.")
		fmt.Println("   Consider separating them for clarity.")
	case domain.Unclear:
		fmt.Println("\n❌ Unable to classify this claim.")
		fmt.Println("   Please provide more context.")
	}

	fmt.Println("\n🎯 Next steps:")
	fmt.Println("   - Search for supporting evidence")
	fmt.Println("   - Verify credibility of sources")
	fmt.Println("   - Check for counter-evidence")
}
