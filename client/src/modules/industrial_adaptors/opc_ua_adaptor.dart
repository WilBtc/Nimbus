import 'package:opcua/opcua.dart'; // Example OPC-UA client library, adjust based on actual library used.
import 'logging_service.dart'; // Assuming a logging service for centralized log handling.

class OPCUAAdaptor {
  final String serverUrl;
  final String securityPolicy;
  final String securityMode;
  final String username;
  final String password;

  late OpcUaClient _client;
  final LoggingService loggingService = LoggingService(); // Logging Service for centralized logging

  OPCUAAdaptor({
    required this.serverUrl,
    this.securityPolicy = 'Basic256Sha256', // Default security policy
    this.securityMode = 'SignAndEncrypt',  // Default security mode
    this.username = '',
    this.password = '',
  });

  /// Initialize and connect the OPC-UA client with security settings.
  Future<void> initialize() async {
    try {
      _client = OpcUaClient(
        endpointUrl: serverUrl,
        securityPolicy: securityPolicy,
        securityMode: securityMode,
      );

      // Configure user authentication if credentials are provided.
      if (username.isNotEmpty && password.isNotEmpty) {
        _client.setUserToken(username, password);
      }

      await _client.connect();
      loggingService.log('OPC-UA client connected to $serverUrl with security: $securityPolicy/$securityMode');
    } catch (e) {
      loggingService.error('Error connecting OPC-UA client: $e');
      rethrow;
    }
  }

  /// Read a specific node's value.
  Future<dynamic> readNode(String nodeId) async {
    try {
      var result = await _client.readNode(nodeId);
      loggingService.log('Read OPC-UA node $nodeId: $result');
      return result;
    } catch (e) {
      loggingService.error('Error reading OPC-UA node $nodeId: $e');
      return null;
    }
  }

  /// Write a value to a specific node.
  Future<void> writeNode(String nodeId, dynamic value) async {
    try {
      await _client.writeNode(nodeId, value);
      loggingService.log('Wrote $value to OPC-UA node $nodeId');
    } catch (e) {
      loggingService.error('Error writing to OPC-UA node $nodeId: $e');
    }
  }

  /// Read multiple nodes at once.
  Future<Map<String, dynamic>> readMultipleNodes(List<String> nodeIds) async {
    Map<String, dynamic> results = {};
    try {
      for (String nodeId in nodeIds) {
        var result = await readNode(nodeId);
        if (result != null) {
          results[nodeId] = result;
        }
      }
      loggingService.log('Read multiple OPC-UA nodes: $results');
      return results;
    } catch (e) {
      loggingService.error('Error reading multiple OPC-UA nodes: $e');
      return {};
    }
  }

  /// Subscribe to changes in a specific node.
  Future<void> subscribeToNode(String nodeId, Function(dynamic) onDataChange) async {
    try {
      await _client.subscribeToNode(nodeId, (value) {
        loggingService.log('Data change detected on node $nodeId: $value');
        onDataChange(value);
      });
    } catch (e) {
      loggingService.error('Error subscribing to OPC-UA node $nodeId: $e');
    }
  }

  /// Unsubscribe from a specific node.
  Future<void> unsubscribeFromNode(String nodeId) async {
    try {
      await _client.unsubscribeFromNode(nodeId);
      loggingService.log('Unsubscribed from OPC-UA node $nodeId');
    } catch (e) {
      loggingService.error('Error unsubscribing from OPC-UA node $nodeId: $e');
    }
  }

  /// Disconnect the OPC-UA client safely.
  void disconnect() {
    try {
      if (_client.isConnected) {
        _client.disconnect();
        loggingService.log('OPC-UA client disconnected.');
      } else {
        loggingService.log('OPC-UA client was already disconnected.');
      }
    } catch (e) {
      loggingService.error('Error during OPC-UA client disconnection: $e');
    }
  }

  /// Check if the client is currently connected.
  bool isConnected() {
    try {
      bool connected = _client.isConnected;
      loggingService.log('OPC-UA client connection status: ${connected ? 'Connected' : 'Disconnected'}');
      return connected;
    } catch (e) {
      loggingService.error('Error checking OPC-UA connection status: $e');
      return false;
    }
  }
}
