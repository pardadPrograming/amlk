import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../auth_controller.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final firstName = TextEditingController();
  final lastName = TextEditingController();
  final selectedCityId = ''.obs;
  final controller = Get.find<AuthController>();

  void _submit() {
    if (!controller.loading.value) {
      controller.completeProfile(
        firstName.text,
        lastName.text,
        selectedCityId.value,
      );
    }
  }

  @override
  void initState() {
    super.initState();
    final currentUser = controller.user.value;
    firstName.text = currentUser?.firstName ?? '';
    lastName.text = currentUser?.lastName ?? '';
    selectedCityId.value = currentUser?.cityId ?? '';
    controller.loadCities();
  }

  @override
  void dispose() {
    firstName.dispose();
    lastName.dispose();
    selectedCityId.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: GlassAppBar(
        title: const Text('تکمیل پروفایل'),
        actions: const [ThemeModeAction()],
      ),
      body: ResponsivePage(
        maxWidth: 640,
        child: Center(
          child: AppCard(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const GradientHeader(
                  title: 'پروفایل کاربری',
                  subtitle:
                      'نام و نام خانوادگی را جدا وارد کنید تا در تیم و دعوت‌نامه‌ها دقیق نمایش داده شود.',
                  icon: Icons.badge_outlined,
                ),
                const SizedBox(height: 22),
                LayoutBuilder(
                  builder: (context, constraints) {
                    final compact = constraints.maxWidth < 520;
                    final fields = [
                      Expanded(
                        child: TextField(
                          controller: firstName,
                          textInputAction: TextInputAction.next,
                          onSubmitted: (_) =>
                              FocusScope.of(context).nextFocus(),
                          decoration: const InputDecoration(
                            labelText: 'نام',
                            prefixIcon: Icon(Icons.person_outline_rounded),
                          ),
                        ),
                      ),
                      SizedBox(
                        width: compact ? 0 : 12,
                        height: compact ? 12 : 0,
                      ),
                      Expanded(
                        child: TextField(
                          controller: lastName,
                          textInputAction: TextInputAction.done,
                          onSubmitted: (_) => _submit(),
                          decoration: const InputDecoration(
                            labelText: 'نام خانوادگی',
                            prefixIcon: Icon(Icons.badge_outlined),
                          ),
                        ),
                      ),
                    ];
                    if (compact) {
                      return Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: fields,
                      );
                    }
                    return Row(children: fields);
                  },
                ),
                const SizedBox(height: 18),
                Obx(
                  () => DropdownButtonFormField<String>(
                    initialValue:
                        selectedCityId.value.isEmpty ||
                            !controller.cities.any(
                              (city) => city.id == selectedCityId.value,
                            )
                        ? null
                        : selectedCityId.value,
                    items: controller.cities
                        .map(
                          (city) => DropdownMenuItem<String>(
                            value: city.id,
                            child: Text(city.name),
                          ),
                        )
                        .toList(),
                    onChanged: controller.loading.value
                        ? null
                        : (value) => selectedCityId.value = value ?? '',
                    decoration: const InputDecoration(
                      labelText: 'شهر',
                      prefixIcon: Icon(Icons.location_city_outlined),
                    ),
                  ),
                ),
                const SizedBox(height: 18),
                Obx(
                  () => PremiumActionButton(
                    icon: Icons.arrow_back_rounded,
                    matchTextDirection: true,
                    onPressed:
                        controller.loading.value || selectedCityId.value.isEmpty
                        ? null
                        : _submit,
                    child: const Text('ادامه'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
