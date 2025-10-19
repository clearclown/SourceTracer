import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/analysis_provider.dart';
import '../widgets/claim_card.dart';

/// Analysis screen with text input and results
class AnalysisScreen extends StatefulWidget {
  const AnalysisScreen({Key? key}) : super(key: key);

  @override
  State<AnalysisScreen> createState() => _AnalysisScreenState();
}

class _AnalysisScreenState extends State<AnalysisScreen> {
  final TextEditingController _textController = TextEditingController();
  bool _includeEvidences = false;
  String? _selectedLLM;
  final List<String> _llmProviders = ['claude', 'openai', 'deepseek'];

  @override
  void dispose() {
    _textController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<AnalysisProvider>(
      builder: (context, provider, child) {
        return Column(
          children: [
            // Input Section
            Container(
              padding: const EdgeInsets.all(16),
              color: Colors.grey[100],
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'テキストを分析',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    controller: _textController,
                    maxLines: 5,
                    decoration: InputDecoration(
                      hintText: 'ここに分析したいテキストを入力してください...\n\n'
                          '例: 「気候変動は加速している。Pythonは最高のプログラミング言語だ。」',
                      border: const OutlineInputBorder(),
                      filled: true,
                      fillColor: Colors.white,
                      enabled: !provider.isLoading,
                    ),
                  ),
                  const SizedBox(height: 12),

                  // Options
                  Wrap(
                    spacing: 16,
                    runSpacing: 8,
                    children: [
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Checkbox(
                            value: _includeEvidences,
                            onChanged: provider.isLoading
                                ? null
                                : (value) {
                                    setState(() {
                                      _includeEvidences = value ?? false;
                                    });
                                  },
                          ),
                          const Text('エビデンスを検索'),
                        ],
                      ),
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Text('LLM: '),
                          DropdownButton<String>(
                            value: _selectedLLM,
                            hint: const Text('自動'),
                            items: _llmProviders.map((String value) {
                              return DropdownMenuItem<String>(
                                value: value,
                                child: Text(value),
                              );
                            }).toList(),
                            onChanged: provider.isLoading
                                ? null
                                : (value) {
                                    setState(() {
                                      _selectedLLM = value;
                                    });
                                  },
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),

                  // Analyze Button
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton.icon(
                      onPressed: provider.isLoading
                          ? null
                          : () => _analyzeText(provider),
                      icon: provider.isLoading
                          ? const SizedBox(
                              width: 16,
                              height: 16,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: Colors.white,
                              ),
                            )
                          : const Icon(Icons.analytics),
                      label: Text(
                        provider.isLoading ? '分析中...' : '分析する',
                        style: const TextStyle(fontSize: 16),
                      ),
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // Results Section
            Expanded(
              child: _buildResultsSection(provider),
            ),
          ],
        );
      },
    );
  }

  Widget _buildResultsSection(AnalysisProvider provider) {
    if (provider.isLoading) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 16),
            Text('分析中...しばらくお待ちください'),
          ],
        ),
      );
    }

    if (provider.hasError) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 64, color: Colors.red),
              const SizedBox(height: 16),
              const Text(
                'エラーが発生しました',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                provider.errorMessage ?? 'Unknown error',
                textAlign: TextAlign.center,
                style: const TextStyle(color: Colors.red),
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => provider.clearResult(),
                child: const Text('閉じる'),
              ),
            ],
          ),
        ),
      );
    }

    if (provider.currentResult == null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.article, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              'テキストを入力して分析を開始してください',
              style: TextStyle(
                fontSize: 16,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      );
    }

    final result = provider.currentResult!;

    return ListView(
      padding: const EdgeInsets.symmetric(vertical: 16),
      children: [
        // Summary Card
        Container(
          margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.blue[50],
            borderRadius: BorderRadius.circular(8),
            border: Border.all(color: Colors.blue[200]!),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  const Icon(Icons.assessment, color: Colors.blue),
                  const SizedBox(width: 8),
                  const Text(
                    '分析サマリー',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
              const Divider(),
              _buildSummaryRow(
                '総クレーム数',
                '${result.claims.length}件',
                Icons.format_list_numbered,
              ),
              _buildSummaryRow(
                '意見',
                '${result.opinionCount}件',
                Icons.comment,
              ),
              _buildSummaryRow(
                '事実',
                '${result.factCount}件',
                Icons.check_circle,
              ),
              _buildSummaryRow(
                '総合信頼性',
                '${(result.overallCredibility * 100).toStringAsFixed(1)}% (${result.credibilityLevel})',
                Icons.star,
              ),
              _buildSummaryRow(
                '処理時間',
                '${result.processingTimeMs}ms',
                Icons.timer,
              ),
              if (result.totalEvidences > 0)
                _buildSummaryRow(
                  'エビデンス',
                  '${result.totalEvidences}件',
                  Icons.library_books,
                ),
            ],
          ),
        ),

        // Claims
        const Padding(
          padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Text(
            'クレーム一覧',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        ...result.claims.map((claim) => ClaimCard(
              claim: claim,
              showEvidences: _includeEvidences,
            )),
      ],
    );
  }

  Widget _buildSummaryRow(String label, String value, IconData icon) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          Icon(icon, size: 18, color: Colors.blue[700]),
          const SizedBox(width: 8),
          Text(
            '$label: ',
            style: const TextStyle(fontWeight: FontWeight.w500),
          ),
          Text(value),
        ],
      ),
    );
  }

  void _analyzeText(AnalysisProvider provider) {
    final text = _textController.text.trim();

    if (text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('テキストを入力してください')),
      );
      return;
    }

    provider.analyzeText(
      text: text,
      includeEvidences: _includeEvidences,
      useLLM: _selectedLLM != null,
      llmProvider: _selectedLLM,
    );
  }
}
