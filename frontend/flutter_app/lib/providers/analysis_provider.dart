import 'package:flutter/foundation.dart';
import '../models/analysis_result.dart';
import '../services/api_service.dart';

/// Analysis state
enum AnalysisState {
  idle,
  loading,
  success,
  error,
}

/// Provider for managing analysis state
class AnalysisProvider with ChangeNotifier {
  final ApiService _apiService;

  AnalysisProvider({ApiService? apiService})
      : _apiService = apiService ?? ApiService();

  AnalysisState _state = AnalysisState.idle;
  AnalysisResult? _currentResult;
  String? _errorMessage;
  List<AnalysisResult> _history = [];

  // Getters
  AnalysisState get state => _state;
  AnalysisResult? get currentResult => _currentResult;
  String? get errorMessage => _errorMessage;
  List<AnalysisResult> get history => _history;
  bool get isLoading => _state == AnalysisState.loading;
  bool get hasError => _state == AnalysisState.error;

  /// Analyze text
  Future<void> analyzeText({
    required String text,
    bool includeEvidences = false,
    int? maxClaims,
    List<String>? searchEngines,
    bool? useLLM,
    String? llmProvider,
  }) async {
    if (text.trim().isEmpty) {
      _state = AnalysisState.error;
      _errorMessage = 'テキストを入力してください';
      notifyListeners();
      return;
    }

    _state = AnalysisState.loading;
    _errorMessage = null;
    notifyListeners();

    try {
      final result = await _apiService.analyzeText(
        text: text,
        includeEvidences: includeEvidences,
        maxClaims: maxClaims,
        searchEngines: searchEngines,
        useLLM: useLLM,
        llmProvider: llmProvider,
      );

      _currentResult = result;
      _state = AnalysisState.success;
      _errorMessage = null;

      // Add to history (local cache)
      _history.insert(0, result);
      if (_history.length > 50) {
        _history = _history.sublist(0, 50);
      }
    } catch (e) {
      _state = AnalysisState.error;
      _errorMessage = e.toString();
      _currentResult = null;
    } finally {
      notifyListeners();
    }
  }

  /// Load history from API
  Future<void> loadHistory({int limit = 20, int offset = 0}) async {
    try {
      final results = await _apiService.getHistory(
        limit: limit,
        offset: offset,
      );
      _history = results;
      notifyListeners();
    } catch (e) {
      if (kDebugMode) {
        print('Failed to load history: $e');
      }
    }
  }

  /// Clear current result
  void clearResult() {
    _currentResult = null;
    _state = AnalysisState.idle;
    _errorMessage = null;
    notifyListeners();
  }

  /// Reset all state
  void reset() {
    _state = AnalysisState.idle;
    _currentResult = null;
    _errorMessage = null;
    _history = [];
    notifyListeners();
  }

  @override
  void dispose() {
    _apiService.dispose();
    super.dispose();
  }
}
