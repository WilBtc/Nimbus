import 'core/at_protocol.dart';
import 'core/security_manager.dart';
import 'core/config.dart';
import 'modules/industrial_adaptors/modbus_adaptor.dart';
import 'modules/industrial_adaptors/opc_ua_adaptor.dart';
import 'modules/industrial_adaptors/mqtt_adaptor.dart';
import 'modules/data_processor.dart';
import 'modules/decision_engine.dart';
import 'logging_service.dart'; // Added centralized logging service

void main() async {
  // Initialize logging service
  final loggingService = LoggingService();

  loggingService.log('Nimbus Client initializing...');
  
  try {
    // Load configuration settings
    final config = Config();
    await config.load();
    loggingService.log('Configuration loaded successfully.');

    // Initialize security manager
    final securityManager = SecurityManager(config.atSign);
    await securityManager.initialize();
    loggingService.log('Security Manager initialized.');

    // Initialize atProtocol client
    final atProtocol = AtProtocol(securityManager);
    await atProtocol.initialize();
    loggingService.log('atProtocol client initialized.');

    // Initialize industrial protocol adaptors
    final modbusAdaptor = ModbusAdaptor(deviceIp: '192.168.1.10', port: 502);
    await modbusAdaptor.initialize();
    loggingService.log('Modbus Adaptor connected.');

    final opcUaAdaptor = OPCUAAdaptor(serverUrl: 'opc.tcp://192.168.1.20:4840');
    await opcUaAdaptor.initialize();
    loggingService.log('OPC-UA Adaptor connected.');

    final mqttAdaptor = MQTTAdaptor(broker: 'broker.hivemq.com', clientId: 'nimbus_client', port: 1883);
    await mqttAdaptor.initialize();
    loggingService.log('MQTT Adaptor connected.');

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
        loggingService.log('Modbus Data: $modbusData');

        // Process data
        List<double> filteredData = dataProcessor.filterNoise([modbusData.toDouble()]);
        double aggregatedData = dataProcessor.aggregateData(filteredData);

        // Decision-making based on processed data
        String decision = decisionEngine.makeDecision(aggregatedData);
        loggingService.log('Decision: $decision');

        // Sending data to atProtocol server
        await atProtocol.sendMessage('aggregated_data', aggregatedData.toString());
        loggingService.log('Data sent to atProtocol: $aggregatedData');

        // Publish decision to MQTT
        mqttAdaptor.publish('nimbus/decisions', decision);
        loggingService.log('Decision published to MQTT: $decision');
      }

      // Reading data from OPC-UA server
      dynamic opcData = await opcUaAdaptor.readNode('ns=2;s=TemperatureSensor1');
      if (opcData != null) {
        loggingService.log('OPC-UA Data: $opcData');

        // Convert OPC-UA data to a double if necessary for processing
        doubl
