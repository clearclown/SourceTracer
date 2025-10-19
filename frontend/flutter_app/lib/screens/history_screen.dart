import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/analysis_provider.dart';
import '../models/analysis_result.dart';

/// History screen showing past analyses
class HistoryScreen extends StatefulWidget {
  const HistoryScreen({Key? key}) : super(key: key);

  @override
  State<HistoryScreen> createState() => _HistoryScreenState();
}

class _HistoryScreenState extends State<HistoryScreen> {
  @override
  void initState() {
    super.initState();
    // Load history when screen is first displayed
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<AnalysisProvider>().loadHistory();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<AnalysisProvider>(
      builder: (context, provider, child) {
        if (provider.history.isEmpty) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.history, size: 64, color: Colors.grey[400]),
                const SizedBox(height: 16),
                Text(
                  '履歴がありません',
                  style: TextStyle(
                    fontSize: 16,
                    color: Colors.grey[600],
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  '分析を実行すると、ここに履歴が表示されます',
                  style: TextStyle(fontSize: 14, color: Colors.grey),
                ),
              ],
            ),
          );
        }

        return RefreshIndicator(
          onRefresh: () => provider.loadHistory(),
          child: ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: provider.history.length,
            itemBuilder: (context, index) {
              final result = provider.history[index];
              return HistoryCard(
                result: result,
                onTap: () => _showAnalysisDetails(context, result),
              );
            },
          ),
        );
      },
    );
  }

  void _showAnalysisDetails(BuildContext context, AnalysisResult result) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.9,
        minChildSize: 0.5,
        maxChildSize: 0.95,
        expand: false,
        builder: (context, scrollController) {
          return Container(
            padding: const EdgeInsets.all(16),
            child: ListView(
              controller: scrollController,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      '分析詳細',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.close),
                      onPressed: () => Navigator.pop(context),
                    ),
                  ],
                ),
                const Divider(),
                _buildDetailSection(result),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildDetailSection(AnalysisResult result) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        _buildDetailRow('分析ID', result.analysisId),
        _buildDetailRow('クレーム数', '${result.claims.length}件'),
        _buildDetailRow('意見', '${result.opinionCount}件'),
        _buildDetailRow('事実', '${result.factCount}件'),
        _buildDetailRow(
          '総合信頼性',
          '${(result.overallCredibility * 100).toStringAsFixed(1)}%',
        ),
        _buildDetailRow('処理時間', '${result.processingTimeMs}ms'),
        if (result.createdAt != null)
          _buildDetailRow('作成日時', result.createdAt!),
        const SizedBox(height: 16),
        const Text(
          'クレーム一覧',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        ...result.claims.map((claim) => Card(
              child: ListTile(
                leading: Icon(
                  claim.type.toLowerCase() == 'opinion'
                      ? Icons.comment
                      : Icons.check_circle,
                  color: claim.type.toLowerCase() == 'opinion'
                      ? Colors.orange
                      : Colors.green,
                ),
                title: Text(claim.text),
                subtitle: Text(
                  '${claim.typeDisplayName} - 信頼度: ${(claim.confidence * 100).toStringAsFixed(0)}%',
                ),
              ),
            )),
      ],
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(
              '$label:',
              style: const TextStyle(fontWeight: FontWeight.w500),
            ),
          ),
          Expanded(
            child: Text(value),
          ),
        ],
      ),
    );
  }
}

/// Card for displaying a history item
class HistoryCard extends StatelessWidget {
  final AnalysisResult result;
  final VoidCallback onTap;

  const HistoryCard({
    Key? key,
    required this.result,
    required this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(4),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Text(
                      'クレーム ${result.claims.length}件',
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  _buildCredibilityBadge(),
                ],
              ),
              const SizedBox(height: 8),

              // Preview of first claim
              if (result.claims.isNotEmpty)
                Text(
                  result.claims.first.text,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    color: Colors.grey[700],
                  ),
                ),
              const SizedBox(height: 8),

              // Metadata
              Wrap(
                spacing: 12,
                runSpacing: 4,
                children: [
                  _buildMetadataChip(
                    Icons.comment,
                    '意見: ${result.opinionCount}',
                    Colors.orange,
                  ),
                  _buildMetadataChip(
                    Icons.check_circle,
                    '事実: ${result.factCount}',
                    Colors.green,
                  ),
                  _buildMetadataChip(
                    Icons.timer,
                    '${result.processingTimeMs}ms',
                    Colors.blue,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCredibilityBadge() {
    Color color;
    if (result.overallCredibility >= 0.8) {
      color = Colors.green;
    } else if (result.overallCredibility >= 0.6) {
      color = Colors.orange;
    } else {
      color = Colors.red;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color),
      ),
      child: Text(
        '${(result.overallCredibility * 100).toStringAsFixed(0)}%',
        style: TextStyle(
          color: color,
          fontWeight: FontWeight.bold,
          fontSize: 14,
        ),
      ),
    );
  }

  Widget _buildMetadataChip(IconData icon, String label, Color color) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: color),
        const SizedBox(width: 4),
        Text(
          label,
          style: TextStyle(fontSize: 12, color: color),
        ),
      ],
    );
  }
}
