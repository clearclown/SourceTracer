import 'package:json_annotation/json_annotation.dart';
import 'evidence.dart';

part 'claim.g.dart';

@JsonSerializable(explicitToJson: true)
class Claim {
  final String id;
  final String text;
  final String type; // opinion, fact, mixed, unclear
  final double confidence;

  final Position? position;
  final List<Evidence> evidences;

  @JsonKey(name: 'counter_evidences')
  final List<Evidence>? counterEvidences;

  final String? summary;
  final String? recommendation; // SUPPORTED, UNSUPPORTED, MIXED, INSUFFICIENT

  Claim({
    required this.id,
    required this.text,
    required this.type,
    required this.confidence,
    this.position,
    this.evidences = const [],
    this.counterEvidences,
    this.summary,
    this.recommendation,
  });

  factory Claim.fromJson(Map<String, dynamic> json) => _$ClaimFromJson(json);

  Map<String, dynamic> toJson() => _$ClaimToJson(this);

  // Helper methods
  String get typeDisplayName {
    switch (type.toLowerCase()) {
      case 'opinion':
        return '意見';
      case 'fact':
        return '事実';
      case 'mixed':
        return '混合';
      case 'unclear':
        return '不明確';
      default:
        return type;
    }
  }

  String get recommendationDisplayName {
    switch (recommendation?.toUpperCase()) {
      case 'SUPPORTED':
        return '支持されている';
      case 'UNSUPPORTED':
        return '支持されていない';
      case 'MIXED':
        return '賛否両論';
      case 'INSUFFICIENT':
        return 'エビデンス不足';
      default:
        return '未評価';
    }
  }

  int get totalEvidences =>
      evidences.length + (counterEvidences?.length ?? 0);
}

@JsonSerializable()
class Position {
  final int start;
  final int end;

  Position({required this.start, required this.end});

  factory Position.fromJson(Map<String, dynamic> json) =>
      _$PositionFromJson(json);

  Map<String, dynamic> toJson() => _$PositionToJson(this);
}
