import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../app/app_theme.dart';
import '../../../shared/responsive.dart';
import '../auth_controller.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final phone = TextEditingController();
  final controller = Get.find<AuthController>();

  void _submit() {
    if (!controller.loading.value) {
      controller.requestOtp(phone.text);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ResponsivePage(
        maxWidth: 960,
        child: Center(
          child: LayoutBuilder(
            builder: (context, constraints) {
              final compact = constraints.maxWidth < 760;
              final form = AppCard(
                padding: 28,
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Row(
                      children: [
                        Container(
                          width: 48,
                          height: 48,
                          decoration: BoxDecoration(
                            gradient: const LinearGradient(
                              colors: [AppColors.primary, AppColors.secondary],
                            ),
                            borderRadius: BorderRadius.circular(14),
                          ),
                          child: const Icon(
                            Icons.apartment_rounded,
                            color: AppColors.accentGold,
                          ),
                        ),
                        const SizedBox(width: 14),
                        Expanded(
                          child: Text(
                            'ورود به مدیریت املاک',
                            style: Theme.of(context).textTheme.headlineSmall
                                ?.copyWith(fontWeight: FontWeight.w900),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 10),
                    Text(
                      'شماره موبایل را وارد کنید تا کد تایید برای شما صادر شود.',
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: AppColors.textSecondary,
                      ),
                    ),
                    const SizedBox(height: 24),
                    TextField(
                      controller: phone,
                      keyboardType: TextInputType.phone,
                      textInputAction: TextInputAction.done,
                      onSubmitted: (_) => _submit(),
                      decoration: const InputDecoration(
                        labelText: 'شماره موبایل',
                        prefixIcon: Icon(Icons.phone_iphone_rounded),
                      ),
                    ),
                    const SizedBox(height: 18),
                    Obx(
                      () => PremiumActionButton(
                        icon: Icons.sms_outlined,
                        onPressed: controller.loading.value ? null : _submit,
                        child: controller.loading.value
                            ? const SizedBox(
                                width: 22,
                                height: 22,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2.4,
                                  color: Colors.white,
                                ),
                              )
                            : const Text('دریافت کد تایید'),
                      ),
                    ),
                  ],
                ),
              );

              final hero = const GradientHeader(
                title: 'Amlak CRM',
                subtitle:
                    'پایه یک CRM حرفه‌ای برای آژانس‌های املاک، با احراز هویت OTP و داشبورد مدیریتی.',
                icon: Icons.real_estate_agent_outlined,
              );

              if (compact) {
                return Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [hero, const SizedBox(height: 18), form],
                );
              }
              return Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(width: 430, height: 260, child: hero),
                  const SizedBox(width: 22),
                  SizedBox(width: 440, child: form),
                ],
              );
            },
          ),
        ),
      ),
    );
  }
}
