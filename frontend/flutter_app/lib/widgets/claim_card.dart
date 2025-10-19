import 'package:flutter/material.dart';
import '../models/claim.dart';
import 'evidence_list.dart';

/// Card displaying a single claim
class ClaimCard extends StatelessWidget {
  final Claim claim;
  final bool showEvidences;

  const ClaimCard({
    Key? key,
    required this.claim,
    this.showEvidences = true,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
      elevation: 2,
      child: ExpansionTile(
        leading: _buildTypeIcon(),
        title: Text(
          claim.text,
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w500,
          ),
        ),
        subtitle: Padding(
          padding: const EdgeInsets.only(top: 8),
          child: Wrap(
            spacing: 8,
            runSpacing: 4,
            children: [
              _buildChip(
                claim.typeDisplayName,
                _getTypeColor(),
              ),
              _buildChip(
                '信頼度: ${(claim.confidence * 100).toStringAsFixed(0)}%',
                Colors.blue,
              ),
              if (claim.totalEvidences > 0)
                _buildChip(
                  'エビデンス: ${claim.totalEvidences}件',
                  Colors.green,
                ),
            ],
          ),
        ),
        children: [
          if (claim.summary != null)
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    '要約',
                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(claim.summary!),
                ],
              ),
            ),
          if (claim.recommendation != null)
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              color: _getRecommendationColor().withOpacity(0.1),
              child: Row(
                children: [
                  Icon(
                    _getRecommendationIcon(),
                    color: _getRecommendationColor(),
                    size: 20,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    '評価: ${claim.recommendationDisplayName}',
                    style: TextStyle(
                      color: _getRecommendationColor(),
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
          if (showEvidences && claim.evidences.isNotEmpty)
            EvidenceList(
              title: '支持するエビデンス',
              evidences: claim.evidences,
            ),
          if (showEvidences &&
              claim.counterEvidences != null &&
              claim.counterEvidences!.isNotEmpty)
            EvidenceList(
              title: '反対するエビデンス',
              evidences: claim.counterEvidences!,
              isCounter: true,
            ),
        ],
      ),
    );
  }

  Widget _buildTypeIcon() {
    IconData icon;
    Color color;

    switch (claim.type.toLowerCase()) {
      case 'opinion':
        icon = Icons.comment;
        color = Colors.orange;
        break;
      case 'fact':
        icon = Icons.check_circle;
        color = Colors.green;
        break;
      case 'mixed':
        icon = Icons.merge_type;
        color = Colors.purple;
        break;
      default:
        icon = Icons.help;
        color = Colors.grey;
    }

    return Icon(icon, color: color);
  }

  Color _getTypeColor() {
    switch (claim.type.toLowerCase()) {
      case 'opinion':
        return Colors.orange;
      case 'fact':
        return Colors.green;
      case 'mixed':
        return Colors.purple;
      default:
        return Colors.grey;
    }
  }

  Color _getRecommendationColor() {
    switch (claim.recommendation?.toUpperCase()) {
      case 'SUPPORTED':
        return Colors.green;
      case 'UNSUPPORTED':
        return Colors.red;
      case 'MIXED':
        return Colors.orange;
      default:
        return Colors.grey;
    }
  }

  IconData _getRecommendationIcon() {
    switch (claim.recommendation?.toUpperCase()) {
      case 'SUPPORTED':
        return Icons.thumb_up;
      case 'UNSUPPORTED':
        return Icons.thumb_down;
      case 'MIXED':
        return Icons.thumbs_up_down;
      default:
        return Icons.help_outline;
    }
  }

  Widget _buildChip(String label, Color color) {
    return Chip(
      label: Text(
        label,
        style: const TextStyle(fontSize: 12),
      ),
      backgroundColor: color.withOpacity(0.1),
      side: BorderSide(color: color, width: 1),
      padding: EdgeInsets.zero,
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    );
  }
}
