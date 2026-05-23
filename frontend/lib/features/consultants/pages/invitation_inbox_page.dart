import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../../properties/properties_controller.dart';
import '../consultants_controller.dart';

class InvitationInboxPage extends StatefulWidget {
  const InvitationInboxPage({super.key});

  @override
  State<InvitationInboxPage> createState() => _InvitationInboxPageState();
}

class _InvitationInboxPageState extends State<InvitationInboxPage> {
  final controller = Get.find<ConsultantsController>();
  final properties = Get.find<PropertiesController>();

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    await Future.wait([
      controller.loadInbox(),
      properties.loadShareRequests(),
      properties.loadNotifications(),
    ]);
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('صندوقچه پیام‌ها'),
      body: ResponsivePage(
        maxWidth: 920,
        child: RefreshIndicator(
          onRefresh: _load,
          child: Obx(() {
            final invitations = controller.pendingInbox;
            final incomingShares = properties.ownerShareRequests
                .where((item) => item.status == 'pending')
                .toList();
            final outgoingShares = properties.requesterShareRequests
                .where(
                  (item) =>
                      item.status == 'approved' || item.status == 'pending',
                )
                .toList();
            final notifications = properties.notifications;
            final total =
                invitations.length +
                incomingShares.length +
                outgoingShares.length +
                notifications.length;
            if (total == 0) {
              return ListView(
                physics: const AlwaysScrollableScrollPhysics(),
                children: const [
                  SizedBox(height: 140),
                  Icon(Icons.inbox_outlined, size: 56),
                  SizedBox(height: 12),
                  Center(child: Text('پیامی برای نمایش وجود ندارد')),
                ],
              );
            }
            return ListView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.only(bottom: 24),
              children: [
                if (notifications.isNotEmpty) ...[
                  const _InboxSectionTitle(
                    title: 'پیام‌ها و فایل‌های پیشنهادی',
                  ),
                  ...notifications.map(
                    (item) => _NotificationInboxCard(
                      notification: item,
                      onRead: () => properties.markNotificationRead(item),
                    ),
                  ),
                ],
                if (invitations.isNotEmpty) ...[
                  const _InboxSectionTitle(title: 'دعوت‌ها'),
                  ...invitations.map(
                    (invite) => _InvitationInboxCard(
                      invite: invite,
                      onAccept: () => controller.accept(invite.id),
                      onReject: () => controller.reject(invite.id),
                    ),
                  ),
                ],
                if (incomingShares.isNotEmpty) ...[
                  const _InboxSectionTitle(title: 'درخواست‌های مشارکت دریافتی'),
                  ...incomingShares.map(
                    (request) => _ShareInboxCard(
                      request: request,
                      incoming: true,
                      onAccept: () => properties.decideShare(request, true),
                      onReject: () => properties.decideShare(request, false),
                    ),
                  ),
                ],
                if (outgoingShares.isNotEmpty) ...[
                  const _InboxSectionTitle(title: 'درخواست‌های مشارکت من'),
                  ...outgoingShares.map(
                    (request) => _ShareInboxCard(
                      request: request,
                      incoming: false,
                      onReceive: request.status == 'approved'
                          ? () => properties.receiveSharedFile(request)
                          : null,
                    ),
                  ),
                ],
              ],
            );
          }),
        ),
      ),
    );
  }
}

class _InboxSectionTitle extends StatelessWidget {
  const _InboxSectionTitle({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(4, 18, 4, 8),
      child: Text(title, style: Theme.of(context).textTheme.titleMedium),
    );
  }
}

class _NotificationInboxCard extends StatelessWidget {
  const _NotificationInboxCard({
    required this.notification,
    required this.onRead,
  });

  final NotificationModel notification;
  final VoidCallback onRead;

  @override
  Widget build(BuildContext context) {
    final unread = notification.readAt.isEmpty;
    return Card(
      child: ListTile(
        leading: Icon(
          unread ? Icons.markunread_outlined : Icons.drafts_outlined,
        ),
        title: Text(notification.title),
        subtitle: Text(notification.body),
        trailing: unread
            ? TextButton(onPressed: onRead, child: const Text('خواندم'))
            : null,
      ),
    );
  }
}

class _InvitationInboxCard extends StatelessWidget {
  const _InvitationInboxCard({
    required this.invite,
    required this.onAccept,
    required this.onReject,
  });

  final Invitation invite;
  final VoidCallback onAccept;
  final VoidCallback onReject;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: const Icon(Icons.person_add_alt_1_outlined),
        title: Text(invite.businessName),
        subtitle: Text(
          'نقش پیشنهادی: ${invite.role} | درصد سود: ${invite.commissionPercent}%',
        ),
        trailing: Wrap(
          spacing: 8,
          children: [
            TextButton(onPressed: onReject, child: const Text('رد')),
            FilledButton(onPressed: onAccept, child: const Text('قبول')),
          ],
        ),
      ),
    );
  }
}

class _ShareInboxCard extends StatelessWidget {
  const _ShareInboxCard({
    required this.request,
    required this.incoming,
    this.onAccept,
    this.onReject,
    this.onReceive,
  });

  final PropertyShareRequestModel request;
  final bool incoming;
  final VoidCallback? onAccept;
  final VoidCallback? onReject;
  final VoidCallback? onReceive;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: const Icon(Icons.handshake_outlined),
        title: Text(request.propertyTitle),
        subtitle: Text(
          incoming
              ? '${request.requesterName} درخواست مشارکت با سهم ${request.commissionPercent}% داده است'
              : 'وضعیت درخواست: ${request.status} | سهم ${request.commissionPercent}%',
        ),
        trailing: incoming
            ? Wrap(
                spacing: 8,
                children: [
                  TextButton(onPressed: onReject, child: const Text('رد')),
                  FilledButton(onPressed: onAccept, child: const Text('قبول')),
                ],
              )
            : onReceive == null
            ? null
            : FilledButton(
                onPressed: onReceive,
                child: const Text('دریافت فایل'),
              ),
      ),
    );
  }
}
