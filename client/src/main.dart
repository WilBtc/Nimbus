// src/main.dart

import 'core/at_protocol.dart';
import 'core/security_manager.dart';
import 'core/config.dart';
import 'modules/industrial_adaptors/modbus_adaptor.dart';
import 'modules/industrial_adaptors/opc_ua_adaptor.dart';
import 'modules/industrial_adaptors/mqtt_adaptor.dart';
import 'modules/data_processor.dart';
import 'modules/decision_engine.dart';

void main() async {
  // Load configuration settings
  final config = Config();
  await config.load();

  // Initialize security manager
  final securityManager = SecurityManager(config.atSign);
  await securityManager.initialize();

  // Initialize atProtocol client
  final atProtocol = AtProtocol(securityManager);
  await atProtocol.initialize();

  // Initialize industrial protocol adaptors
  final modbusAdaptor = ModbusAdaptor(deviceIp: '192.168.1.10', port: 502);
  await modbusAdaptor.initialize();

  final opcUaAdaptor = OPCUAAdaptor(serverUrl: 'opc.tcp://192.168.1.20:4840');
  await opcUaAdaptor.initialize();

  final mqttAdaptor = MQTTAdaptor(broker: 'broker.hivemq.com', clientId: 'nimbus_client', port: 1883);
  await mqttAdaptor.initialize();

  // Initialize data processor and decision engine
  final dataProcessor = DataProcessor(noiseThreshold: 0.05);
  final decisionEngine = DecisionEngine(
    criticalThreshold: 100.0,
    warningThreshold: 75.0,
    temperatureThreshold: 80.0,
    pressureThreshold: 50.0,
  );

  // Sample Workflow: Data Acquisition, Processing, and Decision-Making
  try {
    // Reading data from Modbus device
    int? modbusData = await modbusAdaptor.readRegister(1);
    if (modbusData != null) {
      print('Modbus Data: $modbusData');

      // Process data
      List<double> filteredData = dataProcessor.filterNoise([modbusData.toDouble()]);
      double aggregatedData = dataProcessor.aggregateData(filteredData);

      // Decision-making based on processed data
      String decision = decisionEngine.makeDecision(aggregatedData);
      print('Decision: $decision');

      // Sending data to atProtocol server
      await atProtocol.sendMessage('aggregated_data', aggregatedData.toString());

      // Publish decision to MQTT
      mqttAdaptor.publish('nimbus/decisions', decision);
    }

    // Reading data from OPC-UA server
    dynamic opcData = await opcUaAdaptor.readNode('ns=2;s=TemperatureSensor1');
    if (opcData != null) {
      print('OPC-UA Data: $opcData');

      // Convert OPC-UA data to a double if necessary for processing
      double opcValue = double.tryParse(opcData.toString()) ?? 0.0;

      // Further processing with DataProcessor and DecisionEngine
      List<double> opcFilteredData = dataProcessor.filterNoise([opcValue]);
      double opcAggregatedData = dataProcessor.aggregateData(opcFilteredData);

      // Use predictive decision-making based on historical trends
      String predictiveDecision = decisionEngine.predictiveDecision(
        [modbusData!.toDouble(), opcAggregatedData],
        [opcValue, 55.0], // Example pressure data for demonstration
      );
      print('Predictive Decision: $predictiveDecision');

      // Send predictive decision to the server and publish to MQTT
      await atProtocol.sendMessage('predictive_decision', predictiveDecision);
      mqttAdaptor.publish('nimbus/predictive_decisions', predictiveDecision);
    }
  } catch (e) {
    print('Error in processing workflow: $e');
  } finally {
    // Clean up resources
    modbusAdaptor.disconnect();
    opcUaAdaptor.disconnect();
    mqttAdaptor.disconnect();
    atProtocol.close();
    print('Nimbus Client operations completed.');
  }
}
