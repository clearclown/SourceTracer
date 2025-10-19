# SourceTracer Frontend (Flutter Web)

## Overview

This directory will contain the Flutter Web frontend for SourceTracer.

## Setup

### Prerequisites
- Flutter SDK 3.16+
- Chrome/Edge for web development

### Create Flutter Project
```bash
cd frontend
flutter create flutter_app
cd flutter_app
flutter pub add http provider
```

### Project Structure
```
flutter_app/
├── lib/
│   ├── main.dart              # App entry point
│   ├── models/                # Data models
│   │   ├── claim.dart
│   │   ├── evidence.dart
│   │   └── analysis_result.dart
│   ├── services/              # API services
│   │   └── api_service.dart
│   ├── providers/             # State management
│   │   └── analysis_provider.dart
│   ├── screens/               # UI screens
│   │   ├── home_screen.dart
│   │   ├── analysis_screen.dart
│   │   └── history_screen.dart
│   └── widgets/               # Reusable widgets
│       ├── claim_card.dart
│       ├── evidence_list.dart
│       └── text_input_form.dart
├── web/
│   └── index.html
└── pubspec.yaml
```

## Features to Implement

### Phase 1: Basic UI (MVP)
- [ ] Text input form for analysis
- [ ] Submit button with loading state
- [ ] Display analysis results (claims with types)
- [ ] Basic styling with Material Design

### Phase 2: Enhanced Features
- [ ] Confidence score visualization
- [ ] Evidence display with expandable cards
- [ ] Credibility scoring visualization
- [ ] History page with past analyses

### Phase 3: Advanced Features
- [ ] Search engine selection
- [ ] Real-time analysis progress
- [ ] Export results (JSON, PDF)
- [ ] Dark mode support

## API Integration

### Example API Service

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  final String baseUrl = 'http://localhost:8080/api/v1';

  Future<AnalysisResult> analyzeText(String text) async {
    final response = await http.post(
      Uri.parse('$baseUrl/analyze'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'text': text,
        'options': {
          'include_evidences': true,
          'max_claims': 10,
        },
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return AnalysisResult.fromJson(data['data']);
    } else {
      throw Exception('Failed to analyze text');
    }
  }
}
```

## Development

```bash
# Run in development mode
flutter run -d chrome

# Build for production
flutter build web

# Serve production build
cd build/web
python3 -m http.server 8000
```

## Docker Support

```dockerfile
# Dockerfile for Flutter Web
FROM ubuntu:22.04 AS builder

RUN apt-get update && apt-get install -y curl git wget unzip
RUN git clone https://github.com/flutter/flutter.git /flutter
ENV PATH="/flutter/bin:${PATH}"

WORKDIR /app
COPY . .
RUN flutter config --enable-web
RUN flutter pub get
RUN flutter build web --release

FROM nginx:alpine
COPY --from=builder /app/build/web /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## Environment Configuration

```dart
// lib/config/environment.dart
class Environment {
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080/api/v1',
  );
}
```

## State Management (Provider Pattern)

```dart
import 'package:flutter/foundation.dart';
import '../models/analysis_result.dart';
import '../services/api_service.dart';

class AnalysisProvider with ChangeNotifier {
  final ApiService _apiService = ApiService();

  AnalysisResult? _currentResult;
  bool _isLoading = false;
  String? _error;

  AnalysisResult? get currentResult => _currentResult;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> analyzeText(String text) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _currentResult = await _apiService.analyzeText(text);
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
```

## UI Examples

### Home Screen
```dart
class HomeScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('SourceTracer'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('情報源を、徹底的に追う。'),
            ElevatedButton(
              onPressed: () => Navigator.pushNamed(context, '/analyze'),
              child: Text('Start Analysis'),
            ),
          ],
        ),
      ),
    );
  }
}
```

## Testing

```bash
# Run unit tests
flutter test

# Run integration tests
flutter test integration_test
```

## Deployment

The frontend will be served from the `frontend` service in podman-compose.yml once implemented.

## TODO
- [ ] Initialize Flutter project
- [ ] Implement API service layer
- [ ] Create data models
- [ ] Build text input UI
- [ ] Display analysis results
- [ ] Add history view
- [ ] Implement responsive design
- [ ] Add loading states
- [ ] Error handling
- [ ] Unit tests
