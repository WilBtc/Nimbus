// src/core/security_manager.dart
import 'package:cryptography/cryptography.dart';
import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'logging_service.dart'; // Assuming a logging service for audit purposes

/// Manages security aspects such as key generation, signing, and verification.
/// Uses secure storage to store private keys.
class SecurityManager {
  final String atSign;
  late SimpleKeyPair keyPair;
  late String privateKey;
  late String publicKey;
  final FlutterSecureStorage secureStorage;
  final LoggingService loggingService = LoggingService(); // For logging events

  /// Creates a new instance of [SecurityManager] for the given [atSign].
  SecurityManager(this.atSign)
      : secureStorage = const FlutterSecureStorage();

  /// Initializes the security manager by loading or generating keys.
  Future<void> initialize() async {
    try {
      // Load or generate key pair and initialize security settings
      keyPair = await _loadOrGenerateKeyPair();

      // Extract and store private and public keys as base64 encoded strings
      privateKey = base64Encode(await keyPair.extractPrivateKeyBytes());
      publicKey = base64Encode((await keyPair.extractPublicKey()).bytes);

      // Store private key securely
      await secureStorage.write(key: '${atSign}_privateKey', value: privateKey);

      loggingService.log('SecurityManager initialized for atSign: $atSign');
    } catch (e) {
      loggingService.error('Error initializing SecurityManager for $atSign: $e');
      rethrow;
    }
  }

  /// Loads an existing key pair from secure storage or generates a new one.
  Future<SimpleKeyPair> _loadOrGenerateKeyPair() async {
    final algorithm = X25519();

    // Attempt to load existing private key from secure storage
    String? storedPrivateKey =
        await secureStorage.read(key: '${atSign}_privateKey');

    if (storedPrivateKey != null) {
      loggingService.log('Loading existing key pair from secure storage.');
      final keyBytes = base64Decode(storedPrivateKey);
      return algorithm.newKeyPairFromSeed(keyBytes);
    }

    // Generate a new key pair if none exist
    loggingService.log('Generating new key pair for atSign: $atSign');
    final newKeyPair = await algorithm.newKeyPair();
    final privateKeyBytes = await newKeyPair.extractPrivateKeyBytes();

    // Store the private key securely
    await secureStorage.write(
        key: '${atSign}_privateKey', value: base64Encode(privateKeyBytes));
    loggingService.log('Key pair generated and stored securely.');
    return newKeyPair;
  }

  /// Signs a message using the private key and returns the signature as a base64 string.
  Future<String> signMessage(String message) async {
    try {
      final signatureAlgorithm = Ed25519();
      final messageBytes = utf8.encode(message);
      final signature =
          await signatureAlgorithm.sign(messageBytes, keyPair: keyPair);
      loggingService.log('Message signed successfully for atSign: $atSign.');
      return base64Encode(signature.bytes);
    } catch (e) {
      loggingService.error('Error signing message for $atSign: $e');
      return '';
    }
  }

  /// Verifies a signature using the provided public key and returns the result as a boolean.
  Future<bool> verifySignature(
      String message, String signature, String publicKey) async {
    try {
      final signatureAlgorithm = Ed25519();
      final messageBytes = utf8.encode(message);
      final signatureBytes = base64Decode(signature);
      final publicKeyBytes = base64Decode(publicKey);

      final publicKeyObj =
          SimplePublicKey(publicKeyBytes, type: KeyPairType.ed25519);

      final isValid = await signatureAlgorithm.verify(
        messageBytes,
        signature: Signature(signatureBytes, publicKey: publicKeyObj),
      );

      if (isValid) {
        loggingService.log('Signature verification successful for atSign: $atSign.');
      } else {
        loggingService.warn('Signature verification failed for atSign: $atSign.');
      }
      return isValid;
    } catch (e) {
      loggingService.error('Error verifying signature for $atSign: $e');
      return false;
    }
  }

  /// Deletes the stored private key from secure storage.
  Future<void> deleteKeys() async {
    try {
      await secureStorage.delete(key: '${atSign}_privateKey');
      loggingService.log('Private key deleted for atSign: $atSign');
    } catch (e) {
      loggingService.error('Error deleting private key for $atSign: $e');
    }
  }

  /// Retrieves the stored private key for internal usage or future development.
  Future<String?> retrievePrivateKey() async {
    try {
      String? storedPrivateKey =
          await secureStorage.read(key: '${atSign}_privateKey');
      if (storedPrivateKey != null) {
        loggingService.log('Private key retrieved for atSign: $atSign');
      } else {
        loggingService.warn('No private key found for atSign: $atSign');
      }
      return storedPrivateKey;
    } catch (e) {
      loggingService.error('Error retrieving private key for $atSign: $e');
      return null;
    }
  }
}
