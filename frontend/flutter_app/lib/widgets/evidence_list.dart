import 'package:flutter/material.dart';
import '../models/evidence.dart';
import 'package:intl/intl.dart';

/// List of evidences
class EvidenceList extends StatelessWidget {
  final String title;
  final List<Evidence> evidences;
  final bool isCounter;

  const EvidenceList({
    Key? key,
    required this.title,
    required this.evidences,
    this.isCounter = false,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Row(
            children: [
              Icon(
                isCounter ? Icons.cancel : Icons.check_circle,
                color: isCounter ? Colors.red : Colors.green,
                size: 20,
              ),
              const SizedBox(width: 8),
              Text(
                title,
                style: const TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
              ),
              const SizedBox(width: 8),
              Chip(
                label: Text(
                  '${evidences.length}件',
                  style: const TextStyle(fontSize: 11),
                ),
                padding: EdgeInsets.zero,
                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
              ),
            ],
          ),
        ),
        ...evidences.map((e) => EvidenceCard(evidence: e)),
      ],
    );
  }
}

/// Card displaying a single evidence
class EvidenceCard extends StatelessWidget {
  final Evidence evidence;

  const EvidenceCard({Key? key, required this.evidence}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey[50],
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey[300]!),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header: Source + Type
          Row(
            children: [
              Chip(
                label: Text(
                  evidence.source,
                  style: const TextStyle(fontSize: 11),
                ),
                backgroundColor: Colors.blue[100],
                padding: EdgeInsets.zero,
                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
              ),
              const SizedBox(width: 8),
              Chip(
                label: Text(
                  evidence.sourceTypeDisplayName,
                  style: const TextStyle(fontSize: 11),
                ),
                backgroundColor: _getSourceTypeColor(),
                padding: EdgeInsets.zero,
                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
              ),
            ],
          ),
          const SizedBox(height: 8),

          // Title
          Text(
            evidence.title,
            style: const TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 4),

          // Authors
          if (evidence.authors.isNotEmpty)
            Text(
              evidence.authors.join(', '),
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[700],
                fontStyle: FontStyle.italic,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          const SizedBox(height: 8),

          // Snippet
          if (evidence.snippet.isNotEmpty)
            Text(
              evidence.snippet,
              style: const TextStyle(fontSize: 13),
              maxLines: 3,
              overflow: TextOverflow.ellipsis,
            ),
          const SizedBox(height: 8),

          // Metadata row
          Wrap(
            spacing: 12,
            runSpacing: 4,
            children: [
              _buildMetadata(
                Icons.star,
                '信頼性: ${evidence.credibilityLevel}',
                _getCredibilityColor(),
              ),
              _buildMetadata(
                Icons.link,
                '関連性: ${(evidence.relevance * 100).toStringAsFixed(0)}%',
                Colors.blue,
              ),
              if (evidence.citationCount != null)
                _buildMetadata(
                  Icons.format_quote,
                  '引用: ${evidence.citationCount}',
                  Colors.purple,
                ),
              if (evidence.publishedDate != null)
                _buildMetadata(
                  Icons.calendar_today,
                  _formatDate(evidence.publishedDate!),
                  Colors.grey,
                ),
            ],
          ),

          // URL
          const SizedBox(height: 8),
          InkWell(
            onTap: () {
              // TODO: Open URL in browser
            },
            child: Text(
              evidence.url,
              style: const TextStyle(
                color: Colors.blue,
                decoration: TextDecoration.underline,
                fontSize: 12,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMetadata(IconData icon, String text, Color color) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: color),
        const SizedBox(width: 4),
        Text(
          text,
          style: TextStyle(fontSize: 11, color: color),
        ),
      ],
    );
  }

  Color _getSourceTypeColor() {
    switch (evidence.sourceType) {
      case 'peer_reviewed_journal':
        return Colors.green[100]!;
      case 'preprint':
        return Colors.yellow[100]!;
      case 'government':
        return Colors.blue[100]!;
      case 'news':
        return Colors.orange[100]!;
      default:
        return Colors.grey[100]!;
    }
  }

  Color _getCredibilityColor() {
    if (evidence.credibility >= 0.8) return Colors.green;
    if (evidence.credibility >= 0.6) return Colors.orange;
    return Colors.red;
  }

  String _formatDate(String dateStr) {
    try {
      final date = DateTime.parse(dateStr);
      return DateFormat('yyyy/MM/dd').format(date);
    } catch (e) {
      return dateStr;
    }
  }
}
