// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'analysis_result.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AnalysisResult _$AnalysisResultFromJson(Map<String, dynamic> json) =>
    AnalysisResult(
      analysisId: json['analysis_id'] as String,
      claims: (json['claims'] as List<dynamic>)
          .map((e) => Claim.fromJson(e as Map<String, dynamic>))
          .toList(),
      overallCredibility: (json['overall_credibility'] as num).toDouble(),
      processingTimeMs: (json['processing_time_ms'] as num).toInt(),
      createdAt: json['created_at'] as String?,
    );

Map<String, dynamic> _$AnalysisResultToJson(AnalysisResult instance) =>
    <String, dynamic>{
      'analysis_id': instance.analysisId,
      'claims': instance.claims.map((e) => e.toJson()).toList(),
      'overall_credibility': instance.overallCredibility,
      'processing_time_ms': instance.processingTimeMs,
      'created_at': instance.createdAt,
    };

ApiResponse<T> _$ApiResponseFromJson<T>(
  Map<String, dynamic> json,
  T Function(Object? json) fromJsonT,
) =>
    ApiResponse<T>(
      success: json['success'] as bool,
      data: _$nullableGenericFromJson(json['data'], fromJsonT),
      error: json['error'] == null
          ? null
          : ApiError.fromJson(json['error'] as Map<String, dynamic>),
      metadata: json['metadata'] == null
          ? null
          : ApiMetadata.fromJson(json['metadata'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$ApiResponseToJson<T>(
  ApiResponse<T> instance,
  Object? Function(T value) toJsonT,
) =>
    <String, dynamic>{
      'success': instance.success,
      'data': _$nullableGenericToJson(instance.data, toJsonT),
      'error': instance.error,
      'metadata': instance.metadata,
    };

T? _$nullableGenericFromJson<T>(
  Object? input,
  T Function(Object? json) fromJson,
) =>
    input == null ? null : fromJson(input);

Object? _$nullableGenericToJson<T>(
  T? input,
  Object? Function(T value) toJson,
) =>
    input == null ? null : toJson(input);

ApiError _$ApiErrorFromJson(Map<String, dynamic> json) => ApiError(
      code: json['code'] as String,
      message: json['message'] as String,
      details: json['details'] as Map<String, dynamic>?,
    );

Map<String, dynamic> _$ApiErrorToJson(ApiError instance) => <String, dynamic>{
      'code': instance.code,
      'message': instance.message,
      'details': instance.details,
    };

ApiMetadata _$ApiMetadataFromJson(Map<String, dynamic> json) => ApiMetadata(
      requestId: json['request_id'] as String?,
      timestamp: json['timestamp'] as String?,
      processingTimeMs: (json['processing_time_ms'] as num?)?.toInt(),
      llmUsed: json['llm_used'] as String?,
      sourcesSearched: (json['sources_searched'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
    );

Map<String, dynamic> _$ApiMetadataToJson(ApiMetadata instance) =>
    <String, dynamic>{
      'request_id': instance.requestId,
      'timestamp': instance.timestamp,
      'processing_time_ms': instance.processingTimeMs,
      'llm_used': instance.llmUsed,
      'sources_searched': instance.sourcesSearched,
    };
