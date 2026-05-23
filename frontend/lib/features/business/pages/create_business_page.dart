import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../business_controller.dart';

class CreateBusinessPage extends StatefulWidget {
  const CreateBusinessPage({super.key});

  @override
  State<CreateBusinessPage> createState() => _CreateBusinessPageState();
}

class _CreateBusinessPageState extends State<CreateBusinessPage> {
  final name = TextEditingController();
  final phone = TextEditingController();
  final address = TextEditingController();
  final hours = TextEditingController(text: 'شنبه تا پنجشنبه ۹ تا ۱۸');
  final license = TextEditingController();
  final controller = Get.find<BusinessController>();

  void _submit() {
    if (!controller.loading.value) {
      controller.create(
        name: name.text,
        phone: phone.text,
        address: address.text,
        workingHours: hours.text,
        licenseNumber: license.text,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: GlassAppBar(title: const Text('تعریف کسب‌وکار املاک')),
      body: ResponsivePage(
        maxWidth: 880,
        child: AppCard(
          padding: 28,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const GradientHeader(
                title: 'ساخت دفتر املاک',
                subtitle:
                    'اطلاعات پایه کسب‌وکار، لایسنس آزمایشی و نقش مالک ساخته می‌شود.',
                icon: Icons.storefront_outlined,
              ),
              const SizedBox(height: 24),
              Wrap(
                runSpacing: 16,
                spacing: 16,
                children: [
                  _field(name, 'نام املاک', Icons.apartment_outlined),
                  _field(
                    phone,
                    'شماره تماس اصلی',
                    Icons.phone_in_talk_outlined,
                    keyboard: TextInputType.phone,
                  ),
                  _field(address, 'آدرس', Icons.location_on_outlined),
                  _field(hours, 'ساعت فعالیت', Icons.schedule_outlined),
                  _field(
                    license,
                    'شماره جواز',
                    Icons.workspace_premium_outlined,
                    submit: true,
                  ),
                ],
              ),
              const SizedBox(height: 22),
              Obx(
                () => PremiumActionButton(
                  icon: Icons.dashboard_customize_outlined,
                  onPressed: controller.loading.value ? null : _submit,
                  child: const Text('ساخت کسب‌وکار و ورود به داشبورد'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _field(
    TextEditingController controller,
    String label,
    IconData icon, {
    TextInputType? keyboard,
    bool submit = false,
  }) {
    return SizedBox(
      width: 390,
      child: TextField(
        controller: controller,
        keyboardType: keyboard,
        textInputAction: submit ? TextInputAction.done : TextInputAction.next,
        onSubmitted: (_) =>
            submit ? _submit() : FocusScope.of(context).nextFocus(),
        decoration: InputDecoration(labelText: label, prefixIcon: Icon(icon)),
      ),
    );
  }
}
