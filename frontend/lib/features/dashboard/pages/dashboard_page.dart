import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../app/app.dart';
import '../../../app/app_theme.dart';
import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../../business/business_controller.dart';
import '../dashboard_controller.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  final controller = Get.find<DashboardController>();

  @override
  void initState() {
    super.initState();
    if (Get.arguments is Map<String, dynamic>) {
      Get.find<BusinessController>().selected.value = Business.fromJson(
        Map<String, dynamic>.from(Get.arguments),
      );
    }
    controller.load();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('داشبورد'),
      body: ResponsivePage(
        child: Obx(() {
          final stats = controller.stats;
          final business = Get.find<BusinessController>().selected.value;
          return Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              GradientHeader(
                title: business?.name ?? 'داشبورد املاک',
                subtitle: 'نمای سریع وضعیت دفتر، مشاورین و مسیرهای عملیاتی.',
                icon: Icons.dashboard_customize_outlined,
              ),
              const SizedBox(height: 20),
              LayoutBuilder(
                builder: (context, constraints) {
                  final width = constraints.maxWidth < 720
                      ? constraints.maxWidth
                      : (constraints.maxWidth - 24) / 3;
                  return Wrap(
                    spacing: 12,
                    runSpacing: 12,
                    children: [
                      SizedBox(
                        width: width,
                        child: StatCard(
                          title: 'مشاورین',
                          value: '${stats['consultantsCount'] ?? 0}',
                          icon: Icons.groups_2_outlined,
                        ),
                      ),
                      SizedBox(
                        width: width,
                        child: StatCard(
                          title: 'فایل‌های ملکی',
                          value: '${stats['propertiesCount'] ?? 0}',
                          icon: Icons.home_work_outlined,
                        ),
                      ),
                      SizedBox(
                        width: width,
                        child: StatCard(
                          title: 'قراردادها',
                          value: '${stats['contractsCount'] ?? 0}',
                          icon: Icons.description_outlined,
                        ),
                      ),
                    ],
                  );
                },
              ),
              const SizedBox(height: 24),
              AppCard(
                padding: 18,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'دسترسی سریع',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w800,
                        color: AppColors.primary,
                      ),
                    ),
                    const SizedBox(height: 14),
                    Wrap(
                      spacing: 12,
                      runSpacing: 12,
                      children: [
                        FilledButton.icon(
                          onPressed: () => Get.offNamed(AppRoutes.consultants),
                          icon: const Icon(Icons.manage_accounts_outlined),
                          label: const Text('مدیریت مشاورین'),
                        ),
                        OutlinedButton.icon(
                          onPressed: () =>
                              Get.offNamed(AppRoutes.businessSettings),
                          icon: const Icon(Icons.settings_outlined),
                          label: const Text('تنظیمات کسب‌وکار'),
                        ),
                        OutlinedButton.icon(
                          onPressed: () => Get.offNamed(AppRoutes.properties),
                          icon: const Icon(Icons.apartment_outlined),
                          label: const Text('فایل‌ها'),
                        ),
                        OutlinedButton.icon(
                          onPressed: () => Get.offNamed(AppRoutes.customers),
                          icon: const Icon(Icons.people_outline),
                          label: const Text('مشتریان'),
                        ),
                        OutlinedButton.icon(
                          onPressed: () => Get.offNamed(AppRoutes.contracts),
                          icon: const Icon(Icons.assignment_outlined),
                          label: const Text('قراردادها'),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          );
        }),
      ),
    );
  }
}
