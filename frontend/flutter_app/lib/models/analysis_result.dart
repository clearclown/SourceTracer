import 'package:json_annotation/json_annotation.dart';
import 'claim.dart';

part 'analysis_result.g.dart';

@JsonSerializable(explicitToJson: true)
class AnalysisResult {
  @JsonKey(name: 'analysis_id')
  final String analysisId;

  final List<Claim> claims;

  @JsonKey(name: 'overall_credibility')
  final double overallCredibility;

  @JsonKey(name: 'processing_time_ms')
  final int processingTimeMs;

  @JsonKey(name: 'created_at')
  final String? createdAt;

  AnalysisResult({
    required this.analysisId,
    required this.claims,
    required this.overallCredibility,
    required this.processingTimeMs,
    this.createdAt,
  });

  factory AnalysisResult.fromJson(Map<String, dynamic> json) =>
      _$AnalysisResultFromJson(json);

  Map<String, dynamic> toJson() => _$AnalysisResultToJson(this);

  // Helper methods
  int get totalEvidences => claims.fold(
      0, (sum, claim) => sum + claim.evidences.length);

  int get opinionCount =>
      claims.where((c) => c.type.toLowerCase() == 'opinion').length;

  int get factCount =>
      claims.where((c) => c.type.toLowerCase() == 'fact').length;

  String get credibilityLevel {
    if (overallCredibility >= 0.8) return '高信頼性';
    if (overallCredibility >= 0.6) return '中信頼性';
    if (overallCredibility >= 0.4) return '低信頼性';
    return '要注意';
  }
}

@JsonSerializable(genericArgumentFactories: true)
class ApiResponse<T> {
  final bool success;
  final T? data;
  final ApiError? error;
  final ApiMetadata? metadata;

  ApiResponse({
    required this.success,
    this.data,
    this.error,
    this.metadata,
  });

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) =>
      _$ApiResponseFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(Object? Function(T value) toJsonT) =>
      _$ApiResponseToJson(this, toJsonT);
}

@JsonSerializable()
class ApiError {
  final String code;
  final String message;
  final Map<String, dynamic>? details;

  ApiError({
    required this.code,
    required this.message,
    this.details,
  });

  factory ApiError.fromJson(Map<String, dynamic> json) =>
      _$ApiErrorFromJson(json);

  Map<String, dynamic> toJson() => _$ApiErrorToJson(this);
}

@JsonSerializable()
class ApiMetadata {
  @JsonKey(name: 'request_id')
  final String? requestId;

  final String? timestamp;

  @JsonKey(name: 'processing_time_ms')
  final int? processingTimeMs;

  @JsonKey(name: 'llm_used')
  final String? llmUsed;

  @JsonKey(name: 'sources_searched')
  final List<String>? sourcesSearched;

  ApiMetadata({
    this.requestId,
    this.timestamp,
    this.processingTimeMs,
    this.llmUsed,
    this.sourcesSearched,
  });

  factory ApiMetadata.fromJson(Map<String, dynamic> json) =>
      _$ApiMetadataFromJson(json);

  Map<String, dynamic> toJson() => _$ApiMetadataToJson(this);
}
