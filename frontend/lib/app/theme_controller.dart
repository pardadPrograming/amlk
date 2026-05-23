import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

class ThemeController extends GetxController {
  ThemeController(this._box);

  static const _key = 'theme_mode_v2';
  final GetStorage _box;
  final mode = ThemeMode.system.obs;

  ThemeMode get themeMode => mode.value;

  bool get isDark {
    if (mode.value == ThemeMode.system) {
      return WidgetsBinding.instance.platformDispatcher.platformBrightness ==
          Brightness.dark;
    }
    return mode.value == ThemeMode.dark;
  }

  String get label {
    return switch (mode.value) {
      ThemeMode.system => 'سیستم',
      ThemeMode.dark => 'تاریک',
      ThemeMode.light => 'روشن',
    };
  }

  IconData get icon {
    return switch (mode.value) {
      ThemeMode.system => Icons.brightness_auto_rounded,
      ThemeMode.dark => Icons.dark_mode_rounded,
      ThemeMode.light => Icons.light_mode_rounded,
    };
  }

  @override
  void onInit() {
    super.onInit();
    mode.value = _parse(_box.read(_key));
  }

  void cycle() {
    mode.value = switch (mode.value) {
      ThemeMode.system => ThemeMode.dark,
      ThemeMode.dark => ThemeMode.light,
      ThemeMode.light => ThemeMode.system,
    };
    _save();
  }

  void setMode(ThemeMode value) {
    mode.value = value;
    _save();
  }

  ThemeMode _parse(dynamic value) {
    return switch (value) {
      'dark' => ThemeMode.dark,
      'light' => ThemeMode.light,
      _ => ThemeMode.system,
    };
  }

  void _save() {
    if (mode.value == ThemeMode.system) {
      _box.remove(_key);
    } else {
      _box.write(_key, mode.value == ThemeMode.dark ? 'dark' : 'light');
    }
  }
}
