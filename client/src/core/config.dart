// src/core/config.dart
import 'dart:io';
import 'logging_service.dart'; // Assuming a logging service for audit logs

/// Manages the configuration settings for the application.
class Config {
  late String atSign;
  late String namespace;
  late String rootDomain;
  late String hiveStoragePath;
  late String commitLogPath;
  late int retryCount;
  late Duration timeout;
  late bool enableDiagnostics;

  final LoggingService loggingService = LoggingService(); // Logging for audit purposes

  /// Loads the configuration from environment variables or defaults.
  Future<void> load() async {
    try {
      // Load configuration values from environment variables or use defaults
      atSign = _sanitizeEnvVar('ATSIGN', '@example');
      namespace = _sanitizeEnvVar('NAMESPACE', 'nimbus');
      rootDomain = _sanitizeEnvVar('ROOT_DOMAIN', 'root.atsign.org');
      hiveStoragePath = _sanitizeEnvVar('HIVE_STORAGE_PATH', './nimbus/storage');
      commitLogPath = _sanitizeEnvVar('COMMIT_LOG_PATH', './nimbus/commitLog');
      retryCount = _getEnvVarAsInt('RETRY_COUNT', 3);
      timeout = Duration(seconds: _getEnvVarAsInt('TIMEOUT_SECONDS', 10));
      enableDiagnostics = _getEnvVarAsBool('ENABLE_DIAGNOSTICS', false);

      // Validate essential configurations
      _validateConfiguration();

      loggingService.log('Configuration loaded successfully.');
    } catch (e) {
      loggingService.error('Error loading configuration: $e');
      rethrow;
    }
  }

  /// Retrieves and sanitizes environment variables to prevent injection or formatting issues.
  String _sanitizeEnvVar(String key, String defaultValue) {
    final value = Platform.environment[key] ?? defaultValue;
    if (_isMaliciousInput(value)) {
      throw ArgumentError('Environment variable $key contains invalid characters.');
    }
    return value.trim();
  }

  /// Helper method to check for potentially malicious inputs.
  bool _isMaliciousInput(String value) {
    // A basic check for any unwanted special characters or patterns.
    return value.contains(RegExp(r'[<>;&|]'));
  }

  /// Retrieves an environment variable as an integer with a default fallback.
  int _getEnvVarAsInt(String key, int defaultValue) {
    final value = Platform.environment[key];
    if (value != null) {
      return int.tryParse(value) ?? defaultValue;
    }
    return defaultValue;
  }

  /// Retrieves an environment variable as a boolean with a default fallback.
  bool _getEnvVarAsBool(String key, bool defaultValue) {
    final value = Platform.environment[key];
    if (value != null) {
      return value.toLowerCase() == 'true';
    }
    return defaultValue;
  }

  /// Validates the essential configuration parameters and creates directories if needed.
  void _validateConfiguration() {
    if (atSign.isEmpty || !atSign.startsWith('@')) {
      throw ArgumentError('Invalid atSign configuration: $atSign');
    }

    if (namespace.isEmpty) {
      throw ArgumentError('Namespace must not be empty');
    }

    if (rootDomain.isEmpty) {
      throw ArgumentError('Root domain must not be empty');
    }

    _createDirectoryIfNotExists(hiveStoragePath, 'Hive storage');
    _createDirectoryIfNotExists(commitLogPath, 'Commit log');
  }

  /// Ensures directories exist, creating them if necessary, and logs the action.
  void _createDirectoryIfNotExists(String path, String type) {
    final directory = Directory(path);
    if (!directory.existsSync()) {
      directory.createSync(recursive: true);
      loggingService.log('Created $type directory: $path');
    }
  }
}
