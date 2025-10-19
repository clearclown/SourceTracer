/// Environment configuration for SourceTracer
class Environment {
  /// API Base URL
  /// Can be overridden using --dart-define=API_BASE_URL=<url>
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080/api/v1',
  );

  /// Enable debug logging
  static const bool debugMode = bool.fromEnvironment(
    'DEBUG_MODE',
    defaultValue: true,
  );

  /// Request timeout in seconds
  static const int requestTimeout = int.fromEnvironment(
    'REQUEST_TIMEOUT',
    defaultValue: 30,
  );

  /// Get production URL
  static String get productionUrl => 'https://api.sourcetracer.org/api/v1';

  /// Check if running in production
  static bool get isProduction =>
      apiBaseUrl.contains('sourcetracer.org');
}
