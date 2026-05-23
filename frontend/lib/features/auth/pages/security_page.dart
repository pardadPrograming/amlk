import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../auth_controller.dart';

class SecurityPage extends StatefulWidget {
  const SecurityPage({super.key});

  @override
  State<SecurityPage> createState() => _SecurityPageState();
}

class _SecurityPageState extends State<SecurityPage> {
  final controller = Get.find<AuthController>();

  @override
  void initState() {
    super.initState();
    controller.loadSecurityProfile();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('امنیت حساب'),
      body: ResponsivePage(
        child: Obx(
          () => RefreshIndicator(
            onRefresh: controller.loadSecurityProfile,
            child: ListView(
              children: [
                const GradientHeader(
                  title: 'امنیت حساب',
                  subtitle:
                      'نشست‌های فعال، دستگاه‌ها و آخرین فعالیت حساب کاربری خود را مدیریت کنید.',
                  icon: Icons.verified_user_outlined,
                ),
                const SizedBox(height: 18),
                _ActivityCard(lastActivityAt: controller.lastActivityAt.value),
                const SizedBox(height: 18),
                Text(
                  'نشست‌های فعال',
                  style: Theme.of(
                    context,
                  ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.w900),
                ),
                const SizedBox(height: 12),
                if (controller.sessions.isEmpty)
                  const AppCard(
                    child: Text('نشست فعالی برای نمایش وجود ندارد.'),
                  )
                else
                  ...controller.sessions.map(
                    (session) => Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: _SessionTile(
                        session: session,
                        onRevoke: session.current
                            ? null
                            : () => controller.revokeSession(session.id),
                      ),
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

class _ActivityCard extends StatelessWidget {
  const _ActivityCard({required this.lastActivityAt});

  final DateTime? lastActivityAt;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Row(
        children: [
          const _IconBadge(icon: Icons.timeline_outlined),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'آخرین فعالیت',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
                ),
                const SizedBox(height: 4),
                Text(_formatDate(lastActivityAt)),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _SessionTile extends StatelessWidget {
  const _SessionTile({required this.session, required this.onRevoke});

  final UserSession session;
  final VoidCallback? onRevoke;

  @override
  Widget build(BuildContext context) {
    final icon = session.deviceType == 'mobile'
        ? Icons.smartphone_outlined
        : Icons.desktop_windows_outlined;
    return AppCard(
      padding: 18,
      child: LayoutBuilder(
        builder: (context, constraints) {
          final compact = constraints.maxWidth < 620;
          final content = [
            _IconBadge(icon: icon),
            const SizedBox(width: 14, height: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Flexible(
                        child: Text(
                          session.deviceName,
                          style: Theme.of(context).textTheme.titleMedium
                              ?.copyWith(fontWeight: FontWeight.w800),
                        ),
                      ),
                      if (session.current) ...[
                        const SizedBox(width: 8),
                        const Chip(label: Text('نشست فعلی')),
                      ],
                    ],
                  ),
                  const SizedBox(height: 6),
                  Text('${session.browser} - ${session.os} - ${session.ip}'),
                  const SizedBox(height: 4),
                  Text('آخرین فعالیت: ${_formatDate(session.lastSeenAt)}'),
                ],
              ),
            ),
            const SizedBox(width: 12, height: 12),
            OutlinedButton.icon(
              onPressed: onRevoke,
              icon: const Icon(Icons.power_settings_new_rounded),
              label: const Text('غیرفعال‌سازی'),
            ),
          ];
          if (compact) {
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Row(children: content.take(3).toList()),
                const SizedBox(height: 12),
                content.last,
              ],
            );
          }
          return Row(children: content);
        },
      ),
    );
  }
}

class _IconBadge extends StatelessWidget {
  const _IconBadge({required this.icon});

  final IconData icon;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 48,
      height: 48,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF2F80ED), Color(0xFF4DE1FF)],
          begin: Alignment.topRight,
          end: Alignment.bottomLeft,
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Icon(icon, color: Colors.white),
    );
  }
}

String _formatDate(DateTime? value) {
  if (value == null || value.year <= 1) {
    return 'ثبت نشده';
  }
  final local = value.toLocal();
  String two(int input) => input.toString().padLeft(2, '0');
  return '${local.year}/${two(local.month)}/${two(local.day)} - ${two(local.hour)}:${two(local.minute)}';
}
