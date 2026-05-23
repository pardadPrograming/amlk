import 'dart:async';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';

import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../properties_controller.dart';

class LatestFilesPage extends StatefulWidget {
  const LatestFilesPage({super.key});

  @override
  State<LatestFilesPage> createState() => _LatestFilesPageState();
}

class _LatestFilesPageState extends State<LatestFilesPage> {
  final controller = Get.find<PropertiesController>();
  final scrollController = ScrollController();
  Timer? pollTimer;

  @override
  void initState() {
    super.initState();
    controller.loadLatestFiles(markRead: true);
    scrollController.addListener(_onScroll);
    pollTimer = Timer.periodic(
      const Duration(seconds: 20),
      (_) => controller.checkLatestUnread(),
    );
  }

  @override
  void dispose() {
    pollTimer?.cancel();
    scrollController
      ..removeListener(_onScroll)
      ..dispose();
    super.dispose();
  }

  void _onScroll() {
    if (!scrollController.hasClients) return;
    final position = scrollController.position;
    if (position.pixels >= position.maxScrollExtent - 120) {
      controller.loadLatestFiles(reset: false);
    }
    if (position.pixels <= 40 && controller.latestUnreadCount.value > 0) {
      controller.loadLatestFiles(markRead: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('جدیدترین فایل‌ها'),
      body: ResponsivePage(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const GradientHeader(
              title: 'جدیدترین فایل‌ها',
              subtitle: 'جمع‌بندی فایل‌هایی که برای شما ثبت یا دریافت شده‌اند',
              icon: Icons.dynamic_feed_outlined,
            ),
            const SizedBox(height: 12),
            Obx(
              () => _LatestFilterBar(
                selected: controller.latestFilter.value,
                onChanged: controller.setLatestFilter,
              ),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: Obx(() {
                if (controller.latestLoading.value &&
                    controller.latestFiles.isEmpty) {
                  return const Center(child: CircularProgressIndicator());
                }
                if (controller.latestFiles.isEmpty) {
                  return const AppCard(
                    child: Center(
                      child: Padding(
                        padding: EdgeInsets.all(20),
                        child: Text('هنوز فایلی برای نمایش وجود ندارد.'),
                      ),
                    ),
                  );
                }
                return Stack(
                  children: [
                    ListView.builder(
                      controller: scrollController,
                      reverse: true,
                      padding: const EdgeInsets.fromLTRB(0, 8, 0, 88),
                      itemCount:
                          controller.latestFiles.length +
                          (controller.latestLoadingMore.value ? 1 : 0),
                      itemBuilder: (context, index) {
                        if (index >= controller.latestFiles.length) {
                          return const Padding(
                            padding: EdgeInsets.all(18),
                            child: Center(child: CircularProgressIndicator()),
                          );
                        }
                        final file = controller.latestFiles[index];
                        return _LatestFileTile(
                          file: file,
                          unread: controller.isLatestFileUnread(file),
                        );
                      },
                    ),
                    if (controller.latestUnreadCount.value > 0)
                      PositionedDirectional(
                        bottom: 16,
                        end: 16,
                        child: FilledButton.icon(
                          onPressed: () =>
                              controller.loadLatestFiles(markRead: true),
                          icon: const Icon(Icons.keyboard_arrow_down_rounded),
                          label: Text(
                            '${controller.latestUnreadCount.value} فایل جدید',
                          ),
                        ),
                      ),
                  ],
                );
              }),
            ),
          ],
        ),
      ),
    );
  }
}

class _LatestFilterBar extends StatelessWidget {
  const _LatestFilterBar({required this.selected, required this.onChanged});

  final String selected;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    const filters = [
      _LatestFilter('', 'همه', Icons.all_inclusive_rounded),
      _LatestFilter('sale', 'فروش', Icons.sell_outlined),
      _LatestFilter('partnership', 'مشارکت', Icons.handshake_outlined),
      _LatestFilter('rent_lease', 'رهن و اجاره', Icons.key_outlined),
    ];
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: SegmentedButton<String>(
        selected: {selected},
        showSelectedIcon: false,
        onSelectionChanged: (value) => onChanged(value.first),
        segments: filters
            .map(
              (item) => ButtonSegment<String>(
                value: item.value,
                label: Text(item.label),
                icon: Icon(item.icon),
              ),
            )
            .toList(),
      ),
    );
  }
}

class _LatestFileTile extends StatelessWidget {
  const _LatestFileTile({required this.file, required this.unread});

