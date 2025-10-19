import 'package:json_annotation/json_annotation.dart';

part 'evidence.g.dart';

@JsonSerializable()
class Evidence {
  final String id;
  final String source;
  final String title;
  final List<String> authors;
  final String url;
  final String snippet;
  final double credibility;
  final double relevance;

  @JsonKey(name: 'published_date')
  final String? publishedDate;

  @JsonKey(name: 'source_type')
  final String sourceType;

  @JsonKey(name: 'citation_count')
  final int? citationCount;

  final Map<String, dynamic>? metadata;

  Evidence({
    required this.id,
    required this.source,
    required this.title,
    required this.authors,
    required this.url,
    required this.snippet,
    required this.credibility,
    required this.relevance,
    this.publishedDate,
    required this.sourceType,
    this.citationCount,
    this.metadata,
  });

  factory Evidence.fromJson(Map<String, dynamic> json) =>
      _$EvidenceFromJson(json);

  Map<String, dynamic> toJson() => _$EvidenceToJson(this);

  // Helper methods
  String get sourceTypeDisplayName {
    switch (sourceType) {
      case 'peer_reviewed_journal':
        return '査読論文';
      case 'preprint':
        return 'プレプリント';
      case 'government':
        return '政府機関';
      case 'news':
        return 'ニュース';
      case 'blog':
        return 'ブログ';
      default:
        return '不明';
    }
  }

  String get credibilityLevel {
    if (credibility >= 0.9) return '非常に高い';
    if (credibility >= 0.7) return '高い';
    if (credibility >= 0.5) return '中程度';
    if (credibility >= 0.3) return '低い';
    return '非常に低い';
  }
}
