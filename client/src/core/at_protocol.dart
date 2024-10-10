// src/core/at_protocol.dart
import 'package:at_client_mobile/at_client_mobile.dart';
import 'security_manager.dart';
import 'config.dart';
import 'logging_service.dart'; // Added logging for audit purposes
import 'encryption_service.dart'; // Added encryption module for message security

/// Manages the communication over the atProtocol with security and encryption.
class AtProtocol {
  final SecurityManager securityManager;
  final LoggingService loggingService = LoggingService(); // Logging Service for audit compliance
  final EncryptionService encryptionService = EncryptionService(); // Encryption service for message encryption

  late AtClient atClient;
  late AtClientPreference atClientPreference;

  /// Creates a new instance of [AtProtocol] with the given [securityManager].
  AtProtocol(this.securityManager);

  /// Initializes the atProtocol client.
  Future<void> initialize() async {
    try {
      // Load configuration
      final config = Config();
      await config.load();

      // Setup AtClient preferences
      atClientPreference = AtClientPreference()
        ..namespace = config.namespace
        ..rootDomain = config.rootDomain
        ..hiveStoragePath = config.hiveStoragePath
        ..commitLogPath = config.commitLogPath
        ..isLocalStoreRequired = true;

      // Initialize the AtClient for communication
      atClient = await AtClientImpl.createClient(
          config.atSign, config.namespace, atClientPreference);

      var localSecondary = atClient.getLocalSecondary();

      if (localSecondary == null) {
        throw Exception('Local secondary is null for atSign: ${config.atSign}');
      }

      var connectionResult = await localSecondary.connect();
      if (!connectionResult) {
        throw Exception(
            'Failed to connect to atProtocol with atSign: ${config.atSign}');
      }

      loggingService.log('Connection established with atSign: ${config.atSign}');

    } catch (e) {
      loggingService.error('Error initializing atProtocol: $e');
      rethrow;
    }
  }

  /// Sends an encrypted message with the given [key] and [message].
  Future<void> sendMessage(String key, String message) async {
    var encryptedMessage = encryptionService.encrypt(message); // Encrypt message

    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.put(atKey, encryptedMessage); // Store encrypted message
      if (result) {
        loggingService.log('Message sent successfully with key: $key');
      } else {
        loggingService.error('Failed to send the message with key: $key');
      }
    } catch (e) {
      loggingService.error('Error sending message: $e');
    }
  }

  /// Retrieves and decrypts a message with the given [key].
  Future<String?> retrieveMessage(String key) async {
    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.get(atKey);
      if (result.value != null) {
        var decryptedMessage = encryptionService.decrypt(result.value!); // Decrypt message
        loggingService.log('Message retrieved successfully with key: $key');
        return decryptedMessage;
      } else {
        loggingService.warn('No message found for key: $key');
      }
      return null;
    } catch (e) {
      loggingService.error('Error retrieving message: $e');
      return null;
    }
  }

  /// Deletes a message with the given [key] and logs the deletion event.
  Future<void> deleteMessage(String key) async {
    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.delete(atKey);
      if (result) {
        loggingService.log('Message deleted successfully with key: $key');
      } else {
        loggingService.error('Failed to delete message with key: $key');
      }
    } catch (e) {
      loggingService.error('Error deleting message: $e');
    }
  }

  /// Closes the atProtocol client connection and logs the event.
  void close() {
    try {
      atClient.getLocalSecondary()?.close();
      loggingService.log('Connection closed for atSign: ${securityManager.atSign}');
    } catch (e) {
      loggingService.error('Error closing connection: $e');
    }
  }
}
