// src/modules/decision_engine.dart

class DecisionEngine {
  final double criticalThreshold;
  final double warningThreshold;
  final double temperatureThreshold;
  final double pressureThreshold;

  DecisionEngine({
    this.criticalThreshold = 100.0,
    this.warningThreshold = 75.0,
    this.temperatureThreshold = 80.0,
    this.pressureThreshold = 50.0,
  });

  // Decides based on critical and warning thresholds, useful for automated responses.
  String makeDecision(double data) {
    if (data > criticalThreshold) {
      print('Decision: Critical condition detected, action required.');
      return 'ACTION_REQUIRED';
    } else if (data > warningThreshold) {
      print('Decision: Warning condition detected, monitor closely.');
      return 'WARNING';
    } else {
      print('Decision: Normal operation.');
      return 'NORMAL';
    }
  }

  // Example of a complex decision based on multiple data inputs.
  String multiFactorDecision(double temperature, double pressure) {
    if (temperature > temperatureThreshold && pressure > pressureThreshold) {
      print('Decision: Overheat and overpressure detected, emergency stop.');
      return 'EMERGENCY_STOP';
    } else if (temperature > temperatureThreshold) {
      print('Decision: Overheat detected, reduce load.');
      return 'REDUCE_LOAD';
    } else if (pressure > pressureThreshold) {
      print('Decision: Overpressure detected, vent system.');
      return 'VENT_SYSTEM';
    } else {
      print('Decision: System stable.');
      return 'STABLE';
    }
  }

  // Decision logic considering historical trends, predicting potential failures.
  String predictiveDecision(List<double> temperatureHistory, List<double> pressureHistory) {
    if (temperatureHistory.isEmpty || pressureHistory.isEmpty) {
      print('Decision: Insufficient data for predictive analysis.');
      return 'INSUFFICIENT_DATA';
    }

    double avgTemperature = temperatureHistory.reduce((a, b) => a + b) / temperatureHistory.length;
    double avgPressure = pressureHistory.reduce((a, b) => a + b) / pressureHistory.length;

    if (avgTemperature > temperatureThreshold * 0.9 && avgPressure > pressureThreshold * 0.9) {
      print('Decision: Conditions approaching critical, prepare for emergency procedures.');
      return 'PREPARE_FOR_EMERGENCY';
    } else if (avgTemperature > temperatureThreshold || avgPressure > pressureThreshold) {
      print('Decision: System showing signs of strain, recommended maintenance.');
      return 'RECOMMENDED_MAINTENANCE';
    } else {
      print('Decision: System operating within normal parameters.');
      return 'NORMAL_OPERATION';
    }
  }

  // Adaptive decision making based on external factors, useful for dynamic environments.
  String adaptiveDecision(double load, double vibration, Map<String, dynamic> externalFactors) {
    double loadThreshold = externalFactors['loadThreshold'] ?? 90.0;
    double vibrationThreshold = externalFactors['vibrationThreshold'] ?? 10.0;

    if (load > loadThreshold && vibration > vibrationThreshold) {
      print('Decision: High load and vibration detected, system adjustment required.');
      return 'ADJUST_SYSTEM';
    } else if (load > loadThreshold) {
      print('Decision: High load detected, optimize performance.');
      return 'OPTIMIZE_PERFORMANCE';
    } else if (vibration > vibrationThreshold) {
      print('Decision: High vibration detected, check system stability.');
      return 'CHECK_STABILITY';
    } else {
      print('Decision: System stable under current external conditions.');
      return 'STABLE_CONDITIONS';
    }
  }
}
