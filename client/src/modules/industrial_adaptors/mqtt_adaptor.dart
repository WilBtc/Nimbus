import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';
import 'logging_service.dart'; // Logging service for centralized log handling.

class MQTTAdaptor {
  final String broker;
  final String clientId;
  final int port;
  final int keepAlivePeriod;
  final bool cleanSession;

  late MqttServerClient _client;
  final LoggingService loggingService = LoggingService(); // For logging MQTT events.

  MQTTAdaptor({
    required this.broker,
    required this.clientId,
    this.port = 1883,
    this.keepAlivePeriod = 20,
    this.cleanSession = true,
  });

  /// Initialize the MQTT client and connect to the broker.
  Future<void> initialize() async {
    _client = MqttServerClient(broker, clientId)
      ..port = port
      ..logging(on: true)
      ..keepAlivePeriod = keepAlivePeriod
      ..onConnected = onConnected
      ..onDisconnected = onDisconnected
      ..onSubscribed = onSubscribed
      ..onSubscribeFail = onSubscribeFail
      ..onUnsubscribed = onUnsubscribed
      ..onAutoReconnect = onAutoReconnect
      ..onAutoReconnectAttempt = onAutoReconnectAttempt
      ..connectionMessage = MqttConnectMessage()
          .withClientIdentifier(clientId)
          .withWillTopic('nimbus/will')
          .withWillMessage('Nimbus client disconnected unexpectedly.')
          .startClean() // Set clean session based on configuration.
          .withWillQos(MqttQos.atLeastOnce);

    try {
      await _client.connect();
    } on NoConnectionException catch (e) {
      loggingService.error('MQTT client failed to connect: $e');
    } on SocketException catch (e) {
      loggingService.error('Socket exception during MQTT connection: $e');
    } catch (e) {
      loggingService.error('Error connecting MQTT client: $e');
    }
  }

  /// Callback for successful connection.
  void onConnected() {
    loggingService.log('MQTT client connected successfully to $broker:$port');
  }

  /// Callback for disconnection.
  void onDisconnected() {
    loggingService.log('MQTT client disconnected.');
  }

  /// Callback when a subscription is confirmed.
  void onSubscribed(String topic) {
    loggingService.log('Subscribed to topic: $topic');
  }

  /// Callback when subscription fails.
  void onSubscribeFail(String topic) {
    loggingService.error('Failed to subscribe to topic: $topic');
  }

  /// Callback for unsubscribed topics.
  void onUnsubscribed(String? topic) {
    loggingService.log('Unsubscribed from topic: $topic');
  }

  /// Callback for automatic reconnection.
  void onAutoReconnect() {
    loggingService.log('MQTT client attempting to reconnect...');
  }

  /// Callback for reconnection attempts.
  void onAutoReconnectAttempt() {
    loggingService.log('Attempting to reconnect MQTT client...');
  }

  /// Subscribe to a topic with specified QoS level.
  void subscribe(String topic, {MqttQos qos = MqttQos.atMostOnce}) {
    try {
      _client.subscribe(topic, qos);
      loggingService.log('Subscribed to topic: $topic with QoS: $qos');
      _client.updates?.listen((List<MqttReceivedMessage<MqttMessage>> messages) {
        final recMess = messages[0].payload as MqttPublishMessage;
        final payload = MqttPublishPayload.bytesToStringAsString(recMess.payload.message);
        loggingService.log('MQTT message received on topic ${messages[0].topic}: $payload');
      });
    } catch (e) {
      loggingService.error('Error subscribing to topic $topic: $e');
    }
  }

  /// Publish a message to a topic with specified QoS level.
  void publish(String topic, String message, {MqttQos qos = MqttQos.atMostOnce}) {
    try {
      final builder = MqttClientPayloadBuilder();
      builder.addString(message);
      _client.publishMessage(topic, qos, builder.payload!);
      loggingService.log('MQTT message published to $topic: $message');
    } catch (e) {
      loggingService.error('Error publishing message to topic $topic: $e');
    }
  }

  /// Disconnect the MQTT client gracefully.
  void disconnect() {
    try {
      _client.disconnect();
      loggingService.log('MQTT client disconnected.');
    } catch (e) {
      loggingService.error('Error during MQTT client disconnection: $e');
    }
  }

  /// Check if the client is currently connected.
  bool isConnected() {
    try {
      bool connected = _client.connectionStatus?.state == MqttConnectionState.connected;
      loggingService.log('MQTT client connection status: ${connected ? 'Connected' : 'Disconnected'}');
      return connected;
    } catch (e) {
      loggingService.error('Error checking MQTT connection status: $e');
      return false;
    }
  }
}
