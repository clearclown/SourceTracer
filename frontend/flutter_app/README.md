# SourceTracer Flutter Web Application

## Overview

This is the Flutter Web frontend for SourceTracer - a citation and evidence analysis system.

## Features

- 📝 **Text Analysis**: Analyze text to extract claims and classify as opinion or fact
- 🔍 **Evidence Search**: Search multiple academic databases and sources for supporting evidence
- 📊 **Credibility Scoring**: Automatic credibility assessment based on multiple factors
- 📜 **Analysis History**: View and manage past analyses
- 🎨 **Modern UI**: Clean, responsive Material Design 3 interface

## Project Structure

```
flutter_app/
├── lib/
│   ├── config/              # Environment configuration
│   │   └── environment.dart
│   ├── models/              # Data models
│   │   ├── claim.dart
│   │   ├── evidence.dart
│   │   └── analysis_result.dart
│   ├── services/            # API services
│   │   └── api_service.dart
│   ├── providers/           # State management (Provider)
│   │   └── analysis_provider.dart
│   ├── screens/             # UI screens
│   │   ├── home_screen.dart
│   │   ├── analysis_screen.dart
│   │   └── history_screen.dart
│   ├── widgets/             # Reusable widgets
│   │   ├── claim_card.dart
│   │   └── evidence_list.dart
│   └── main.dart            # App entry point
├── web/                     # Web-specific files
│   ├── index.html
│   └── manifest.json
├── test/                    # Tests
├── pubspec.yaml             # Dependencies
└── analysis_options.yaml    # Linter configuration
```

## Prerequisites

- Flutter SDK 3.0.0 or higher
- Chrome/Edge browser for development

## Installation

### 1. Install Flutter SDK

Follow the official Flutter installation guide:
https://docs.flutter.dev/get-started/install

### 2. Clone Repository

```bash
cd /home/ablaze/research/SourceTracer/frontend/flutter_app
```

### 3. Install Dependencies

```bash
flutter pub get
```

### 4. Generate Code

Generate JSON serialization code:

```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

## Development

### Run in Development Mode

```bash
# Default (Chrome)
flutter run -d chrome

# With custom API URL
flutter run -d chrome --dart-define=API_BASE_URL=http://localhost:8080/api/v1

# Enable debug mode
flutter run -d chrome --dart-define=DEBUG_MODE=true
```

### Hot Reload

Press `r` in the terminal to hot reload changes during development.

## Building

### Build for Production

```bash
flutter build web --release
```

The output will be in `build/web/`.

### Build with Custom Configuration

```bash
flutter build web --release \
  --dart-define=API_BASE_URL=https://api.sourcetracer.org/api/v1 \
  --dart-define=DEBUG_MODE=false
```

### Serve Production Build Locally

```bash
cd build/web
python3 -m http.server 8000
```

Then open http://localhost:8000 in your browser.

## Testing

### Run Unit Tests

```bash
flutter test
```

### Run Tests with Coverage

```bash
flutter test --coverage
```

### Analyze Code

```bash
flutter analyze
```

## Docker Deployment

### Dockerfile

Create a `Dockerfile` in `frontend/flutter_app/`:

```dockerfile
# Build stage
FROM ubuntu:22.04 AS builder

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    && rm -rf /var/lib/apt/lists/*

# Install Flutter
RUN git clone https://github.com/flutter/flutter.git /flutter
ENV PATH="/flutter/bin:${PATH}"

WORKDIR /app
COPY . .

# Get dependencies and build
RUN flutter config --enable-web
RUN flutter pub get
RUN flutter pub run build_runner build --delete-conflicting-outputs
RUN flutter build web --release \
    --dart-define=API_BASE_URL=http://api:8080/api/v1

# Runtime stage
FROM nginx:alpine

# Copy built app to nginx
COPY --from=builder /app/build/web /usr/share/nginx/html

# Custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### nginx.conf

Create `nginx.conf`:

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # Enable gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### Build Docker Image

```bash
docker build -t sourcetracer-web .
```

### Run Container

```bash
docker run -d -p 3000:80 sourcetracer-web
```

## Environment Variables

Configure the app using `--dart-define` flags:

| Variable | Default | Description |
|----------|---------|-------------|
| `API_BASE_URL` | `http://localhost:8080/api/v1` | Backend API URL |
| `DEBUG_MODE` | `true` | Enable debug logging |
| `REQUEST_TIMEOUT` | `30` | API request timeout (seconds) |

## Integration with Backend

The frontend communicates with the Go backend API. Make sure the backend is running:

```bash
cd ../../backend
go run cmd/api/main.go
```

Or use Podman Compose to run the full stack:

```bash
cd ../..
podman compose up -d
```

## API Integration

### Analyze Text

```dart
final provider = context.read<AnalysisProvider>();

await provider.analyzeText(
  text: "Your text here",
  includeEvidences: true,
  useLLM: true,
  llmProvider: "claude",
);

// Access result
final result = provider.currentResult;
```

### Load History

```dart
await provider.loadHistory(limit: 20, offset: 0);
final history = provider.history;
```

## Troubleshooting

### CORS Issues

If you encounter CORS errors, ensure the backend API has CORS enabled:

```go
// In backend/cmd/api/main.go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
    AllowCredentials: true,
}))
```

### Build Runner Issues

If code generation fails:

```bash
# Clean and regenerate
flutter clean
flutter pub get
flutter pub run build_runner clean
flutter pub run build_runner build --delete-conflicting-outputs
```

### Web Renderer Issues

Try different renderers:

```bash
# HTML renderer (better compatibility)
flutter run -d chrome --web-renderer html

# CanvasKit renderer (better performance)
flutter run -d chrome --web-renderer canvaskit
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linter
5. Submit a pull request

## License

MIT License - see LICENSE file

## Support

- Issue Tracker: https://github.com/yourusername/sourcetracer/issues
- Documentation: https://github.com/yourusername/sourcetracer/tree/main/docs

---

**SourceTracer** - 情報源を、徹底的に追う。
