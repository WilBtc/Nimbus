import 'package:modbus/modbus.dart'; // Example library for Modbus communication, adjust as needed.
import 'logging_service.dart'; // Assuming a logging service for logging Modbus operations.

class ModbusAdaptor {
  final String deviceIp;
  final int port;
  final int timeout; // Timeout in seconds for Modbus operations

  ModbusAdaptor({required this.deviceIp, this.port = 502, this.timeout = 5});

  late ModbusClient _client;
  final LoggingService loggingService = LoggingService(); // Logging Service for centralized logging

  /// Initialize and connect the Modbus client.
  Future<void> initialize() async {
    try {
      _client = ModbusClient(deviceIp, port: port)
        ..timeout = Duration(seconds: timeout);
      await _client.connect();
      loggingService.log('Modbus client connected to $deviceIp:$port');
    } catch (e) {
      loggingService.error('Error connecting Modbus client: $e');
      rethrow;
    }
  }

  /// Read a single register from the Modbus device.
  Future<int?> readRegister(int registerAddress) async {
    try {
      var result = await _client.readHoldingRegisters(registerAddress, 1);
      loggingService.log('Read Modbus register $registerAddress: ${result[0]}');
      return result[0];
    } catch (e) {
      loggingService.error('Error reading Modbus register $registerAddress: $e');
      return null;
    }
  }

  /// Write a single value to a specific register.
  Future<void> writeRegister(int registerAddress, int value) async {
    try {
      await _client.writeSingleRegister(registerAddress, value);
      loggingService.log('Wrote $value to Modbus register $registerAddress');
    } catch (e) {
      loggingService.error('Error writing to Modbus register $registerAddress: $e');
    }
  }

  /// Read multiple registers for bulk data acquisition.
  Future<List<int>?> readMultipleRegisters(int startAddress, int count) async {
    try {
      var result = await _client.readHoldingRegisters(startAddress, count);
      loggingService.log('Read Modbus registers from $startAddress to ${startAddress + count - 1}: $result');
      return result;
    } catch (e) {
      loggingService.error('Error reading multiple Modbus registers starting at $startAddress: $e');
      return null;
    }
  }

  /// Write multiple registers in a single operation.
  Future<void> writeMultipleRegisters(int startAddress, List<int> values) async {
    try {
      await _client.writeMultipleRegisters(startAddress, values);
      loggingService.log('Wrote values $values to Modbus registers starting at $startAddress');
    } catch (e) {
      loggingService.error('Error writing multiple Modbus registers starting at $startAddress: $e');
    }
  }

  /// Check the connection status of the Modbus client.
  bool isConnected() {
    try {
      bool connected = _client.isConnected;
      loggingService.log('Modbus client connection status: ${connected ? 'Connected' : 'Disconnected'}');
      return connected;
    } catch (e) {
      loggingService.error('Error checking Modbus connection status: $e');
      return false;
    }
  }

  /// Disconnect the Modbus client safely.
  void disconnect() {
    try {
      if (_client.isConnected) {
        _client.disconnect();
        loggingService.log('Modbus client disconnected.');
      } else {
        loggingService.log('Modbus client was already disconnected.');
      }
    } catch (e) {
      loggingService.error('Error during Modbus client disconnection: $e');
    }
  }
}
