import 'package:test/test.dart';
import 'package:nimbus_client/modules/data_processor.dart';

void main() {
  group('DataProcessor', () {
    final dataProcessor = DataProcessor(noiseThreshold: 0.05);

    test('filterNoise removes values below threshold', () {
      final data = [0.01, 0.02, 0.06, 0.1];
      final result = dataProcessor.filterNoise(data);
      expect(result, [0.06, 0.1]);
    });

    test('filterNoise handles empty list', () {
      final data = [];
      final result = dataProcessor.filterNoise(data);
      expect(result, []);
    });

    test('aggregateData sums values correctly', () {
      final data = [1.0, 2.0, 3.0];
      final result = dataProcessor.aggregateData(data);
      expect(result, 6.0);
    });

    test('aggregateData handles empty list', () {
      final data = [];
      final result = dataProcessor.aggregateData(data);
      expect(result, 0.0);
    });

    test('averageData calculates mean correctly', () {
      final data = [2.0, 4.0, 6.0];
      final result = dataProcessor.averageData(data);
      expect(result, 4.0);
    });

    test('averageData handles single value', () {
      final data = [5.0];
      final result = dataProcessor.averageData(data);
      expect(result, 5.0);
    });

    test('medianData calculates median correctly for odd length', () {
      final data = [3.0, 1.0, 2.0];
      final result = dataProcessor.medianData(data);
      expect(result, 2.0);
    });

    test('medianData calculates median correctly for even length', () {
      final data = [1.0, 2.0, 3.0, 4.0];
      final result = dataProcessor.medianData(data);
      expect(result, 2.5);
    });

    test('standardDeviation calculates correctly', () {
      final data = [2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0];
      final stdDev = dataProcessor.standardDeviation(data);
      expect(stdDev, closeTo(2.0, 0.001));
    });

    test('normalizeData normalizes correctly', () {
      final data = [2.0, 4.0, 6.0];
      final result = dataProcessor.normalizeData(data);
      expect(result, [0.0, 0.5, 1.0]);
    });

    test('normalizeData handles equal values', () {
      final data = [5.0, 5.0, 5.0];
      final result = dataProcessor.normalizeData(data);
      expect(result, [5.0, 5.0, 5.0]);
    });

    test('normalizeData handles single value', () {
      final data = [7.0];
      final result = dataProcessor.normalizeData(data);
      expect(result, [7.0]);
    });
  });
}
