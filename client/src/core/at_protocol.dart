// src/core/at_protocol.dart
import 'package:at_client_mobile/at_client_mobile.dart';
import 'security_manager.dart';
import 'config.dart';

/// Manages the communication over the atProtocol.
class AtProtocol {
  final SecurityManager securityManager;
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
      print(
          'Successfully connected to atProtocol with atSign: ${config.atSign}');
    } catch (e) {
      print('Error initializing atProtocol: $e');
      rethrow;
    }
  }

  /// Sends a message with the given [key] and [message].
  Future<void> sendMessage(String key, String message) async {
    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.put(atKey, message);
      if (result) {
        print('Message sent successfully with key: $key');
      } else {
        print('Failed to send the message with key: $key');
      }
    } catch (e) {
      print('Error sending message: $e');
    }
  }

  /// Retrieves a message with the given [key].
  Future<String?> retrieveMessage(String key) async {
    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.get(atKey);
      if (result.value != null) {
        print('Message retrieved successfully with key: $key');
      } else {
        print('No message found for key: $key');
      }
      return result.value;
    } catch (e) {
      print('Error retrieving message: $e');
      return null;
    }
  }

  /// Deletes a message with the given [key].
  Future<void> deleteMessage(String key) async {
    var atKey = AtKey()
      ..key = key
      ..sharedWith = securityManager.atSign;

    try {
      var result = await atClient.delete(atKey);
      if (result) {
        print('Message deleted successfully with key: $key');
      } else {
        print('Failed to delete message with key: $key');
      }
    } catch (e) {
      print('Error deleting message: $e');
    }
  }

  /// Closes the atProtocol client connection.
  void close() {
    try {
      atClient.getLocalSecondary()?.close();
      print('Connection closed for atSign: ${securityManager.atSign}');
    } catch (e) {
      print('Error closing connection: $e');
    }
  }
}
