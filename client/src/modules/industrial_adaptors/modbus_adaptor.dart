// src/modules/industrial_adaptors/modbus_adaptor.dart
import 'package:modbus/modbus.dart'; // Example library for Modbus communication, adjust as needed.

class ModbusAdaptor {
  final String deviceIp;
  final int port;
  final int timeout; // Timeout in seconds for Modbus operations

  ModbusAdaptor({required this.deviceIp, this.port = 502, this.timeout = 5});

  late ModbusClient _client;

  // Initialize and connect the Modbus client
  Future<void> initialize() async {
    try {
      _client = ModbusClient(deviceIp, port: port)
        ..timeout = Duration(seconds: timeout);
      await _client.connect();
      print('Modbus client connected to $deviceIp:$port');
    } catch (e) {
      print('Error connecting Modbus client: $e');
    }
  }

  // Read a single register from the Modbus device
  Future<int?> readRegister(int registerAddress) async {
    try {
      var result = await _client.readHoldingRegisters(registerAddress, 1);
      print('Read Modbus register $registerAddress: ${result[0]}');
      return result[0];
    } catch (e) {
      print('Error reading Modbus register: $e');
      return null;
    }
  }

  // Write a single value to a specific register
  Future<void> writeRegister(int registerAddress, int value) async {
    try {
      await _client.writeSingleRegister(registerAddress, value);
      print('Wrote $value to Modbus register $registerAddress');
    } catch (e) {
      print('Error writing to Modbus register: $e');
    }
  }

  // Read multiple registers, useful for bulk data acquisition
  Future<List<int>?> readMultipleRegisters(int startAddress, int count) async {
    try {
      var result = await _client.readHoldingRegisters(startAddress, count);
      print('Read Modbus registers from $startAddress to ${startAddress + count - 1}: $result');
      return result;
    } catch (e) {
      print('Error reading multiple Modbus registers: $e');
      return null;
    }
  }

  // Write multiple registers in a single operation
  Future<void> writeMultipleRegisters(int startAddress, List<int> values) async {
    try {
      await _client.writeMultipleRegisters(startAddress, values);
      print('Wrote values $values to Modbus registers starting at $startAddress');
    } catch (e) {
      print('Error writing multiple Modbus registers: $e');
    }
  }

  // Check the connection status of the Modbus client
  bool isConnected() {
    bool connected = _client.isConnected;
    print('Modbus client connection status: ${connected ? 'Connected' : 'Disconnected'}');
    return connected;
  }

  // Disconnect the Modbus client safely
  void disconnect() {
    try {
      if (_client.isConnected) {
        _client.disconnect();
        print('Modbus client disconnected.');
      } else {
        print('Modbus client was already disconnected.');
      }
    } catch (e) {
      print('Error during Modbus client disconnection: $e');
    }
  }
}
