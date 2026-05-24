import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../app/app.dart';
import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../properties_controller.dart';

class PropertyOffersPage extends StatefulWidget {
  const PropertyOffersPage({super.key});

  @override
  State<PropertyOffersPage> createState() => _PropertyOffersPageState();
}

class _PropertyOffersPageState extends State<PropertyOffersPage>
    with SingleTickerProviderStateMixin {
  late final TabController _tabs;
  final controller = Get.find<PropertiesController>();

  @override
  void initState() {
    super.initState();
    _tabs = TabController(length: 2, vsync: this);
    controller.loadOffers();
  }

  @override
  void dispose() {
    _tabs.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('پیشنهادها'),
      body: ResponsivePage(
        maxWidth: 1100,
        child: Column(
          children: [
            Material(
              color: Colors.transparent,
              child: TabBar(
                controller: _tabs,
                tabs: const [
                  Tab(text: 'دریافتی'),
                  Tab(text: 'ارسالی'),
                ],
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: RefreshIndicator(
                onRefresh: controller.loadOffers,
                child: TabBarView(
                  controller: _tabs,
                  children: [
                    Obx(
                      () => _OfferList(
                        offers: controller.incomingOffers,
                        incoming: true,
                        controller: controller,
                      ),
                    ),
                    Obx(
                      () => _OfferList(
                        offers: controller.outgoingOffers,
                        incoming: false,
                        controller: controller,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _OfferList extends StatelessWidget {
  const _OfferList({
    required this.offers,
    required this.incoming,
    required this.controller,
  });

  final List<PropertyOfferModel> offers;
  final bool incoming;
  final PropertiesController controller;

  @override
  Widget build(BuildContext context) {
    if (offers.isEmpty) {
      return ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        children: const [
          SizedBox(height: 120),
          Icon(Icons.local_offer_outlined, size: 54),
          SizedBox(height: 10),
          Center(child: Text('پیشنهادی برای نمایش وجود ندارد')),
        ],
      );
    }
    return ListView.separated(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.only(bottom: 24),
      itemCount: offers.length,
      separatorBuilder: (_, _) => const SizedBox(height: 10),
      itemBuilder: (context, index) => _OfferCard(
        offer: offers[index],
        incoming: incoming,
        controller: controller,
      ),
    );
  }
}

class _OfferCard extends StatelessWidget {
  const _OfferCard({
    required this.offer,
    required this.incoming,
    required this.controller,
  });

  final PropertyOfferModel offer;
  final bool incoming;
  final PropertiesController controller;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(_statusIcon(offer.status)),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    offer.propertyTitle,
                    style: theme.textTheme.titleMedium,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                _StatusChip(status: offer.status),
              ],
            ),
            const SizedBox(height: 10),
            Text(
              incoming
                  ? 'درخواست: ${offer.requestTitle} | سهم پیشنهادی ${offer.commissionPercent.toStringAsFixed(0)}%'
                  : 'گیرنده: ${offer.contactName.isEmpty ? offer.requesterName : offer.contactName} | مچ ${offer.score}%',
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                if (!incoming &&
                    (offer.status == 'candidate' || offer.status == 'sent'))
                  FilledButton.icon(
                    onPressed: () => _send(context),
                    icon: const Icon(Icons.send_outlined),
                    label: Text(
                      offer.status == 'sent' ? 'ارسال دوباره' : 'ارسال پیشنهاد',
                    ),
                  ),
                if (incoming && offer.status == 'sent') ...[
                  FilledButton.icon(
                    onPressed: () => controller.respondOffer(offer, true),
                    icon: const Icon(Icons.check_outlined),
                    label: const Text('قبول پیشنهاد'),
                  ),
                  TextButton.icon(
                    onPressed: () => controller.respondOffer(offer, false),
                    icon: const Icon(Icons.close_outlined),
                    label: const Text('رد'),
                  ),
                ],
                if (!incoming && offer.status == 'requester_approved') ...[
                  FilledButton.icon(
                    onPressed: () => controller.finalizeOffer(offer, true),
                    icon: const Icon(Icons.verified_outlined),
                    label: const Text('تایید نهایی'),
                  ),
                  TextButton.icon(
                    onPressed: () => controller.finalizeOffer(offer, false),
                    icon: const Icon(Icons.block_outlined),
                    label: const Text('رد نهایی'),
                  ),
                ],
                if (offer.chatChannelId.isNotEmpty)
                  OutlinedButton.icon(
                    onPressed: () => Get.toNamed(AppRoutes.chats),
                    icon: const Icon(Icons.chat_bubble_outline),
                    label: const Text('چت فایل'),
                  ),
                OutlinedButton.icon(
                  onPressed: () => _history(context),
                  icon: const Icon(Icons.history_outlined),
                  label: const Text('تاریخچه'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _send(BuildContext context) {
    final field = TextEditingController(
      text: offer.commissionPercent.toStringAsFixed(0),
    );
    Get.dialog(
      AlertDialog(
        title: const Text('ارسال پیشنهاد'),
        content: TextField(
          controller: field,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(labelText: 'درصد سهم همکاری'),
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          FilledButton(
            onPressed: () {
              final value =
                  double.tryParse(field.text.trim()) ?? offer.commissionPercent;
              Get.back();
              controller.sendOffer(offer, value);
            },
            child: const Text('ارسال'),
          ),
        ],
      ),
    );
  }

  void _history(BuildContext context) {
    Get.bottomSheet(
      SafeArea(
        child: Material(
          borderRadius: const BorderRadius.vertical(top: Radius.circular(18)),
          child: ListView(
            shrinkWrap: true,
            padding: const EdgeInsets.all(16),
            children: [
              Text(
                'تاریخچه پیشنهاد',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 12),
              if (offer.history.isEmpty) const Text('تاریخچه‌ای ثبت نشده است'),
              ...offer.history.map(
                (item) => ListTile(
                  leading: const Icon(Icons.circle_outlined, size: 14),
                  title: Text(_historyAction(item.action)),
                  subtitle: Text(
                    item.note.isEmpty
                        ? item.createdAt
                        : '${item.note}\n${item.createdAt}',
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatusChip extends StatelessWidget {
  const _StatusChip({required this.status});

  final String status;

  @override
  Widget build(BuildContext context) {
    return Chip(label: Text(_statusText(status)));
  }
}

IconData _statusIcon(String status) => switch (status) {
  'approved' => Icons.verified_outlined,
  'rejected' => Icons.block_outlined,
  'requester_approved' => Icons.pending_actions_outlined,
  'sent' => Icons.outgoing_mail,
  _ => Icons.tips_and_updates_outlined,
};

String _statusText(String status) => switch (status) {
  'candidate' => 'قابل پیشنهاد',
  'sent' => 'ارسال شده',
  'requester_approved' => 'تایید گیرنده',
  'approved' => 'نهایی شده',
  'rejected' => 'رد شده',
  _ => status,
};

String _historyAction(String action) => switch (action) {
  'send' => 'ارسال پیشنهاد',
  'requester_approve' => 'تایید گیرنده',
  'requester_reject' => 'رد گیرنده',
  'owner_approve' => 'تایید نهایی',
  'owner_reject' => 'رد نهایی',
  _ => action,
};
