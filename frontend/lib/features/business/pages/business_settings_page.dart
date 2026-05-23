import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../business_controller.dart';

class BusinessSettingsPage extends StatelessWidget {
  const BusinessSettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<BusinessController>();
    return PanelScaffold(
      title: const Text('تنظیمات کسب‌وکار'),
      body: ResponsivePage(
        child: Obx(() {
          final business = controller.selected.value;
          return Card(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    business?.name ?? 'کسب‌وکار انتخاب نشده',
                    style: Theme.of(context).textTheme.headlineSmall,
                  ),
                  const SizedBox(height: 12),
                  Text('وضعیت لایسنس: ${business?.licenseStatus ?? '-'}'),
                  const SizedBox(height: 16),
                  OutlinedButton.icon(
                    onPressed: () {},
                    icon: const Icon(Icons.image_outlined),
                    label: const Text('بارگذاری لوگو'),
                  ),
                ],
              ),
            ),
          );
        }),
      ),
    );
  }
}
