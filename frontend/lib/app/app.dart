import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../features/auth/pages/login_page.dart';
import '../features/auth/pages/otp_page.dart';
import '../features/auth/pages/profile_page.dart';
import '../features/auth/pages/security_page.dart';
import '../features/admin/pages/admin_page.dart';
import '../features/business/pages/business_settings_page.dart';
import '../features/chats/chats_controller.dart';
import '../features/chats/pages/chats_page.dart';
import '../features/business/pages/create_business_page.dart';
import '../features/consultants/pages/consultants_page.dart';
import '../features/consultants/pages/invitation_inbox_page.dart';
import '../features/customers/pages/customers_page.dart';
import '../features/dashboard/pages/dashboard_page.dart';
import '../features/locations/pages/locations_page.dart';
import '../features/placeholders/placeholder_page.dart';
import '../features/properties/pages/latest_files_page.dart';
import '../features/properties/pages/properties_page.dart';
import '../features/settings/settings_page.dart';
import 'app_theme.dart';
import 'theme_controller.dart';
import 'translations.dart';

class AmlakApp extends StatelessWidget {
  const AmlakApp({super.key});

  @override
  Widget build(BuildContext context) {
    final themeController = Get.find<ThemeController>();
    return Obx(
      () => GetMaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'Amlak CRM',
        locale: const Locale('fa', 'IR'),
        fallbackLocale: const Locale('en', 'US'),
        translations: AppTranslations(),
        textDirection: TextDirection.rtl,
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        themeMode: themeController.themeMode,
        initialRoute: AppRoutes.login,
        getPages: [
          GetPage(name: AppRoutes.login, page: () => const LoginPage()),
          GetPage(name: AppRoutes.otp, page: () => const OtpPage()),
          GetPage(name: AppRoutes.profile, page: () => const ProfilePage()),
          GetPage(name: AppRoutes.settings, page: () => const SettingsPage()),
          GetPage(name: AppRoutes.security, page: () => const SecurityPage()),
          GetPage(name: AppRoutes.admin, page: () => const AdminPage()),
          GetPage(
            name: AppRoutes.createBusiness,
            page: () => const CreateBusinessPage(),
          ),
          GetPage(name: AppRoutes.dashboard, page: () => const DashboardPage()),
          GetPage(
            name: AppRoutes.consultants,
            page: () => const ConsultantsPage(),
          ),
          GetPage(
            name: AppRoutes.inbox,
            page: () => const InvitationInboxPage(),
          ),
          GetPage(name: AppRoutes.locations, page: () => const LocationsPage()),
          GetPage(
            name: AppRoutes.businessSettings,
            page: () => const BusinessSettingsPage(),
          ),
          GetPage(
            name: AppRoutes.properties,
            page: () => const PropertiesPage(),
          ),
          GetPage(
            name: AppRoutes.latestFiles,
            page: () => const LatestFilesPage(),
          ),
          GetPage(
            name: AppRoutes.chats,
            page: () =>
                const ChatsPage(initialSection: ChatSection.privateChats),
          ),
          GetPage(
            name: AppRoutes.vaults,
            page: () =>
                const ChatsPage(initialSection: ChatSection.personalVaults),
          ),
          GetPage(
            name: AppRoutes.propertyCreate,
            page: () => const PropertyCreatePage(),
          ),
          GetPage(name: AppRoutes.customers, page: () => const CustomersPage()),
          GetPage(
            name: AppRoutes.contracts,
            page: () => const PlaceholderPage(title: 'قراردادها'),
          ),
        ],
      ),
    );
  }
}

class AppRoutes {
  static const login = '/login';
  static const otp = '/otp';
  static const profile = '/profile';
  static const settings = '/settings';
  static const security = '/settings/security';
  static const admin = '/admin';
  static const createBusiness = '/business/create';
  static const dashboard = '/dashboard';
  static const consultants = '/consultants';
  static const inbox = '/invitations';
  static const locations = '/locations';
  static const businessSettings = '/business/settings';
  static const properties = '/properties';
  static const latestFiles = '/properties/latest';
  static const chats = '/chats';
  static const vaults = '/vaults';
  static const propertyCreate = '/properties/create';
  static const customers = '/customers';
  static const contracts = '/contracts';
}
