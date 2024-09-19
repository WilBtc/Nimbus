// src/core/security_manager.dart
import 'package:cryptography/cryptography.dart';
import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Manages security aspects such as key generation, signing, and verification.
/// Uses secure storage to store private keys.
class SecurityManager {
  final String atSign;
  late SimpleKeyPair keyPair;
  late String privateKey;
  late String publicKey;
  final FlutterSecureStorage secureStorage;

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

      print('SecurityManager initialized for atSign: $atSign');
    } catch (e) {
      print('Error initializing SecurityManager: $e');
      rethrow;
    }
  }

  Future<SimpleKeyPair> _loadOrGenerateKeyPair() async {
    final algorithm = X25519();

    // Attempt to load existing private key from secure storage
    String? storedPrivateKey =
        await secureStorage.read(key: '${atSign}_privateKey');

    if (storedPrivateKey != null) {
      print('Loading existing key pair from secure storage.');
      final keyBytes = base64Decode(storedPrivateKey);
      return algorithm.newKeyPairFromSeed(keyBytes);
    }

    // Generate a new key pair if none exist
    print('Generating new key pair for atSign: $atSign');
    final newKeyPair = await algorithm.newKeyPair();
    final privateKeyBytes = await newKeyPair.extractPrivateKeyBytes();

    // Store the private key securely
    await secureStorage.write(
        key: '${atSign}_privateKey', value: base64Encode(privateKeyBytes));
    print('Key pair generated and stored securely.');
    return newKeyPair;
  }

  /// Signs a message using the private key.
  Future<String> signMessage(String message) async {
    try {
      final signatureAlgorithm = Ed25519();
      final messageBytes = utf8.encode(message);
      final signature =
          await signatureAlgorithm.sign(messageBytes, keyPair: keyPair);
      print('Message signed successfully.');
      return base64Encode(signature.bytes);
    } catch (e) {
      print('Error signing message: $e');
      return '';
    }
  }

  /// Verifies a signature using the provided public key.
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
        print('Signature verification successful.');
      } else {
        print('Signature verification failed.');
      }
      return isValid;
    } catch (e) {
      print('Error verifying signature: $e');
      return false;
    }
  }

  /// Deletes the stored private key from secure storage.
  Future<void> deleteKeys() async {
    try {
      await secureStorage.delete(key: '${atSign}_privateKey');
      print('Private key deleted for atSign: $atSign');
    } catch (e) {
      print('Error deleting private key: $e');
    }
  }
}
