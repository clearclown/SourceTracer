// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'evidence.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Evidence _$EvidenceFromJson(Map<String, dynamic> json) => Evidence(
      id: json['id'] as String,
      source: json['source'] as String,
      title: json['title'] as String,
      authors:
          (json['authors'] as List<dynamic>).map((e) => e as String).toList(),
      url: json['url'] as String,
      snippet: json['snippet'] as String,
      credibility: (json['credibility'] as num).toDouble(),
      relevance: (json['relevance'] as num).toDouble(),
      publishedDate: json['published_date'] as String?,
      sourceType: json['source_type'] as String,
      citationCount: (json['citation_count'] as num?)?.toInt(),
      metadata: json['metadata'] as Map<String, dynamic>?,
    );

Map<String, dynamic> _$EvidenceToJson(Evidence instance) => <String, dynamic>{
      'id': instance.id,
      'source': instance.source,
      'title': instance.title,
      'authors': instance.authors,
      'url': instance.url,
      'snippet': instance.snippet,
      'credibility': instance.credibility,
      'relevance': instance.relevance,
      'published_date': instance.publishedDate,
      'source_type': instance.sourceType,
      'citation_count': instance.citationCount,
      'metadata': instance.metadata,
    };
