import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../shared/responsive.dart';
import '../consultants_controller.dart';

class ConsultantsPage extends StatefulWidget {
  const ConsultantsPage({super.key});

  @override
  State<ConsultantsPage> createState() => _ConsultantsPageState();
}

class _ConsultantsPageState extends State<ConsultantsPage> {
  final controller = Get.find<ConsultantsController>();

  @override
  void initState() {
    super.initState();
    controller.load();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('مدیریت مشاورین'),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Get.dialog(const InviteConsultantDialog()),
        icon: const Icon(Icons.person_add_alt_1),
        label: const Text('دعوت مشاور'),
      ),
      body: ResponsivePage(
        child: Obx(
          () => Column(
            children: [
              LayoutBuilder(
                builder: (context, constraints) {
                  final compact = constraints.maxWidth < 760;
                  final search = TextField(
                    onChanged: (v) => controller.query.value = v,
                    decoration: const InputDecoration(
                      prefixIcon: Icon(Icons.search),
                      labelText: 'جستجو بر اساس نام یا شماره',
                    ),
                  );
                  final filter = SegmentedButton<String>(
                    segments: const [
                      ButtonSegment(value: 'all', label: Text('همه')),
                      ButtonSegment(value: 'owner', label: Text('مالک')),
                      ButtonSegment(value: 'manager', label: Text('مدیر')),
                      ButtonSegment(value: 'consultant', label: Text('مشاور')),
                    ],
                    selected: {controller.roleFilter.value},
                    onSelectionChanged: (v) =>
                        controller.roleFilter.value = v.first,
                  );
                  if (compact) {
                    return Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [search, const SizedBox(height: 12), filter],
                    );
                  }
                  return Row(
                    children: [
                      Expanded(child: search),
                      const SizedBox(width: 12),
                      filter,
                    ],
                  );
                },
              ),
              const SizedBox(height: 16),
              Expanded(
                child: Card(
                  child: ListView.separated(
                    itemCount: controller.filteredMembers.length,
                    separatorBuilder: (_, _) => const Divider(height: 1),
                    itemBuilder: (context, index) {
                      final member = controller.filteredMembers[index];
                      return ListTile(
                        leading: const CircleAvatar(
                          child: Icon(Icons.person_outline),
                        ),
                        title: Text(
                          member.userDisplayName.isEmpty
                              ? member.userPhone
                              : member.userDisplayName,
                        ),
                        subtitle: Text(
                          '${member.userPhone} | نقش: ${member.role} | درصد سود: ${member.commissionPercent}% | وضعیت: ${member.status}',
                        ),
                        trailing: PopupMenuButton<String>(
                          onSelected: (value) {
                            if (value == 'manager') {
                              controller.updateMember(member, role: 'manager');
                            }
                            if (value == 'consultant') {
                              controller.updateMember(
                                member,
                                role: 'consultant',
                              );
                            }
                            if (value == 'disable') {
                              controller.updateMember(
                                member,
                                status: 'disabled',
                              );
                            }
                          },
                          itemBuilder: (_) => const [
                            PopupMenuItem(
                              value: 'manager',
                              child: Text('ارتقا به مدیر'),
                            ),
                            PopupMenuItem(
                              value: 'consultant',
                              child: Text('تبدیل به مشاور'),
                            ),
                            PopupMenuItem(
                              value: 'disable',
                              child: Text('غیرفعال کردن'),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
                ),
              ),
              if (controller.invitations
                  .where((e) => e.status == 'pending')
                  .isNotEmpty) ...[
                const SizedBox(height: 12),
                Align(
                  alignment: Alignment.centerRight,
                  child: Text(
                    "دعوت‌های در انتظار: ${controller.invitations.where((e) => e.status == 'pending').length}",
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class InviteConsultantDialog extends StatefulWidget {
  const InviteConsultantDialog({super.key});

  @override
  State<InviteConsultantDialog> createState() => _InviteConsultantDialogState();
}

class _InviteConsultantDialogState extends State<InviteConsultantDialog> {
  final phone = TextEditingController();
  final commission = TextEditingController(text: '40');
  String role = 'consultant';

  void _submit() {
    Get.find<ConsultantsController>().invite(
      phone.text,
      double.tryParse(commission.text) ?? 0,
      role,
    );
  }

  @override
  void dispose() {
    phone.dispose();
    commission.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('دعوت مشاور'),
      content: SizedBox(
        width: 420,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: phone,
              keyboardType: TextInputType.phone,
              textInputAction: TextInputAction.next,
              onSubmitted: (_) => FocusScope.of(context).nextFocus(),
              decoration: const InputDecoration(labelText: 'شماره موبایل'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: commission,
              keyboardType: TextInputType.number,
              textInputAction: TextInputAction.done,
              onSubmitted: (_) => _submit(),
              decoration: const InputDecoration(labelText: 'درصد سود پیش‌فرض'),
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<String>(
              initialValue: role,
              items: const [
                DropdownMenuItem(value: 'consultant', child: Text('مشاور')),
                DropdownMenuItem(value: 'manager', child: Text('مدیر')),
              ],
              onChanged: (v) => setState(() => role = v ?? 'consultant'),
              decoration: const InputDecoration(labelText: 'نقش پیشنهادی'),
            ),
          ],
        ),
      ),
      actions: [
        TextButton(onPressed: Get.back, child: const Text('انصراف')),
        ElevatedButton(onPressed: _submit, child: const Text('ارسال دعوت')),
      ],
    );
  }
}
