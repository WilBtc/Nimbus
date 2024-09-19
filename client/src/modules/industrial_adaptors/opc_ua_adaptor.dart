// src/modules/industrial_adaptors/opc_ua_adaptor.dart
import 'package:opcua/opcua.dart'; // Example OPC-UA client library, adjust based on actual library used.

class OPCUAAdaptor {
  final String serverUrl;
  final String securityPolicy;
  final String securityMode;
  final String username;
  final String password;

  OPCUAAdaptor({
    required this.serverUrl,
    this.securityPolicy = 'Basic256Sha256', // Default security policy
    this.securityMode = 'SignAndEncrypt',  // Default security mode
    this.username = '',
    this.password = '',
  });

  late OpcUaClient _client;

  // Initialize and connect the OPC-UA client with security settings
  Future<void> initialize() async {
    try {
      _client = OpcUaClient(
        endpointUrl: serverUrl,
        securityPolicy: securityPolicy,
        securityMode: securityMode,
      );

      // Configure user authentication if credentials are provided
      if (username.isNotEmpty && password.isNotEmpty) {
        _client.setUserToken(username, password);
      }

      await _client.connect();
      print('OPC-UA client connected to $serverUrl with security: $securityPolicy/$securityMode');
    } catch (e) {
      print('Error connecting OPC-UA client: $e');
    }
  }

  // Read a specific node's value
  Future<dynamic> readNode(String nodeId) async {
    try {
      var result = await _client.readNode(nodeId);
      print('Read OPC-UA node $nodeId: $result');
      return result;
    } catch (e) {
      print('Error reading OPC-UA node $nodeId: $e');
      return null;
    }
  }

  // Write a value to a specific node
  Future<void> writeNode(String nodeId, dynamic value) async {
    try {
      await _client.writeNode(nodeId, value);
      print('Wrote $value to OPC-UA node $nodeId');
    } catch (e) {
      print('Error writing to OPC-UA node $nodeId: $e');
    }
  }

  // Read multiple nodes at once
  Future<Map<String, dynamic>> readMultipleNodes(List<String> nodeIds) async {
    Map<String, dynamic> results = {};
    try {
      for (String nodeId in nodeIds) {
        var result = await readNode(nodeId);
        if (result != null) {
          results[nodeId] = result;
        }
      }
      print('Read multiple OPC-UA nodes: $results');
      return results;
    } catch (e) {
      print('Error reading multiple OPC-UA nodes: $e');
      return {};
    }
  }

  // Subscribe to changes in a specific node
  Future<void> subscribeToNode(String nodeId, Function(dynamic) onDataChange) async {
    try {
      await _client.subscribeToNode(nodeId, (value) {
        print('Data change detected on node $nodeId: $value');
        onDataChange(value);
      });
    } catch (e) {
      print('Error subscribing to OPC-UA node $nodeId: $e');
    }
  }

  // Unsubscribe from a specific node
  Future<void> unsubscribeFromNode(String nodeId) async {
    try {
      await _client.unsubscribeFromNode(nodeId);
      print('Unsubscribed from OPC-UA node $nodeId');
    } catch (e) {
      print('Error unsubscribing from OPC-UA node $nodeId: $e');
    }
  }

  // Disconnect the OPC-UA client safely
  void disconnect() {
    try {
      if (_client.isConnected) {
        _client.disconnect();
        print('OPC-UA client disconnected.');
      } else {
        print('OPC-UA client was already disconnected.');
      }
    } catch (e) {
      print('Error during OPC-UA client disconnection: $e');
    }
  }

  // Check if the client is currently connected
  bool isConnected() {
    bool connected = _client.isConnected;
    print('OPC-UA client connection status: ${connected ? 'Connected' : 'Disconnected'}');
    return connected;
  }
}
