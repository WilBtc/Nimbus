// src/modules/industrial_adaptors/mqtt_adaptor.dart
import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';

class MQTTAdaptor {
  final String broker;
  final String clientId;
  final int port;
  final int keepAlivePeriod;
  final bool cleanSession;

  MQTTAdaptor({
    required this.broker,
    required this.clientId,
    this.port = 1883,
    this.keepAlivePeriod = 20,
    this.cleanSession = true,
  });

  late MqttServerClient _client;

  // Initialize the MQTT client and connect to the broker
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
          .startClean() // Set clean session based on configuration
          .withWillQos(MqttQos.atLeastOnce);

    try {
      await _client.connect();
    } on NoConnectionException catch (e) {
      print('MQTT client failed to connect: $e');
    } on SocketException catch (e) {
      print('Socket exception during MQTT connection: $e');
    } catch (e) {
      print('Error connecting MQTT client: $e');
    }
  }

  // Callback for successful connection
  void onConnected() {
    print('MQTT client connected successfully to $broker:$port');
  }

  // Callback for disconnection
  void onDisconnected() {
    print('MQTT client disconnected.');
  }

  // Callback when a subscription is confirmed
  void onSubscribed(String topic) {
    print('Subscribed to topic: $topic');
  }

  // Callback when subscription fails
  void onSubscribeFail(String topic) {
    print('Failed to subscribe to topic: $topic');
  }

  // Callback for unsubscribed topics
  void onUnsubscribed(String? topic) {
    print('Unsubscribed from topic: $topic');
  }

  // Callback for automatic reconnection
  void onAutoReconnect() {
    print('MQTT client attempting to reconnect...');
  }

  // Callback for reconnection attempts
  void onAutoReconnectAttempt() {
    print('Attempting to reconnect MQTT client...');
  }

  // Subscribe to a topic with specified QoS level
  void subscribe(String topic, {MqttQos qos = MqttQos.atMostOnce}) {
    try {
      _client.subscribe(topic, qos);
      print('Subscribed to topic: $topic with QoS: $qos');
      _client.updates?.listen((List<MqttReceivedMessage<MqttMessage>> messages) {
        final recMess = messages[0].payload as MqttPublishMessage;
        final payload = MqttPublishPayload.bytesToStringAsString(recMess.payload.message);
        print('MQTT message received on topic ${messages[0].topic}: $payload');
      });
    } catch (e) {
      print('Error subscribing to topic $topic: $e');
    }
  }

  // Publish a message to a topic with specified QoS level
  void publish(String topic, String message, {MqttQos qos = MqttQos.atMostOnce}) {
    try {
      final builder = MqttClientPayloadBuilder();
      builder.addString(message);
      _client.publishMessage(topic, qos, builder.payload!);
      print('MQTT message published to $topic: $message');
    } catch (e) {
      print('Error publishing message to topic $topic: $e');
    }
  }

  // Disconnect the MQTT client gracefully
  void disconnect() {
    try {
      _client.disconnect();
      print('MQTT client disconnected.');
    } catch (e) {
      print('Error during MQTT client disconnection: $e');
    }
  }

  // Check if the client is currently connected
  bool isConnected() {
    bool connected = _client.connectionStatus?.state == MqttConnectionState.connected;
    print('MQTT client connection status: ${connected ? 'Connected' : 'Disconnected'}');
    return connected;
  }
}
