import 'package:flutter/foundation.dart';

class AppLogger {
  static void error(
    String scope,
    Object error, {
    StackTrace? stackTrace,
    Object? details,
  }) {
    final buffer = StringBuffer()
      ..writeln('[AmlakCRM][$scope] ERROR')
      ..writeln(error);
    if (details != null) {
      buffer
        ..writeln('details:')
        ..writeln(details);
    }
    if (stackTrace != null) {
      buffer
        ..writeln('stack:')
        ..writeln(stackTrace);
    }
    debugPrint(buffer.toString());
  }

  static void info(String scope, Object message) {
    debugPrint('[AmlakCRM][$scope] $message');
  }
}