  final PropertyFileModel file;
  final bool unread;

  @override
  Widget build(BuildContext context) {
    final address = file.addresses.isEmpty ? null : file.addresses.first;
    final colors = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: AppCard(
        padding: 14,
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            CircleAvatar(
              backgroundColor: colors.secondaryContainer,
              child: Icon(_fileIcon(file), color: colors.onSecondaryContainer),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          file.title,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: Theme.of(context).textTheme.titleMedium
                              ?.copyWith(fontWeight: FontWeight.w900),
                        ),
                      ),
                      if (unread)
                        Container(
                          width: 9,
                          height: 9,
                          decoration: BoxDecoration(
                            color: colors.primary,
                            shape: BoxShape.circle,
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 6),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: [
                      _MiniChip(
                        icon: _fileIcon(file),
                        label: (file.types.isEmpty ? [file.type] : file.types)
                            .map(_fileTypeLabel)
                            .join(' / '),
                      ),
                      _MiniChip(
                        icon: Icons.flag_outlined,
                        label: _statusLabel(file.status),
                      ),
                      if (file.isPartnershipCopy)
                        _MiniChip(
                          icon: Icons.handshake_outlined,
                          label:
                              'مشارکتی ${file.partnershipCommissionPercent.toStringAsFixed(0)}%',
                        ),
                      if (file.salePrice > 0)
                        _MiniChip(
                          icon: Icons.sell_outlined,
                          label: _money(file.salePrice),
                        ),
                      if (file.depositPrice > 0)
                        _MiniChip(
                          icon: Icons.account_balance_wallet_outlined,
                          label: 'رهن ${_money(file.depositPrice)}',
                        ),
                      if (file.rentPrice > 0)
                        _MiniChip(
                          icon: Icons.payments_outlined,
                          label: 'اجاره ${_money(file.rentPrice)}',
                        ),
                    ],
                  ),
                  if (address != null) ...[
                    const SizedBox(height: 8),
                    Text(
                      [
                        address.areaName,
                        address.streetName,
                        address.neighborhoodName,
                      ].where((item) => item.trim().isNotEmpty).join('، '),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                  if (file.description.isNotEmpty) ...[
                    const SizedBox(height: 8),
                    Text(
                      file.description,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: Theme.of(context).textTheme.bodySmall,
                    ),
                  ],
                  const SizedBox(height: 8),
                  Align(
                    alignment: AlignmentDirectional.centerEnd,
                    child: Text(
                      _relativeTime(file.createdAt),
                      style: Theme.of(context).textTheme.labelSmall?.copyWith(
                        color: Theme.of(context).hintColor,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _MiniChip extends StatelessWidget {
  const _MiniChip({required this.icon, required this.label});

  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Chip(
      visualDensity: VisualDensity.compact,
      avatar: Icon(icon, size: 16),
      label: Text(label),
    );
  }
}

class _LatestFilter {
  const _LatestFilter(this.value, this.label, this.icon);

  final String value;
  final String label;
  final IconData icon;
}

IconData _fileIcon(PropertyFileModel file) {
  final types = file.types.isEmpty ? [file.type] : file.types;
  if (types.contains('rent_lease')) return Icons.key_outlined;
  if (types.contains('partnership')) return Icons.handshake_outlined;
  return Icons.sell_outlined;
}

String _fileTypeLabel(String type) => switch (type) {
  'sale' => 'فروش',
  'partnership' => 'مشارکت',
  'rent_lease' => 'رهن و اجاره',
  _ => type,
};

String _statusLabel(String status) => switch (status) {
  'done' => 'انجام شده',
  'inactive' => 'غیرفعال',
  'suspended' => 'تعلیق شده',
  _ => 'فعال',
};

String _money(int value) {
  if (value <= 0) return '0';
  return NumberFormat.decimalPattern('fa_IR').format(value);
}

String _relativeTime(String value) {
  final date = DateTime.tryParse(value)?.toLocal();
  if (date == null) return '';
  final now = DateTime.now();
  final diff = now.difference(date);
  if (diff.inMinutes < 1) return 'همین الان';
  if (diff.inHours < 1) return '${diff.inMinutes} دقیقه پیش';
  if (diff.inDays < 1) return '${diff.inHours} ساعت پیش';
  if (diff.inDays < 7) return '${diff.inDays} روز پیش';
  return DateFormat('yyyy/MM/dd', 'fa_IR').format(date);
}
