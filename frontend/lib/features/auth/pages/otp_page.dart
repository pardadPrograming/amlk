import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../app/app_theme.dart';
import '../../../shared/responsive.dart';
import '../auth_controller.dart';

class OtpPage extends StatefulWidget {
  const OtpPage({super.key});

  @override
  State<OtpPage> createState() => _OtpPageState();
}

class _OtpPageState extends State<OtpPage> {
  final code = TextEditingController();
  final controller = Get.find<AuthController>();

  void _submit() {
    if (!controller.loading.value) {
      controller.verifyOtp(code.text);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: GlassAppBar(title: const Text('تایید شماره')),
      body: ResponsivePage(
        maxWidth: 520,
        child: Center(
          child: AppCard(
            padding: 28,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const GradientHeader(
                  title: 'کد تایید',
                  subtitle: 'کد شش رقمی ارسال‌شده را وارد کنید.',
                  icon: Icons.verified_user_outlined,
                ),
                const SizedBox(height: 22),
                Obx(
                  () => Text(
                    controller.phone.value.isEmpty
                        ? 'برای تست می‌توانید آخرین OTP را از سرور بگیرید.'
                        : 'ارسال شده به ${controller.phone.value}',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: AppColors.textSecondary,
                    ),
                  ),
                ),
                const SizedBox(height: 10),
                Obx(
                  () => AnimatedSwitcher(
                    duration: const Duration(milliseconds: 250),
                    child: controller.devCode.value.isEmpty
                        ? OutlinedButton.icon(
                            onPressed: controller.loading.value
                                ? null
                                : controller.loadLatestTestOtp,
                            icon: const Icon(Icons.bolt_outlined),
                            label: const Text('دریافت آخرین OTP تست'),
                          )
                        : Container(
                            key: ValueKey(controller.devCode.value),
                            padding: const EdgeInsets.symmetric(
                              horizontal: 14,
                              vertical: 12,
                            ),
                            decoration: BoxDecoration(
                              color: AppColors.accentGold.withValues(
                                alpha: 0.18,
                              ),
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: AppColors.accentGold.withValues(
                                  alpha: 0.48,
                                ),
                              ),
                            ),
                            child: Row(
                              children: [
                                Expanded(
                                  child: Text(
                                    'کد تست: ${controller.devCode.value}',
                                    style: const TextStyle(
                                      fontWeight: FontWeight.w800,
                                      color: AppColors.primary,
                                    ),
                                  ),
                                ),
                                TextButton(
                                  onPressed: controller.loading.value
                                      ? null
                                      : controller.loadLatestTestOtp,
                                  child: const Text('بروزرسانی'),
                                ),
                              ],
                            ),
                          ),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: code,
                  keyboardType: TextInputType.number,
                  textInputAction: TextInputAction.done,
                  onSubmitted: (_) => _submit(),
                  decoration: const InputDecoration(
                    labelText: 'کد تایید',
                    prefixIcon: Icon(Icons.password_rounded),
                  ),
                ),
                const SizedBox(height: 18),
                Obx(
                  () => PremiumActionButton(
                    icon: Icons.login_rounded,
                    onPressed: controller.loading.value ? null : _submit,
                    child: const Text('ورود'),
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
