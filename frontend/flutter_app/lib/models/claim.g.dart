// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'claim.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Claim _$ClaimFromJson(Map<String, dynamic> json) => Claim(
      id: json['id'] as String,
      text: json['text'] as String,
      type: json['type'] as String,
      confidence: (json['confidence'] as num).toDouble(),
      position: json['position'] == null
          ? null
          : Position.fromJson(json['position'] as Map<String, dynamic>),
      evidences: (json['evidences'] as List<dynamic>?)
              ?.map((e) => Evidence.fromJson(e as Map<String, dynamic>))
              .toList() ??
          const [],
      counterEvidences: (json['counter_evidences'] as List<dynamic>?)
          ?.map((e) => Evidence.fromJson(e as Map<String, dynamic>))
          .toList(),
      summary: json['summary'] as String?,
      recommendation: json['recommendation'] as String?,
    );

Map<String, dynamic> _$ClaimToJson(Claim instance) => <String, dynamic>{
      'id': instance.id,
      'text': instance.text,
      'type': instance.type,
      'confidence': instance.confidence,
      'position': instance.position?.toJson(),
      'evidences': instance.evidences.map((e) => e.toJson()).toList(),
      'counter_evidences':
          instance.counterEvidences?.map((e) => e.toJson()).toList(),
      'summary': instance.summary,
      'recommendation': instance.recommendation,
    };

Position _$PositionFromJson(Map<String, dynamic> json) => Position(
      start: (json['start'] as num).toInt(),
      end: (json['end'] as num).toInt(),
    );

Map<String, dynamic> _$PositionToJson(Position instance) => <String, dynamic>{
      'start': instance.start,
      'end': instance.end,
    };
