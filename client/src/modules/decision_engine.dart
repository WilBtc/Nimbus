// src/modules/decision_engine.dart

import 'logging_service.dart';  // Assuming you have a logging service for better log management

class DecisionEngine {
  final double criticalThreshold;
  final double warningThreshold;
  final double temperatureThreshold;
  final double pressureThreshold;

  final LoggingService loggingService = LoggingService();  // For logging decision events

  DecisionEngine({
    this.criticalThreshold = 100.0,
    this.warningThreshold = 75.0,
    this.temperatureThreshold = 80.0,
    this.pressureThreshold = 50.0,
  });

  /// Decides based on critical and warning thresholds, useful for automated responses.
  String makeDecision(double data) {
    try {
      if (data > criticalThreshold) {
        loggingService.log('Decision: Critical condition detected, action required.');
        return 'ACTION_REQUIRED';
      } else if (data > warningThreshold) {
        loggingService.log('Decision: Warning condition detected, monitor closely.');
        return 'WARNING';
      } else {
        loggingService.log('Decision: Normal operation.');
        return 'NORMAL';
      }
    } catch (e) {
      loggingService.error('Error in makeDecision: $e');
      return 'ERROR';
    }
  }

  /// Example of a complex decision based on multiple data inputs.
  String multiFactorDecision(double temperature, double pressure) {
    try {
      if (temperature > temperatureThreshold && pressure > pressureThreshold) {
        loggingService.log('Decision: Overheat and overpressure detected, emergency stop.');
        return 'EMERGENCY_STOP';
      } else if (temperature > temperatureThreshold) {
        loggingService.log('Decision: Overheat detected, reduce load.');
        return 'REDUCE_LOAD';
      } else if (pressure > pressureThreshold) {
        loggingService.log('Decision: Overpressure detected, vent system.');
        return 'VENT_SYSTEM';
      } else {
        loggingService.log('Decision: System stable.');
        return 'STABLE';
      }
    } catch (e) {
      loggingService.error('Error in multiFactorDecision: $e');
      return 'ERROR';
    }
  }

  /// Decision logic considering historical trends, predicting potential failures.
  String predictiveDecision(List<double> temperatureHistory, List<double> pressureHistory) {
    try {
      if (temperatureHistory.isEmpty || pressureHistory.isEmpty) {
        loggingService.warn('Decision: Insufficient data for predictive analysis.');
        return 'INSUFFICIENT_DATA';
      }

      double avgTemperature = _calculateAverage(temperatureHistory);
      double avgPressure = _calculateAverage(pressureHistory);

      if (avgTemperature > temperatureThreshold * 0.9 && avgPressure > pressureThreshold * 0.9) {
        loggingService.log('Decision: Conditions approaching critical, prepare for emergency procedures.');
        return 'PREPARE_FOR_EMERGENCY';
      } else if (avgTemperature > temperatureThreshold || avgPressure > pressureThreshold) {
        loggingService.log('Decision: System showing signs of strain, recommended maintenance.');
        return 'RECOMMENDED_MAINTENANCE';
      } else {
        loggingService.log('Decision: System operating within normal parameters.');
        return 'NORMAL_OPERATION';
      }
    } catch (e) {
      loggingService.error('Error in predictiveDecision: $e');
      return 'ERROR';
    }
  }

  /// Adaptive decision making based on external factors, useful for dynamic environments.
  String adaptiveDecision(double load, double vibration, Map<String, dynamic> externalFactors) {
    try {
      double loadThreshold = externalFactors['loadThreshold'] ?? 90.0;
      double vibrationThreshold = externalFactors['vibrationThreshold'] ?? 10.0;

      if (load > loadThreshold && vibration > vibrationThreshold) {
        loggingService.log('Decision: High load and vibration detected, system adjustment required.');
        return 'ADJUST_SYSTEM';
      } else if (load > loadThreshold) {
        loggingService.log('Decision: High load detected, optimize performance.');
        return 'OPTIMIZE_PERFORMANCE';
      } else if (vibration > vibrationThreshold) {
        loggingService.log('Decision: High vibration detected, check system stability.');
        return 'CHECK_STABILITY';
      } else {
        loggingService.log('Decision: System stable under current external conditions.');
        return 'STABLE_CONDITIONS';
      }
    } catch (e) {
      loggingService.error('Error in adaptiveDecision: $e');
      return 'ERROR';
    }
  }

  /// Helper method to calculate the average of a list of doubles.
  double _calculateAverage(List<double> data) {
    return data.reduce((a, b) => a + b) / data.length;
  }
}
