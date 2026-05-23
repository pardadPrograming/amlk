import 'dart:async';
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

import 'app/app.dart';
import 'app/theme_controller.dart';
import 'core/api/api_client.dart';
import 'core/logging/app_logger.dart';
import 'core/storage/token_storage.dart';
import 'features/auth/auth_controller.dart';
import 'features/admin/admin_controller.dart';
import 'features/business/business_controller.dart';
import 'features/chats/chats_controller.dart';
import 'features/consultants/consultants_controller.dart';
import 'features/customers/customers_controller.dart';
import 'features/dashboard/dashboard_controller.dart';
import 'features/locations/locations_controller.dart';
import 'features/properties/properties_controller.dart';

Future<void> main() async {
  runZonedGuarded(
    () async {
      WidgetsFlutterBinding.ensureInitialized();

      FlutterError.onError = (details) {
        FlutterError.presentError(details);
        AppLogger.error(
          'flutter',
          details.exception,
          stackTrace: details.stack,
          details: details.context,
        );
      };

      PlatformDispatcher.instance.onError = (error, stack) {
        AppLogger.error('platform', error, stackTrace: stack);
        return true;
      };

      await GetStorage.init();

      Get.put(ThemeController(GetStorage()), permanent: true);
      Get.put(TokenStorage(), permanent: true);
      Get.put(ApiClient(Get.find<TokenStorage>()), permanent: true);
      Get.put(
        AuthController(Get.find<ApiClient>(), Get.find<TokenStorage>()),
        permanent: true,
      );
      Get.put(AdminController(Get.find<ApiClient>()), permanent: true);
      Get.put(BusinessController(Get.find<ApiClient>()), permanent: true);
      Get.put(DashboardController(Get.find<ApiClient>()), permanent: true);
      Get.put(ChatsController(Get.find<ApiClient>()), permanent: true);
      Get.put(ConsultantsController(Get.find<ApiClient>()), permanent: true);
      Get.put(CustomersController(Get.find<ApiClient>()), permanent: true);
      Get.put(LocationsController(Get.find<ApiClient>()), permanent: true);
      Get.put(PropertiesController(Get.find<ApiClient>()), permanent: true);

      runApp(const AmlakApp());
    },
    (error, stackTrace) {
      AppLogger.error('zone', error, stackTrace: stackTrace);
    },
  );
}
