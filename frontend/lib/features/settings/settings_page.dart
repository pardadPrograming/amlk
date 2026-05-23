import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../app/app.dart';
import '../../app/theme_controller.dart';
import '../../shared/responsive.dart';
import '../business/business_controller.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final themeController = Get.find<ThemeController>();
    final businessController = Get.find<BusinessController>();
    return PanelScaffold(
      title: const Text('تنظیمات'),
      body: ResponsivePage(
        maxWidth: 820,
        child: ListView(
          children: [
            const GradientHeader(
              title: 'تنظیمات',
              subtitle: 'مدیریت حساب، امنیت، کسب‌وکار و ظاهر برنامه.',
              icon: Icons.tune_rounded,
            ),
            const SizedBox(height: 18),
            AppCard(
              padding: 8,
              child: Column(
                children: [
                  _SettingsTile(
                    icon: Icons.account_circle_outlined,
                    title: 'پروفایل کاربری',
                    subtitle:
                        'نام، نام خانوادگی و شهر فعال برای موقعیت های مکانی',
                    onTap: () => Get.toNamed(AppRoutes.profile),
                  ),
                  const Divider(height: 1),
                  _SettingsTile(
                    icon: Icons.verified_user_outlined,
                    title: 'امنیت حساب',
                    subtitle: 'نشست‌های فعال، دستگاه‌ها و آخرین فعالیت',
                    onTap: () => Get.toNamed(AppRoutes.security),
                  ),
                  const Divider(height: 1),
                  _SettingsTile(
                    icon: Icons.business_center_outlined,
                    title: 'تنظیمات کسب‌وکار',
                    subtitle: 'اطلاعات دفتر، لایسنس و لوگو',
                    onTap: () => Get.toNamed(AppRoutes.businessSettings),
                  ),
                  const Divider(height: 1),
                  _SettingsTile(
                    icon: Icons.admin_panel_settings_outlined,
                    title: 'مدیریت کل پروژه',
                    subtitle: 'نقش‌های پلتفرمی، شهرها و تایید لوکیشن‌ها',
                    onTap: () => Get.toNamed(AppRoutes.admin),
                  ),
                  const Divider(height: 1),
                  Obx(
                    () => _SettingsTile(
                      icon: themeController.icon,
                      title: 'ظاهر برنامه',
                      subtitle: 'تم فعلی: ${themeController.label}',
                      onTap: themeController.cycle,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 18),
            Obx(() {
              final business = businessController.selected.value;
              if (business == null) return const SizedBox.shrink();
              return AppCard(
                padding: 8,
                child: _SettingsTile(
                  icon: Icons.logout_rounded,
                  title: 'خروج از ${business.name}',
                  subtitle:
                      'تمام دسترسی‌ها و اتصال شما به این کسب‌وکار قطع می‌شود.',
                  danger: true,
                  onTap: () => _confirmLeaveBusiness(context),
                ),
              );
            }),
          ],
        ),
      ),
    );
  }

  void _confirmLeaveBusiness(BuildContext context) {
    final controller = Get.find<BusinessController>();
    final business = controller.selected.value;
    if (business == null) return;
    Get.dialog(
      AlertDialog(
        title: const Text('خروج از کسب‌وکار'),
        content: Text(
          'با خروج از ${business.name} تمام دسترسی‌ها، نقش و اتصال شما به این کسب‌وکار قطع می‌شود. ادامه می‌دهید؟',
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          Obx(
            () => FilledButton.icon(
              onPressed: controller.loading.value
                  ? null
                  : controller.leaveSelectedBusiness,
              icon: const Icon(Icons.logout_rounded),
              label: const Text('خروج'),
              style: FilledButton.styleFrom(
                backgroundColor: Theme.of(context).colorScheme.error,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.onTap,
    this.danger = false,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final VoidCallback onTap;
  final bool danger;

  @override
  Widget build(BuildContext context) {
    final color = danger
        ? Theme.of(context).colorScheme.error
        : Theme.of(context).colorScheme.primary;
    return ListTile(
      leading: Container(
        width: 44,
        height: 44,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: danger
                ? [color, color.withValues(alpha: 0.72)]
                : const [Color(0xFF2F80ED), Color(0xFF4DE1FF)],
            begin: Alignment.topRight,
            end: Alignment.bottomLeft,
          ),
          borderRadius: BorderRadius.circular(14),
        ),
        child: Icon(icon, color: Colors.white),
      ),
      title: Text(
        title,
        style: TextStyle(
          fontWeight: FontWeight.w800,
          color: danger ? color : null,
        ),
      ),
      subtitle: Text(subtitle),
      trailing: Text(
        '›',
        textDirection: TextDirection.ltr,
        style: Theme.of(
          context,
        ).textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.w800),
      ),
      onTap: onTap,
    );
  }
}
