// src/core/config.dart
import 'dart:io';

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

  /// Loads the configuration from environment variables or defaults.
  Future<void> load() async {
    try {
      // Load configuration values from environment variables or use defaults
      atSign = _getEnvVar('ATSIGN', '@example');
      namespace = _getEnvVar('NAMESPACE', 'nimbus');
      rootDomain = _getEnvVar('ROOT_DOMAIN', 'root.atsign.org');
      hiveStoragePath = _getEnvVar('HIVE_STORAGE_PATH', './nimbus/storage');
      commitLogPath = _getEnvVar('COMMIT_LOG_PATH', './nimbus/commitLog');
      retryCount = _getEnvVarAsInt('RETRY_COUNT', 3);
      timeout = Duration(seconds: _getEnvVarAsInt('TIMEOUT_SECONDS', 10));
      enableDiagnostics = _getEnvVarAsBool('ENABLE_DIAGNOSTICS', false);

      // Validate essential configurations
      _validateConfiguration();

      print('Configuration loaded successfully.');
    } catch (e) {
      print('Error loading configuration: $e');
      rethrow;
    }
  }

  String _getEnvVar(String key, String defaultValue) {
    return Platform.environment[key] ?? defaultValue;
  }

  int _getEnvVarAsInt(String key, int defaultValue) {
    final value = Platform.environment[key];
    if (value != null) {
      return int.tryParse(value) ?? defaultValue;
    }
    return defaultValue;
  }

  bool _getEnvVarAsBool(String key, bool defaultValue) {
    final value = Platform.environment[key];
    if (value != null) {
      return value.toLowerCase() == 'true';
    }
    return defaultValue;
  }

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

    if (!Directory(hiveStoragePath).existsSync()) {
      Directory(hiveStoragePath)
          .createSync(recursive: true);
      print('Created Hive storage directory: $hiveStoragePath');
    }

    if (!Directory(commitLogPath).existsSync()) {
      Directory(commitLogPath)
          .createSync(recursive: true);
      print('Created commit log directory: $commitLogPath');
    }
  }
}
