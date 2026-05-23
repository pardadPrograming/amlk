import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../admin_controller.dart';

class AdminPage extends StatefulWidget {
  const AdminPage({super.key});

  @override
  State<AdminPage> createState() => _AdminPageState();
}

class _AdminPageState extends State<AdminPage> {
  final controller = Get.find<AdminController>();

  @override
  void initState() {
    super.initState();
    controller.load();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('مدیریت کل'),
      body: ResponsivePage(
        child: Obx(() {
          if (controller.loading.value && controller.account.value == null) {
            return const Center(child: CircularProgressIndicator());
          }
          if (!controller.isAdmin) {
            return const AppCard(
              child: Text('این بخش فقط برای مدیران سطح کل پروژه فعال است.'),
            );
          }
          return ListView(
            children: [
              const GradientHeader(
                title: 'مدیریت کل پروژه',
                subtitle:
                    'مدیریت نقش‌های پلتفرمی، شهرها و پیشنهادهای لوکیشن کاربران',
                icon: Icons.admin_panel_settings_outlined,
              ),
              const SizedBox(height: 16),
              Wrap(
                spacing: 12,
                runSpacing: 12,
                children: [
                  SizedBox(
                    width: 260,
                    child: StatCard(
                      title: 'کاربران',
                      value: '${controller.users.length}',
                      icon: Icons.people_alt_outlined,
                    ),
                  ),
                  SizedBox(
                    width: 260,
                    child: StatCard(
                      title: 'کسب‌وکارها',
                      value: '${controller.businesses.length}',
                      icon: Icons.business_center_outlined,
                    ),
                  ),
                  SizedBox(
                    width: 260,
                    child: StatCard(
                      title: 'پیشنهادهای لوکیشن',
                      value: '${controller.suggestions.length}',
                      icon: Icons.add_location_alt_outlined,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              _AdminIdentity(),
              const SizedBox(height: 16),
              _PlatformSettingsPanel(),
              const SizedBox(height: 16),
              _CitiesPanel(),
              const SizedBox(height: 16),
              _SuggestionsPanel(),
            ],
          );
        }),
      ),
    );
  }
}

class _AdminIdentity extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final controller = Get.find<AdminController>();
    return AppCard(
      child: Obx(() {
        final account = controller.account.value;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'سطح دسترسی فعلی',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(account?.roles.join(', ') ?? '-'),
          ],
        );
      }),
    );
  }
}

class _PlatformSettingsPanel extends StatefulWidget {
  @override
  State<_PlatformSettingsPanel> createState() => _PlatformSettingsPanelState();
}

class _PlatformSettingsPanelState extends State<_PlatformSettingsPanel> {
  final otpApiKey = TextEditingController();
  final serviceSmsApiKey = TextEditingController();
  var showOtp = false;
  var showSms = false;

  @override
  void dispose() {
    otpApiKey.dispose();
    serviceSmsApiKey.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<AdminController>();
    return AppCard(
      child: Obx(() {
        final settings = controller.settings.value;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'تنظیمات مجموعه',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 6),
            Text(
              'API Keys',
              style: Theme.of(
                context,
              ).textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w800),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: otpApiKey,
              obscureText: !showOtp,
              decoration: InputDecoration(
                labelText: 'API Key OTP',
                helperText: settings.otpApiKeyMasked.isEmpty
                    ? 'هنوز ثبت نشده'
                    : 'کلید فعلی: ${settings.otpApiKeyMasked}',
                suffixIcon: IconButton(
                  tooltip: showOtp ? 'مخفی کردن' : 'نمایش',
                  onPressed: () => setState(() => showOtp = !showOtp),
                  icon: Icon(
                    showOtp
                        ? Icons.visibility_off_outlined
                        : Icons.visibility_outlined,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: serviceSmsApiKey,
              obscureText: !showSms,
              decoration: InputDecoration(
                labelText: 'API Key SMS خدماتی',
                helperText: settings.serviceSmsApiKeyMasked.isEmpty
                    ? 'هنوز ثبت نشده'
                    : 'کلید فعلی: ${settings.serviceSmsApiKeyMasked}',
                suffixIcon: IconButton(
                  tooltip: showSms ? 'مخفی کردن' : 'نمایش',
                  onPressed: () => setState(() => showSms = !showSms),
                  icon: Icon(
                    showSms
                        ? Icons.visibility_off_outlined
                        : Icons.visibility_outlined,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 14),
            Align(
              alignment: Alignment.centerLeft,
              child: FilledButton.icon(
                onPressed: controller.loading.value
                    ? null
                    : () {
                        controller.saveSettings(
                          otpApiKey: otpApiKey.text,
                          serviceSmsApiKey: serviceSmsApiKey.text,
                        );
                        otpApiKey.clear();
                        serviceSmsApiKey.clear();
                      },
                icon: const Icon(Icons.save_outlined),
                label: const Text('ذخیره API Keys'),
              ),
            ),
          ],
        );
      }),
    );
  }
}

class _CitiesPanel extends StatefulWidget {
  @override
  State<_CitiesPanel> createState() => _CitiesPanelState();
}

class _CitiesPanelState extends State<_CitiesPanel> {
  final name = TextEditingController();

  @override
  void dispose() {
    name.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<AdminController>();
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('شهرهای سیستمی', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: TextField(
                  controller: name,
                  decoration: const InputDecoration(labelText: 'نام شهر'),
                ),
              ),
              const SizedBox(width: 8),
              FilledButton.icon(
                onPressed: () {
                  controller.createCity(name.text);
                  name.clear();
                },
                icon: const Icon(Icons.add_rounded),
                label: const Text('افزودن'),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Obx(
            () => Wrap(
              spacing: 8,
              runSpacing: 8,
              children: controller.cities
                  .map((city) => Chip(label: Text(city.name)))
                  .toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _SuggestionsPanel extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final controller = Get.find<AdminController>();
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'پیشنهادهای در انتظار تایید',
            style: Theme.of(context).textTheme.titleMedium,
          ),
          const SizedBox(height: 12),
          Obx(() {
            if (controller.suggestions.isEmpty) {
              return const Text('پیشنهاد در انتظار بررسی وجود ندارد.');
            }
            return Column(
              children: controller.suggestions
                  .map(
                    (item) => ListTile(
                      leading: const Icon(Icons.location_on_outlined),
                      title: Text(item.name),
                      subtitle: Text(
                        item.manualParentName.isEmpty
                            ? '${item.type} / ${item.status}'
                            : '${item.type} / ${item.status} / ${item.manualParentName}',
                      ),
                      trailing: Wrap(
                        spacing: 8,
                        children: [
                          IconButton.filledTonal(
                            tooltip: 'تایید',
                            onPressed: () =>
                                controller.approveSuggestion(item.id),
                            icon: const Icon(Icons.check_rounded),
                          ),
                          IconButton.outlined(
                            tooltip: 'رد',
                            onPressed: () =>
                                controller.rejectSuggestion(item.id, ''),
                            icon: const Icon(Icons.close_rounded),
                          ),
                        ],
                      ),
                    ),
                  )
                  .toList(),
            );
          }),
        ],
      ),
    );
  }
}
