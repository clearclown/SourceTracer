import 'dart:convert';
import 'package:http/http.dart' as http;
import '../config/environment.dart';
import '../models/analysis_result.dart';

/// Exception thrown when API request fails
class ApiException implements Exception {
  final String message;
  final int? statusCode;
  final ApiError? error;

  ApiException(this.message, {this.statusCode, this.error});

  @override
  String toString() => 'ApiException: $message (status: $statusCode)';
}

/// API Service for SourceTracer backend
class ApiService {
  final String baseUrl;
  final http.Client client;

  ApiService({
    String? baseUrl,
    http.Client? client,
  })  : baseUrl = baseUrl ?? Environment.apiBaseUrl,
        client = client ?? http.Client();

  /// Analyze text and get claims with evidences
  Future<AnalysisResult> analyzeText({
    required String text,
    bool includeEvidences = false,
    int? maxClaims,
    int? maxEvidences,
    List<String>? searchEngines,
    bool? useLLM,
    String? llmProvider,
  }) async {
    final url = Uri.parse('$baseUrl/analyze');

    final requestBody = {
      'text': text,
      'options': {
        'include_evidences': includeEvidences,
        if (maxClaims != null) 'max_claims': maxClaims,
        if (maxEvidences != null) 'max_evidences': maxEvidences,
        if (searchEngines != null && searchEngines.isNotEmpty)
          'search_engines': searchEngines,
        if (useLLM != null) 'use_llm': useLLM,
        if (llmProvider != null) 'llm_provider': llmProvider,
      },
    };

    try {
      final response = await client
          .post(
            url,
            headers: {'Content-Type': 'application/json'},
            body: jsonEncode(requestBody),
          )
          .timeout(Duration(seconds: Environment.requestTimeout));

      final jsonData = jsonDecode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          jsonData,
          (json) => json as Map<String, dynamic>,
        );

        if (apiResponse.success && apiResponse.data != null) {
          return AnalysisResult.fromJson(apiResponse.data!);
        } else {
          throw ApiException(
            apiResponse.error?.message ?? 'Unknown error',
            statusCode: response.statusCode,
            error: apiResponse.error,
          );
        }
      } else {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          jsonData,
          (json) => json as Map<String, dynamic>,
        );
        throw ApiException(
          apiResponse.error?.message ?? 'Request failed',
          statusCode: response.statusCode,
          error: apiResponse.error,
        );
      }
    } catch (e) {
      if (e is ApiException) rethrow;
      throw ApiException('Network error: ${e.toString()}');
    }
  }

  /// Get analysis history
  Future<List<AnalysisResult>> getHistory({
    int limit = 20,
    int offset = 0,
  }) async {
    final url = Uri.parse('$baseUrl/history?limit=$limit&offset=$offset');

    try {
      final response = await client
          .get(url)
          .timeout(Duration(seconds: Environment.requestTimeout));

      final jsonData = jsonDecode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          jsonData,
          (json) => json as Map<String, dynamic>,
        );

        if (apiResponse.success && apiResponse.data != null) {
          final analyses = apiResponse.data!['analyses'] as List;
          return analyses
              .map((a) => AnalysisResult.fromJson(a as Map<String, dynamic>))
              .toList();
        }
      }

      throw ApiException(
        'Failed to fetch history',
        statusCode: response.statusCode,
      );
    } catch (e) {
      if (e is ApiException) rethrow;
      throw ApiException('Network error: ${e.toString()}');
    }
  }

  /// Get specific analysis by ID
  Future<AnalysisResult> getAnalysis(String id) async {
    final url = Uri.parse('$baseUrl/history/$id');

    try {
      final response = await client
          .get(url)
          .timeout(Duration(seconds: Environment.requestTimeout));

      final jsonData = jsonDecode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          jsonData,
          (json) => json as Map<String, dynamic>,
        );

        if (apiResponse.success && apiResponse.data != null) {
          return AnalysisResult.fromJson(apiResponse.data!);
        }
      }

      throw ApiException(
        'Analysis not found',
        statusCode: response.statusCode,
      );
    } catch (e) {
      if (e is ApiException) rethrow;
      throw ApiException('Network error: ${e.toString()}');
    }
  }

  /// Health check
  Future<bool> healthCheck() async {
    final url = Uri.parse('$baseUrl/health');

    try {
      final response = await client
          .get(url)
          .timeout(const Duration(seconds: 5));

      if (response.statusCode == 200) {
        final jsonData = jsonDecode(response.body) as Map<String, dynamic>;
        return jsonData['status'] == 'ok';
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  /// Dispose resources
  void dispose() {
    client.close();
  }
}
