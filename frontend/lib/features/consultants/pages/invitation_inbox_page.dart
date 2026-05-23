import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../consultants_controller.dart';

class InvitationInboxPage extends StatefulWidget {
  const InvitationInboxPage({super.key});

  @override
  State<InvitationInboxPage> createState() => _InvitationInboxPageState();
}

class _InvitationInboxPageState extends State<InvitationInboxPage> {
  final controller = Get.find<ConsultantsController>();

  @override
  void initState() {
    super.initState();
    controller.loadInbox();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('دعوت‌نامه‌ها'),
      body: ResponsivePage(
        child: Obx(() {
          final pending = controller.pendingInbox;
          if (pending.isEmpty) {
            return const SizedBox.shrink();
          }
          return ListView.separated(
            itemCount: pending.length,
            separatorBuilder: (_, _) => const SizedBox(height: 10),
            itemBuilder: (context, index) {
              final invite = pending[index];
              return Card(
                child: ListTile(
                  title: Text(invite.businessName),
                  subtitle: Text(
                    'نقش پیشنهادی: ${invite.role} | درصد سود: ${invite.commissionPercent}%',
                  ),
                  trailing: Wrap(
                    spacing: 8,
                    children: [
                      TextButton(
                        onPressed: () => controller.reject(invite.id),
                        child: const Text('رد'),
                      ),
                      FilledButton(
                        onPressed: () => controller.accept(invite.id),
                        child: const Text('قبول'),
                      ),
                    ],
                  ),
                ),
              );
            },
          );
        }),
      ),
    );
  }
}
