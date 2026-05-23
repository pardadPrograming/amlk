import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';

import '../../../app/app.dart';
import '../../../app/app_theme.dart';
import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../../locations/locations_controller.dart';
import '../properties_controller.dart';

class PropertiesPage extends StatefulWidget {
  const PropertiesPage({super.key});

  @override
  State<PropertiesPage> createState() => _PropertiesPageState();
}

class _PropertiesPageState extends State<PropertiesPage> {
  final controller = Get.find<PropertiesController>();

  @override
  void initState() {
    super.initState();
    controller.load();
    controller.loadShareRequests();
    controller.loadNotifications();
    controller.ensureLocations();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text(
        'Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã™â€¦Ã™â€žÃšÂ©Ã›Å’',
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Get.to(() => const PropertyCreatePage()),
        icon: const Icon(Icons.add_home_work_outlined),
        label: const Text('Ã˜Â§Ã™ÂÃ˜Â²Ã™Ë†Ã˜Â¯Ã™â€  Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž'),
      ),
      body: ResponsivePage(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const GradientHeader(
              title: 'Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã™â€¦Ã™â€žÃšÂ©Ã›Å’',
              subtitle:
                  'Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã™ÂÃ˜Â±Ã™Ë†Ã˜Â´Ã˜Å’ Ã™â€¦Ã˜Â´Ã˜Â§Ã˜Â±ÃšÂ©Ã˜ÂªÃ˜Å’ Ã˜Â±Ã™â€¡Ã™â€  Ã™Ë† Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡ Ã˜Â±Ã˜Â§ Ã˜Â¨Ã˜Â§ Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³Ã¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â³Ã˜Â§Ã˜Â®Ã˜ÂªÃ˜Â§Ã˜Â±Ã›Å’Ã˜Â§Ã™ÂÃ˜ÂªÃ™â€¡ Ã™Ë† Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â§Ã›Å’ Ã˜Â¨Ã™â€¡Ã›Å’Ã™â€ Ã™â€¡ Ã˜Â«Ã˜Â¨Ã˜Âª ÃšÂ©Ã™â€ Ã›Å’Ã˜Â¯.',
              icon: Icons.apartment_outlined,
            ),
            const SizedBox(height: 14),
            AppCard(
              padding: 14,
              child: LayoutBuilder(
                builder: (context, constraints) {
                  final compact = constraints.maxWidth < 620;
                  final message = Row(
                    children: const [
                      Icon(Icons.map_outlined),
                      SizedBox(width: 12),
                      Expanded(
                        child: Text(
                          'Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã˜Â§Ã™â€ Ã˜ÂªÃ˜Â®Ã˜Â§Ã˜Â¨ Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§Ã˜Å’ Ã™â€¦Ã™â€ Ã˜Â§Ã˜Â·Ã™â€šÃ˜Å’ Ã˜Â®Ã›Å’Ã˜Â§Ã˜Â¨Ã˜Â§Ã™â€ Ã¢â‚¬Å’Ã™â€¡Ã˜Â§ Ã™Ë† Ã™â€¦Ã˜Â­Ã™â€žÃ™â€¡Ã¢â‚¬Å’Ã™â€¡Ã˜Â§ Ã˜Â±Ã˜Â§ Ã˜Â§Ã›Å’Ã™â€ Ã˜Â¬Ã˜Â§ Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â±Ã›Å’Ã˜Âª ÃšÂ©Ã™â€ Ã›Å’Ã˜Â¯.',
                        ),
                      ),
                    ],
                  );
                  final button = OutlinedButton.icon(
                    onPressed: () => Get.toNamed(AppRoutes.locations),
                    icon: const Icon(Icons.tune_outlined),
                    label: const Text(
                      'Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â±Ã›Å’Ã˜Âª Ã™â€¦Ã™â€ Ã˜Â§Ã˜Â·Ã™â€š',
                    ),
                  );
                  if (compact) {
                    return Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [message, const SizedBox(height: 12), button],
                    );
                  }
                  return Row(
                    children: [
                      Expanded(child: message),
                      const SizedBox(width: 12),
                      button,
                    ],
                  );
                },
              ),
            ),
            const SizedBox(height: 18),
            Obx(() => _ShareRequestsPanel(controller: controller)),
            const SizedBox(height: 18),
            Obx(() {
              if (controller.loading.value && controller.files.isEmpty) {
                return const Center(
                  child: Padding(
                    padding: EdgeInsets.all(32),
                    child: CircularProgressIndicator(),
                  ),
                );
              }
              if (controller.files.isEmpty) {
                return const AppCard(
                  child: Column(
                    children: [
                      Icon(Icons.apartment_outlined, size: 48),
                      SizedBox(height: 12),
                      Text(
                        'Ã™â€¡Ã™â€ Ã™Ë†Ã˜Â² Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã™â€¦Ã™â€žÃšÂ©Ã›Å’ Ã˜Â«Ã˜Â¨Ã˜Âª Ã™â€ Ã˜Â´Ã˜Â¯Ã™â€¡ Ã˜Â§Ã˜Â³Ã˜Âª.',
                        style: TextStyle(fontWeight: FontWeight.w900),
                      ),
                      SizedBox(height: 6),
                      Text(
                        'Ã˜Â§Ã˜Â² Ã˜Â¯ÃšÂ©Ã™â€¦Ã™â€¡ Ã˜Â§Ã™ÂÃ˜Â²Ã™Ë†Ã˜Â¯Ã™â€  Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â´Ã˜Â±Ã™Ë†Ã˜Â¹ ÃšÂ©Ã™â€ Ã›Å’Ã˜Â¯.',
                      ),
                    ],
                  ),
                );
              }
              return LayoutBuilder(
                builder: (context, constraints) {
                  final width = constraints.maxWidth < 760
                      ? constraints.maxWidth
                      : (constraints.maxWidth - 16) / 2;
                  return Wrap(
                    spacing: 16,
                    runSpacing: 16,
                    children: controller.files
                        .map(
                          (file) => SizedBox(
                            width: width,
                            child: _PropertyCard(file: file),
                          ),
                        )
                        .toList(),
                  );
                },
              );
            }),
          ],
        ),
      ),
    );
  }
}

class _PropertyCard extends StatelessWidget {
  const _PropertyCard({required this.file});

  final PropertyFileModel file;

  @override
  Widget build(BuildContext context) {
    final address = file.addresses.isEmpty ? null : file.addresses.first;
    return AppCard(
      padding: 16,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                _typeIcon(file.type),
                color: Theme.of(context).colorScheme.secondary,
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  file.title,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w900,
                  ),
                ),
              ),
              Chip(
                label: Text(
                  (file.types.isEmpty ? [file.type] : file.types)
                      .map(_typeLabel)
                      .join(' / '),
                ),
              ),
              IconButton(
                tooltip: 'صندوقچه‌ها',
                onPressed: _showVaultPicker,
                icon: const Icon(Icons.inventory_2_outlined),
              ),
            ],
          ),
          if (file.isPartnershipCopy) ...[
            const SizedBox(height: 8),
            Chip(
              avatar: const Icon(Icons.handshake_outlined, size: 17),
              label: Text(
                'مشارکت - سهم ${file.partnershipCommissionPercent.toStringAsFixed(0)}%',
              ),
            ),
          ],
          const SizedBox(height: 10),
          if (address != null)
            Text(
              '${address.areaName}Ã˜Å’ ${address.streetName}Ã˜Å’ ${address.neighborhoodName}',
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          const SizedBox(height: 10),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _InfoChip(
                icon: Icons.flag_outlined,
                label: _propertyStatusLabel(file.status),
              ),
              if (file.vaultPlacements.isNotEmpty)
                _InfoChip(
                  icon: Icons.percent_rounded,
                  label:
                      'سهم همکار ${file.vaultPlacements.map((e) => e.commissionPercent.toStringAsFixed(0)).join('/')}%',
                ),
              _InfoChip(
                icon: Icons.image_outlined,
                label: '${file.media.length} Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â§',
              ),
              if (file.salePrice > 0)
                _InfoChip(
                  icon: Icons.sell_outlined,
                  label: '${file.salePrice}',
                ),
              if (file.depositPrice > 0)
                _InfoChip(
                  icon: Icons.account_balance_wallet_outlined,
                  label: 'Ã˜Â±Ã™â€¡Ã™â€  ${file.depositPrice}',
                ),
              if (file.rentPrice > 0)
                _InfoChip(
                  icon: Icons.payments_outlined,
                  label: 'Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡ ${file.rentPrice}',
                ),
            ],
          ),
          if (file.description.isNotEmpty) ...[
            const SizedBox(height: 10),
            _DescriptionPreview(
              icon: Icons.campaign_outlined,
              label: 'توضیح مشتری',
              text: file.description,
            ),
          ],
          if (file.internalDescription.isNotEmpty) ...[
            const SizedBox(height: 8),
            _DescriptionPreview(
              icon: Icons.lock_outline_rounded,
              label: 'یادداشت مشاور',
              text: file.internalDescription,
              muted: true,
            ),
          ],
          const SizedBox(height: 10),
          Align(
            alignment: Alignment.centerLeft,
            child: OutlinedButton.icon(
              onPressed: file.isPartnershipCopy ? null : _requestShare,
              icon: const Icon(Icons.handshake_outlined),
              label: const Text('درخواست مشارکت'),
            ),
          ),
          if (file.sharingHistory.isNotEmpty) ...[
            const Divider(height: 20),
            Text(
              'تاریخچه اشتراک‌گذاری',
              style: Theme.of(
                context,
              ).textTheme.labelLarge?.copyWith(fontWeight: FontWeight.w900),
            ),
            const SizedBox(height: 6),
            ...file.sharingHistory
                .take(3)
                .map(
                  (item) => ListTile(
                    dense: true,
                    contentPadding: EdgeInsets.zero,
                    leading: const Icon(Icons.history_outlined),
                    title: Text(
                      '${item.userName.isEmpty ? item.userPhone : item.userName} - ${item.commissionPercent.toStringAsFixed(0)}%',
                    ),
                    subtitle: Text(_shareStatusLabel(item.status)),
                  ),
                ),
          ],
        ],
      ),
    );
  }

  void _requestShare() {
    final percent = TextEditingController(text: '25');
    Get.dialog(
      AlertDialog(
        title: const Text('درخواست مشارکت'),
        content: TextField(
          controller: percent,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(
            labelText: 'درصد سهم مشارکت',
            suffixText: '%',
            prefixIcon: Icon(Icons.percent_rounded),
          ),
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          FilledButton.icon(
            onPressed: () {
              Get.find<PropertiesController>().requestShare(
                file,
                percent: double.tryParse(percent.text) ?? 25,
              );
              Get.back();
            },
            icon: const Icon(Icons.send_outlined),
            label: const Text('ارسال'),
          ),
        ],
      ),
    );
  }

  Future<void> _showVaultPicker() async {
    final controller = Get.find<PropertiesController>();
    await controller.loadVaults();
    final selected = _initialVaultCommissions(file);
    final status = file.status.obs;
    Get.dialog(
      AlertDialog(
        title: const Text('صندوقچه‌های فایل'),
        content: SizedBox(
          width: 520,
          child: Obx(
            () => Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                DropdownButtonFormField<String>(
                  initialValue: status.value,
                  decoration: const InputDecoration(
                    labelText: 'وضعیت فایل',
                    prefixIcon: Icon(Icons.flag_outlined),
                  ),
                  items: _propertyStatusOptions
                      .map(
                        (item) => DropdownMenuItem(
                          value: item.value,
                          child: Text(item.label),
                        ),
                      )
                      .toList(),
                  onChanged: (value) {
                    if (value != null) status.value = value;
                  },
                ),
                const SizedBox(height: 12),
                _VaultSelector(
                  vaults: controller.vaults,
                  selectedVaultCommissions: selected,
                ),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          FilledButton.icon(
            onPressed: () async {
              await controller.updateVaults(
                file,
                _vaultPlacements(selected),
                status.value,
              );
              Get.back();
            },
            icon: const Icon(Icons.check_rounded),
            label: const Text('ذخیره'),
          ),
        ],
      ),
    );
  }
}

RxMap<String, double> _initialVaultCommissions(PropertyFileModel file) {
  final result = <String, double>{};
  if (file.vaultPlacements.isNotEmpty) {
    for (final placement in file.vaultPlacements) {
      if (placement.vaultId.isNotEmpty) {
        result[placement.vaultId] = placement.commissionPercent;
      }
    }
  } else {
    for (final vaultId in file.vaultIds) {
      if (vaultId.isNotEmpty) {
        result[vaultId] = 25;
      }
    }
  }
  return result.obs;
}

List<PropertyVaultPlacementModel> _vaultPlacements(
  Map<String, double> commissions,
) => commissions.entries
    .where((entry) => entry.key.trim().isNotEmpty)
    .map(
      (entry) => PropertyVaultPlacementModel(
        vaultId: entry.key,
        commissionPercent: entry.value,
      ),
    )
    .toList();

class _PropertyStatusOption {
  const _PropertyStatusOption(this.value, this.label);

  final String value;
  final String label;
}

const _propertyStatusOptions = [
  _PropertyStatusOption('active', 'فعال'),
  _PropertyStatusOption('done', 'انجام شده'),
  _PropertyStatusOption('suspended', 'تعلیق شده'),
  _PropertyStatusOption('expired', 'اتمام شده'),
];

String _propertyStatusLabel(String value) => _propertyStatusOptions
    .firstWhere(
      (item) => item.value == value,
      orElse: () => _propertyStatusOptions.first,
    )
    .label;

class _ShareRequestsPanel extends StatelessWidget {
  const _ShareRequestsPanel({required this.controller});

  final PropertiesController controller;

  @override
  Widget build(BuildContext context) {
    final incoming = controller.ownerShareRequests;
    final outgoing = controller.requesterShareRequests;
    if (incoming.isEmpty && outgoing.isEmpty) {
      return const SizedBox.shrink();
    }
    return AppCard(
      padding: 14,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'درخواست‌های مشارکت',
            style: Theme.of(
              context,
            ).textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w900),
          ),
          const SizedBox(height: 8),
          ...incoming
              .take(5)
              .map(
                (request) => ListTile(
                  contentPadding: EdgeInsets.zero,
                  leading: const Icon(Icons.inbox_outlined),
                  title: Text(request.propertyTitle),
                  subtitle: Text(
                    '${request.requesterName.isEmpty ? request.requesterPhone : request.requesterName} - ${request.commissionPercent.toStringAsFixed(0)}% - ${_shareStatusLabel(request.status)}',
                  ),
                  trailing: request.status == 'pending'
                      ? Wrap(
                          spacing: 6,
                          children: [
                            IconButton(
                              tooltip: 'تایید',
                              onPressed: () =>
                                  controller.decideShare(request, true),
                              icon: const Icon(Icons.check_circle_outline),
                            ),
                            IconButton(
                              tooltip: 'رد',
                              onPressed: () =>
                                  controller.decideShare(request, false),
                              icon: const Icon(Icons.cancel_outlined),
                            ),
                          ],
                        )
                      : null,
                ),
              ),
          ...outgoing
              .take(5)
              .map(
                (request) => ListTile(
                  contentPadding: EdgeInsets.zero,
                  leading: const Icon(Icons.outbox_outlined),
                  title: Text(request.propertyTitle),
                  subtitle: Text(
                    'درخواست شما - ${request.commissionPercent.toStringAsFixed(0)}% - ${_shareStatusLabel(request.status)}',
                  ),
                  trailing: request.status == 'approved'
                      ? FilledButton.icon(
                          onPressed: () =>
                              controller.receiveSharedFile(request),
                          icon: const Icon(Icons.file_open_outlined),
                          label: const Text('دریافت فایل'),
                        )
                      : null,
                ),
              ),
        ],
      ),
    );
  }
}

String _shareStatusLabel(String value) => switch (value) {
  'pending' => 'در انتظار',
  'approved' => 'تایید شده',
  'rejected' => 'رد شده',
  'received' => 'دریافت شده',
  _ => value,
};

class _InfoChip extends StatelessWidget {
  const _InfoChip({required this.icon, required this.label});

  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Chip(
      avatar: Icon(icon, size: 17),
      label: Text(label),
      visualDensity: VisualDensity.compact,
    );
  }
}

class _DescriptionPreview extends StatelessWidget {
  const _DescriptionPreview({
    required this.icon,
    required this.label,
    required this.text,
    this.muted = false,
  });

  final IconData icon;
  final String label;
  final String text;
  final bool muted;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(10),
      decoration: BoxDecoration(
        color:
            (muted ? scheme.surfaceContainerHighest : scheme.primaryContainer)
                .withValues(alpha: 0.35),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(
          color: scheme.outlineVariant.withValues(alpha: 0.45),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                icon,
                size: 16,
                color: muted ? scheme.outline : scheme.primary,
              ),
              const SizedBox(width: 6),
              Text(
                label,
                style: TextStyle(
                  fontWeight: FontWeight.w900,
                  color: muted ? scheme.outline : scheme.primary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 6),
          Text(text, maxLines: 3, overflow: TextOverflow.ellipsis),
        ],
      ),
    );
  }
}

class PropertyCreatePage extends StatefulWidget {
  const PropertyCreatePage({super.key});

  @override
  State<PropertyCreatePage> createState() => _PropertyCreatePageState();
}

class _PropertyCreatePageState extends State<PropertyCreatePage> {
  final controller = Get.find<PropertiesController>();
  final locations = Get.find<LocationsController>();
  final formKey = GlobalKey<FormState>();
  final title = TextEditingController();
  final description = TextEditingController();
  final internalDescription = TextEditingController();
  final salePrice = TextEditingController();
  final finalPrice = TextEditingController();
  final depositPrice = TextEditingController();
  final rentPrice = TextEditingController();
  final maxConvertibleDeposit = TextEditingController();
  final areaM2 = TextEditingController();
  final bedrooms = TextEditingController();
  final floor = TextEditingController();
  final totalFloors = TextEditingController();
  final ageYears = TextEditingController();
  final terraceCount = TextEditingController();
  final backyardAreaM2 = TextEditingController();
  final gardenBuildingAreaM2 = TextEditingController();
  final gardenBuildingFloors = TextEditingController();
  final masterServiceCount = TextEditingController();
  final parking = false.obs;
  final elevator = false.obs;
  final storage = false.obs;
  final renovated = false.obs;
  final separateEntrance = false.obs;
  final rentWithOwner = false.obs;
  final gardenBuilding = false.obs;
  final pool = false.obs;
  final waterUtility = false.obs;
  final electricityUtility = false.obs;
  final gasUtility = false.obs;
  final waterRight = false.obs;
  final permit = false.obs;
  final stoveTop = false.obs;
  final parkingDisturbance = false.obs;
  final parkingInDeed = false.obs;
  final parkingCommon = false.obs;
  final doubleGlazedWindows = false.obs;
  final convertible = false.obs;
  final fileTypes = <String>['sale'].obs;
  final propertyStatus = 'active'.obs;
  final propertyType = 'apartment'.obs;
  final documentType = 'six_dang'.obs;
  final directions = <String>[].obs;
  final cooling = <String>[].obs;
  final flooring = 'ceramic'.obs;
  final heating = 'radiator'.obs;
  final cabinetType = 'mdf'.obs;
  final wallCovering = 'paint'.obs;
  final kitchenType = 'open'.obs;
  final addresses = <_AddressDraft>[].obs;
  final selectedVaultCommissions = <String, double>{}.obs;

  @override
  void initState() {
    super.initState();
    addresses.add(_AddressDraft());
    locations.load();
    controller.loadVaults();
  }

  @override
  void dispose() {
    title.dispose();
    description.dispose();
    internalDescription.dispose();
    salePrice.dispose();
    finalPrice.dispose();
    depositPrice.dispose();
    rentPrice.dispose();
    maxConvertibleDeposit.dispose();
    areaM2.dispose();
    bedrooms.dispose();
    floor.dispose();
    totalFloors.dispose();
    ageYears.dispose();
    terraceCount.dispose();
    backyardAreaM2.dispose();
    gardenBuildingAreaM2.dispose();
    gardenBuildingFloors.dispose();
    masterServiceCount.dispose();
    for (final item in addresses) {
      item.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('Ã˜Â§Ã™ÂÃ˜Â²Ã™Ë†Ã˜Â¯Ã™â€  Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž'),
      body: ResponsivePage(
        child: SingleChildScrollView(
          child: Form(
            key: formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const GradientHeader(
                  title: 'Ã˜Â«Ã˜Â¨Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â¬Ã˜Â¯Ã›Å’Ã˜Â¯',
                  subtitle:
                      'Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ Ã˜Â¨Ã˜Â§Ã›Å’Ã˜Â¯ Ã˜Â§Ã˜Â² Ã™â€¦Ã™â€ Ã˜Â§Ã˜Â·Ã™â€š Ã˜ÂªÃ˜Â¹Ã˜Â±Ã›Å’Ã™ÂÃ¢â‚¬Å’Ã˜Â´Ã˜Â¯Ã™â€¡ Ã˜Â§Ã™â€ Ã˜ÂªÃ˜Â®Ã˜Â§Ã˜Â¨ Ã˜Â´Ã™Ë†Ã˜Â¯ Ã™Ë† Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ Ã˜Â¯Ã™â€šÃ›Å’Ã™â€š Ã˜Â¯Ã˜Â³Ã˜ÂªÃ›Å’ Ã™â€¡Ã™â€¦ Ã™â€šÃ˜Â§Ã˜Â¨Ã™â€ž Ã˜Â«Ã˜Â¨Ã˜Âª Ã˜Â§Ã˜Â³Ã˜Âª.',
                  icon: Icons.add_home_work_outlined,
                ),
                const SizedBox(height: 18),
                AppCard(
                  child: Column(
                    children: [
                      Obx(
                        () => _FileTypeSelector(
                          types: fileTypes,
                          propertyType: propertyType,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Obx(
                        () => DropdownButtonFormField<String>(
                          initialValue: propertyStatus.value,
                          decoration: const InputDecoration(
                            labelText: 'وضعیت فایل',
                            prefixIcon: Icon(Icons.flag_outlined),
                          ),
                          items: _propertyStatusOptions
                              .map(
                                (item) => DropdownMenuItem(
                                  value: item.value,
                                  child: Text(item.label),
                                ),
                              )
                              .toList(),
                          onChanged: (value) {
                            if (value != null) {
                              propertyStatus.value = value;
                            }
                          },
                        ),
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: title,
                        decoration: const InputDecoration(
                          labelText: 'Ã˜Â¹Ã™â€ Ã™Ë†Ã˜Â§Ã™â€  Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž',
                        ),
                        textInputAction: TextInputAction.next,
                        validator: (value) =>
                            value == null || value.trim().isEmpty
                            ? 'Ã˜Â¹Ã™â€ Ã™Ë†Ã˜Â§Ã™â€  Ã˜Â§Ã™â€žÃ˜Â²Ã˜Â§Ã™â€¦Ã›Å’ Ã˜Â§Ã˜Â³Ã˜Âª'
                            : null,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: description,
                        decoration: const InputDecoration(
                          labelText: 'توضیحات قابل نمایش به مشتری',
                          alignLabelWithHint: true,
                          prefixIcon: Icon(Icons.campaign_outlined),
                        ),
                        textAlignVertical: TextAlignVertical.top,
                        minLines: 4,
                        maxLines: 8,
                      ),
                      const SizedBox(height: 12),
                      TextFormField(
                        controller: internalDescription,
                        decoration: const InputDecoration(
                          labelText: 'یادداشت داخلی مشاور',
                          hintText:
                              'نکات مذاکره، حساسیت مالک، وضعیت فایل یا توضیحاتی که برای مشتری نمایش داده نمی‌شود',
                          alignLabelWithHint: true,
                          prefixIcon: Icon(Icons.lock_outline_rounded),
                        ),
                        textAlignVertical: TextAlignVertical.top,
                        minLines: 4,
                        maxLines: 8,
                      ),
                      const SizedBox(height: 12),
                      Obx(
                        () => _PriceFields(
                          types: fileTypes,
                          propertyType: propertyType,
                          salePrice: salePrice,
                          finalPrice: finalPrice,
                          depositPrice: depositPrice,
                          rentPrice: rentPrice,
                          convertible: convertible,
                          maxConvertibleDeposit: maxConvertibleDeposit,
                          rentWithOwner: rentWithOwner,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 14),
                AppCard(
                  child: _HouseInfoForm(
                    propertyType: propertyType,
                    areaM2: areaM2,
                    bedrooms: bedrooms,
                    floor: floor,
                    totalFloors: totalFloors,
                    ageYears: ageYears,
                    terraceCount: terraceCount,
                    backyardAreaM2: backyardAreaM2,
                    gardenBuildingAreaM2: gardenBuildingAreaM2,
                    gardenBuildingFloors: gardenBuildingFloors,
                    documentType: documentType,
                    parking: parking,
                    elevator: elevator,
                    storage: storage,
                    renovated: renovated,
                    separateEntrance: separateEntrance,
                    stoveTop: stoveTop,
                    parkingDisturbance: parkingDisturbance,
                    parkingInDeed: parkingInDeed,
                    parkingCommon: parkingCommon,
                    masterServiceCount: masterServiceCount,
                    doubleGlazedWindows: doubleGlazedWindows,
                    directions: directions,
                    cooling: cooling,
                    flooring: flooring,
                    heating: heating,
                    cabinetType: cabinetType,
                    wallCovering: wallCovering,
                    kitchenType: kitchenType,
                    gardenBuilding: gardenBuilding,
                    pool: pool,
                    waterUtility: waterUtility,
                    electricityUtility: electricityUtility,
                    gasUtility: gasUtility,
                    waterRight: waterRight,
                    permit: permit,
                  ),
                ),
                const SizedBox(height: 14),
                AppCard(
                  child: Obx(
                    () => _VaultSelector(
                      vaults: controller.vaults,
                      selectedVaultCommissions: selectedVaultCommissions,
                    ),
                  ),
                ),
                const SizedBox(height: 14),
                AppCard(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const Expanded(
                            child: Text(
                              'Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³Ã¢â‚¬Å’Ã™â€¡Ã˜Â§',
                              style: TextStyle(fontWeight: FontWeight.w900),
                            ),
                          ),
                          Obx(
                            () => IconButton(
                              tooltip:
                                  'Ã˜Â§Ã™ÂÃ˜Â²Ã™Ë†Ã˜Â¯Ã™â€  Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³',
                              onPressed: addresses.length >= 5
                                  ? null
                                  : () => addresses.add(_AddressDraft()),
                              icon: const Icon(Icons.add_location_alt_outlined),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Obx(
                        () => Column(
                          children: addresses
                              .asMap()
                              .entries
                              .map(
                                (entry) => _AddressEditor(
                                  index: entry.key,
                                  draft: entry.value,
                                  canRemove: addresses.length > 1,
                                  onRemove: () {
                                    entry.value.dispose();
                                    addresses.remove(entry.value);
                                  },
                                ),
                              )
                              .toList(),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 14),
                AppCard(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const Expanded(
                            child: Text(
                              'Ã˜Â¹ÃšÂ©Ã˜Â³ Ã™Ë† Ã™Ë†Ã›Å’Ã˜Â¯Ã˜Â¦Ã™Ë†',
                              style: TextStyle(fontWeight: FontWeight.w900),
                            ),
                          ),
                          OutlinedButton.icon(
                            onPressed: controller.pickMedia,
                            icon: const Icon(Icons.upload_file_outlined),
                            label: const Text(
                              'Ã˜Â§Ã™â€ Ã˜ÂªÃ˜Â®Ã˜Â§Ã˜Â¨ Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â§',
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      const Text(
                        'Ã˜Â­Ã˜Â¯Ã˜Â§ÃšÂ©Ã˜Â«Ã˜Â± Ã›Â²Ã›Â° Ã™â€¦Ã˜Â¯Ã›Å’Ã˜Â§Ã˜Å’ Ã˜Â­Ã˜Â¯Ã˜Â§ÃšÂ©Ã˜Â«Ã˜Â± Ã›Â² Ã™Ë†Ã›Å’Ã˜Â¯Ã˜Â¦Ã™Ë†. Ã˜Â¹ÃšÂ©Ã˜Â³Ã¢â‚¬Å’Ã™â€¡Ã˜Â§ Ã˜Â±Ã™Ë†Ã›Å’ Ã˜Â³Ã˜Â±Ã™Ë†Ã˜Â± Ã˜ÂªÃ˜Â§ Ã›ÂµÃ›Â°Ã›Â°KB Ã™Ë† Ã™Ë†Ã›Å’Ã˜Â¯Ã˜Â¦Ã™Ë†Ã™â€¡Ã˜Â§ Ã˜Â¨Ã™â€¡ 480p Ã˜ÂªÃ˜Â¨Ã˜Â¯Ã›Å’Ã™â€ž Ã™â€¦Ã›Å’Ã¢â‚¬Å’Ã˜Â´Ã™Ë†Ã™â€ Ã˜Â¯.',
                      ),
                      const SizedBox(height: 10),
                      Obx(
                        () => Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: controller.selectedUploads
                              .map(
                                (file) => InputChip(
                                  avatar: Icon(_uploadIcon(file.name)),
                                  label: Text(
                                    file.name,
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                  onDeleted: () =>
                                      controller.removeUpload(file),
                                ),
                              )
                              .toList(),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 18),
                Obx(
                  () => PremiumActionButton(
                    onPressed: controller.loading.value ? null : _submit,
                    icon: Icons.check_rounded,
                    child: controller.loading.value
                        ? const SizedBox(
                            width: 22,
                            height: 22,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : const Text('Ã˜Â«Ã˜Â¨Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _submit() async {
    if (!(formKey.currentState?.validate() ?? false)) return;
    if (propertyType.value == 'apartment' &&
        fileTypes.contains('partnership')) {
      fileTypes.remove('partnership');
    }
    final builtAddresses = <PropertyAddressModel>[];
    for (final draft in addresses) {
      if (draft.areaId.value.isEmpty ||
          draft.streetId.value.isEmpty ||
          draft.neighborhoodId.value.isEmpty) {
        Get.snackbar(
          'Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ Ã™â€ Ã˜Â§Ã™â€šÃ˜Âµ',
          'Ã™â€¦Ã™â€ Ã˜Â·Ã™â€šÃ™â€¡Ã˜Å’ Ã˜Â®Ã›Å’Ã˜Â§Ã˜Â¨Ã˜Â§Ã™â€  Ã™Ë† Ã™â€¦Ã˜Â­Ã™â€žÃ™â€¡ Ã˜Â±Ã˜Â§ Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã™â€¡Ã™â€¦Ã™â€¡ Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³Ã¢â‚¬Å’Ã™â€¡Ã˜Â§ Ã˜Â§Ã™â€ Ã˜ÂªÃ˜Â®Ã˜Â§Ã˜Â¨ ÃšÂ©Ã™â€ Ã›Å’Ã˜Â¯.',
        );
        return;
      }
      builtAddresses.add(
        PropertyAddressModel(
          areaId: draft.areaId.value,
          streetId: draft.streetId.value,
          neighborhoodId: draft.neighborhoodId.value,
          manualExactAddress: draft.manual.text.trim(),
        ),
      );
    }
    await controller.create(
      type: fileTypes.first,
      status: propertyStatus.value,
      types: fileTypes.toList(),
      title: title.text.trim(),
      description: description.text.trim(),
      internalDescription: internalDescription.text.trim(),
      salePrice: _number(salePrice.text),
      finalPrice: _number(finalPrice.text),
      depositPrice: _number(depositPrice.text),
      rentPrice: _number(rentPrice.text),
      convertible: fileTypes.contains('rent_lease') && convertible.value,
      maxConvertibleDeposit: fileTypes.contains('rent_lease')
          ? _number(maxConvertibleDeposit.text)
          : 0,
      rentWithOwner:
          fileTypes.contains('rent_lease') &&
          propertyType.value == 'apartment' &&
          rentWithOwner.value,
      houseInfo: {
        'propertyType': propertyType.value,
        'areaM2': _number(areaM2.text),
        'bedrooms': _number(bedrooms.text),
        'floor': _number(floor.text),
        'totalFloors': _number(totalFloors.text),
        'ageYears': _number(ageYears.text),
        'renovated': renovated.value,
        'parking': parking.value,
        'elevator': elevator.value,
        'storage': storage.value,
        'terraceCount': _number(terraceCount.text),
        'backyardAreaM2': _number(backyardAreaM2.text),
        'separateEntrance':
            propertyType.value != 'villa' && separateEntrance.value,
        'documentType': _documentTypeLabel(documentType.value),
        'directions': directions
            .map((value) => _optionLabel(_directionOptions, value))
            .toList(),
        'flooring': _optionLabel(_flooringOptions, flooring.value),
        'heating': _optionLabel(_heatingOptions, heating.value),
        'cabinetType': _optionLabel(_cabinetOptions, cabinetType.value),
        'cooling': cooling
            .map((value) => _optionLabel(_coolingOptions, value))
            .toList(),
        'wallCovering': _optionLabel(_wallCoveringOptions, wallCovering.value),
        'kitchenType': _optionLabel(_kitchenOptions, kitchenType.value),
        'stoveTop': stoveTop.value,
        'parkingDisturbance': parkingDisturbance.value,
        'parkingInDeed': parkingInDeed.value,
        'parkingCommon': parkingCommon.value,
        'masterServiceCount': _number(masterServiceCount.text),
        'doubleGlazedWindows': doubleGlazedWindows.value,
        'gardenBuilding':
            propertyType.value == 'garden' && gardenBuilding.value,
        'gardenBuildingAreaM2':
            propertyType.value == 'garden' && gardenBuilding.value
            ? _number(gardenBuildingAreaM2.text)
            : 0,
        'gardenBuildingFloors':
            propertyType.value == 'garden' && gardenBuilding.value
            ? _number(gardenBuildingFloors.text)
            : 0,
        'pool': propertyType.value == 'garden' && pool.value,
        'waterUtility':
            (propertyType.value == 'garden' || propertyType.value == 'land') &&
            waterUtility.value,
        'electricityUtility':
            (propertyType.value == 'garden' || propertyType.value == 'land') &&
            electricityUtility.value,
        'gasUtility':
            (propertyType.value == 'garden' || propertyType.value == 'land') &&
            gasUtility.value,
        'waterRight': propertyType.value == 'garden' && waterRight.value,
        'permit':
            (propertyType.value == 'garden' || propertyType.value == 'land') &&
            permit.value,
        'landType': propertyType.value == 'land' ? 'unspecified' : '',
      },
      addresses: builtAddresses,
      vaultPlacements: _vaultPlacements(selectedVaultCommissions),
    );
  }
}

class _VaultSelector extends StatelessWidget {
  const _VaultSelector({
    required this.vaults,
    required this.selectedVaultCommissions,
  });

  final List<ChannelVaultModel> vaults;
  final RxMap<String, double> selectedVaultCommissions;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        Row(
          children: [
            Icon(
              Icons.inventory_2_outlined,
              color: theme.colorScheme.secondary,
            ),
            const SizedBox(width: 8),
            Expanded(
              child: Text(
                'صندوقچه‌های فایل',
                style: theme.textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.w900,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text(
          'می‌توانید فایل را همزمان در چند صندوقچه شخصی یا صندوقچه املاک قرار دهید.',
          style: theme.textTheme.bodySmall?.copyWith(color: theme.hintColor),
        ),
        const SizedBox(height: 10),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [
            OutlinedButton.icon(
              onPressed: () => _promptCreateVault(business: false),
              icon: const Icon(Icons.person_add_alt_1_outlined),
              label: const Text('صندوقچه شخصی جدید'),
            ),
            OutlinedButton.icon(
              onPressed: () => _promptCreateVault(business: true),
              icon: const Icon(Icons.add_business_outlined),
              label: const Text('صندوقچه املاک جدید'),
            ),
          ],
        ),
        const SizedBox(height: 12),
        if (vaults.isEmpty)
          Text(
            'هنوز صندوقچه‌ای برای انتخاب آماده نیست.',
            style: theme.textTheme.bodyMedium,
          )
        else
          Column(
            children: vaults.map((vault) {
              final selected = selectedVaultCommissions.containsKey(vault.id);
              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Row(
                  children: [
                    Expanded(
                      child: CheckboxListTile(
                        contentPadding: EdgeInsets.zero,
                        controlAffinity: ListTileControlAffinity.leading,
                        secondary: Icon(
                          vault.isBusinessVault
                              ? Icons.business_center_outlined
                              : Icons.person_outline_rounded,
                        ),
                        title: Text(_vaultTitle(vault)),
                        subtitle: const Text('سهم پیش‌فرض همکار ۲۵٪ است'),
                        value: selected,
                        onChanged: (value) {
                          if (value == true) {
                            selectedVaultCommissions[vault.id] =
                                selectedVaultCommissions[vault.id] ?? 25;
                          } else {
                            selectedVaultCommissions.remove(vault.id);
                          }
                        },
                      ),
                    ),
                    SizedBox(
                      width: 112,
                      child: TextFormField(
                        enabled: selected,
                        initialValue: (selectedVaultCommissions[vault.id] ?? 25)
                            .toStringAsFixed(0),
                        decoration: const InputDecoration(
                          labelText: 'درصد',
                          suffixText: '%',
                        ),
                        keyboardType: TextInputType.number,
                        onChanged: (value) {
                          final parsed = double.tryParse(value) ?? 0;
                          selectedVaultCommissions[vault.id] = parsed;
                        },
                      ),
                    ),
                  ],
                ),
              );
            }).toList(),
          ),
      ],
    );
  }

  String _vaultTitle(ChannelVaultModel vault) {
    final title = vault.title.trim().isEmpty ? 'صندوقچه' : vault.title.trim();
    if (vault.isMain && vault.isBusinessVault) return '$title اصلی املاک';
    if (vault.isMain) return '$title اصلی';
    return title;
  }

  void _promptCreateVault({required bool business}) {
    final controller = Get.find<PropertiesController>();
    final title = TextEditingController();
    Get.dialog(
      AlertDialog(
        title: Text(business ? 'صندوقچه املاک جدید' : 'صندوقچه شخصی جدید'),
        content: TextField(
          controller: title,
          autofocus: true,
          decoration: const InputDecoration(
            labelText: 'نام صندوقچه',
            prefixIcon: Icon(Icons.inventory_2_outlined),
          ),
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          FilledButton.icon(
            onPressed: () async {
              if (business) {
                await controller.createBusinessVault(title.text);
              } else {
                await controller.createUserVault(title.text);
              }
              Get.back();
            },
            icon: const Icon(Icons.check_rounded),
            label: const Text('ساخت'),
          ),
        ],
      ),
    );
  }
}

class _FileTypeSelector extends StatelessWidget {
  const _FileTypeSelector({required this.types, required this.propertyType});

  final RxList<String> types;
  final RxString propertyType;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final accent = isDark ? AppColors.electricCyan : AppColors.secondary;
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topRight,
          end: Alignment.bottomLeft,
          colors: isDark
              ? const [Color(0xFF13263D), Color(0xFF0C1728)]
              : const [Color(0xFFFFFFFF), Color(0xFFEAF2FF)],
        ),
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: accent.withValues(alpha: 0.28)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'ÙˆØ¶Ø¹ÛŒØªâ€ŒÙ‡Ø§ÛŒ ÙØ§ÛŒÙ„',
            style: Theme.of(context).textTheme.labelMedium?.copyWith(
              color: Theme.of(context).hintColor,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 10),
          Obx(
            () => Wrap(
              spacing: 8,
              runSpacing: 8,
              children: _fileTypeOptions.map((option) {
                final selected = types.contains(option.value);
                final disabled =
                    propertyType.value == 'apartment' &&
                    option.value == 'partnership';
                return FilterChip(
                  avatar: Icon(option.icon, size: 18),
                  label: Text(option.title),
                  selected: selected,
                  disabledColor: Theme.of(
                    context,
                  ).disabledColor.withValues(alpha: 0.08),
                  onSelected: disabled
                      ? null
                      : (value) {
                          if (value) {
                            types.add(option.value);
                          } else if (types.length > 1) {
                            types.remove(option.value);
                          }
                        },
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _FileTypeOption {
  const _FileTypeOption({
    required this.value,
    required this.title,
    required this.subtitle,
    required this.icon,
  });

  final String value;
  final String title;
  final String subtitle;
  final IconData icon;
}

const _fileTypeOptions = [
  _FileTypeOption(
    value: 'sale',
    title: 'Ã™ÂÃ˜Â±Ã™Ë†Ã˜Â´',
    subtitle:
        'Ã™â€šÃ›Å’Ã™â€¦Ã˜Âª Ã™ÂÃ˜Â±Ã™Ë†Ã˜Â´ Ã™Ë† Ã™â€šÃ›Å’Ã™â€¦Ã˜Âª Ã™â€ Ã™â€¡Ã˜Â§Ã›Å’Ã›Å’ Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â«Ã˜Â¨Ã˜Âª Ã™â€¦Ã›Å’Ã¢â‚¬Å’Ã˜Â´Ã™Ë†Ã˜Â¯.',
    icon: Icons.sell_outlined,
  ),
  _FileTypeOption(
    value: 'partnership',
    title: 'Ã™â€¦Ã˜Â´Ã˜Â§Ã˜Â±ÃšÂ©Ã˜Âª',
    subtitle:
        'Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â³Ã˜Â§Ã˜Â®Ã˜ÂªÃ¢â‚¬Å’Ã™Ë†Ã˜Â³Ã˜Â§Ã˜Â² Ã™Ë† Ã™â€¦Ã˜Â´Ã˜Â§Ã˜Â±ÃšÂ©Ã˜Âª Ã˜Â¯Ã˜Â± Ã˜Â³Ã˜Â§Ã˜Â®Ã˜Âª.',
    icon: Icons.handshake_outlined,
  ),
  _FileTypeOption(
    value: 'rent_lease',
    title: 'Ã˜Â±Ã™â€¡Ã™â€  Ã™Ë† Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡',
    subtitle:
        'Ã™â€¦Ã˜Â¨Ã™â€žÃ˜Âº Ã˜Â±Ã™â€¡Ã™â€  Ã™Ë† Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡ Ã™â€¦Ã˜Â§Ã™â€¡Ã˜Â§Ã™â€ Ã™â€¡ Ã˜Â«Ã˜Â¨Ã˜Âª Ã™â€¦Ã›Å’Ã¢â‚¬Å’Ã˜Â´Ã™Ë†Ã˜Â¯.',
    icon: Icons.key_outlined,
  ),
];

class _PriceFields extends StatelessWidget {
  const _PriceFields({
    required this.types,
    required this.propertyType,
    required this.salePrice,
    required this.finalPrice,
    required this.depositPrice,
    required this.rentPrice,
    required this.convertible,
    required this.maxConvertibleDeposit,
    required this.rentWithOwner,
  });

  final RxList<String> types;
  final RxString propertyType;
  final TextEditingController salePrice;
  final TextEditingController finalPrice;
  final TextEditingController depositPrice;
  final TextEditingController rentPrice;
  final RxBool convertible;
  final TextEditingController maxConvertibleDeposit;
  final RxBool rentWithOwner;

  @override
  Widget build(BuildContext context) {
    Widget pair(Widget first, Widget second) => LayoutBuilder(
      builder: (context, constraints) {
        final compact = constraints.maxWidth < 560;
        if (compact) {
          return Column(children: [first, const SizedBox(height: 12), second]);
        }
        return Row(
          children: [
            Expanded(child: first),
            const SizedBox(width: 12),
            Expanded(child: second),
          ],
        );
      },
    );
    return Obx(() {
      final children = <Widget>[];
      if (types.contains('sale') || types.contains('partnership')) {
        children.add(
          pair(
            _NumberField(
              controller: salePrice,
              label: types.contains('partnership') && !types.contains('sale')
                  ? 'Ø§Ø±Ø²Ø´ Ù…Ù„Ú©'
                  : 'Ù‚ÛŒÙ…Øª ÙØ±ÙˆØ´',
            ),
            _NumberField(controller: finalPrice, label: 'Ù‚ÛŒÙ…Øª Ù†Ù‡Ø§ÛŒÛŒ'),
          ),
        );
      }
      if (types.contains('rent_lease')) {
        if (children.isNotEmpty) children.add(const SizedBox(height: 12));
        children.addAll([
          pair(
            _NumberField(
              controller: depositPrice,
              label: 'Ã™â€¦Ã˜Â¨Ã™â€žÃ˜Âº Ã˜Â±Ã™â€¡Ã™â€ ',
            ),
            _NumberField(controller: rentPrice, label: 'Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡'),
          ),
          const SizedBox(height: 12),
          SwitchListTile(
            value: convertible.value,
            onChanged: (value) => convertible.value = value,
            contentPadding: EdgeInsets.zero,
            title: const Text('Ù‚Ø§Ø¨Ù„ ØªØ¨Ø¯ÛŒÙ„'),
            subtitle: const Text(
              'Ø¨Ø¹Ø¯Ø§ Ù‚ÛŒÙ…Øª Ø±Ù‡Ù†/Ø§Ø¬Ø§Ø±Ù‡ Ø¨Ø§ Ù†Ø³Ø¨Øª Ù‡Ø± Û±Û°Û° Ù…ÛŒÙ„ÛŒÙˆÙ†ØŒ Û³ Ù…ÛŒÙ„ÛŒÙˆÙ† Ø§Ø¬Ø§Ø±Ù‡ Ù…Ø­Ø§Ø³Ø¨Ù‡ Ù…ÛŒâ€ŒØ´ÙˆØ¯.',
            ),
          ),
          if (convertible.value)
            _NumberField(
              controller: maxConvertibleDeposit,
              label: 'Ø­Ø¯Ø§Ú©Ø«Ø± Ø±Ù‡Ù† Ù‚Ø§Ø¨Ù„ ØªØ¨Ø¯ÛŒÙ„',
            ),
          if (propertyType.value == 'apartment')
            SwitchListTile(
              value: rentWithOwner.value,
              onChanged: (value) => rentWithOwner.value = value,
              contentPadding: EdgeInsets.zero,
              title: const Text('Ù…Ù„Ú© Ù‡Ù…Ø±Ø§Ù‡ Ø¨Ø§ Ù…Ø§Ù„Ú© Ø§Ø³Øª'),
            ),
        ]);
      }
      return Column(children: children);
    });
  }
}

class _HouseInfoForm extends StatelessWidget {
  const _HouseInfoForm({
    required this.propertyType,
    required this.areaM2,
    required this.bedrooms,
    required this.floor,
    required this.totalFloors,
    required this.ageYears,
    required this.terraceCount,
    required this.backyardAreaM2,
    required this.gardenBuildingAreaM2,
    required this.gardenBuildingFloors,
    required this.documentType,
    required this.parking,
    required this.elevator,
    required this.storage,
    required this.renovated,
    required this.separateEntrance,
    required this.stoveTop,
    required this.parkingDisturbance,
    required this.parkingInDeed,
    required this.parkingCommon,
    required this.masterServiceCount,
    required this.doubleGlazedWindows,
    required this.directions,
    required this.cooling,
    required this.flooring,
    required this.heating,
    required this.cabinetType,
    required this.wallCovering,
    required this.kitchenType,
    required this.gardenBuilding,
    required this.pool,
    required this.waterUtility,
    required this.electricityUtility,
    required this.gasUtility,
    required this.waterRight,
    required this.permit,
  });

  final RxString propertyType;
  final TextEditingController areaM2;
  final TextEditingController bedrooms;
  final TextEditingController floor;
  final TextEditingController totalFloors;
  final TextEditingController ageYears;
  final TextEditingController terraceCount;
  final TextEditingController backyardAreaM2;
  final TextEditingController gardenBuildingAreaM2;
  final TextEditingController gardenBuildingFloors;
  final RxString documentType;
  final RxBool parking;
  final RxBool elevator;
  final RxBool storage;
  final RxBool renovated;
  final RxBool separateEntrance;
  final RxBool stoveTop;
  final RxBool parkingDisturbance;
  final RxBool parkingInDeed;
  final RxBool parkingCommon;
  final TextEditingController masterServiceCount;
  final RxBool doubleGlazedWindows;
  final RxList<String> directions;
  final RxList<String> cooling;
  final RxString flooring;
  final RxString heating;
  final RxString cabinetType;
  final RxString wallCovering;
  final RxString kitchenType;
  final RxBool gardenBuilding;
  final RxBool pool;
  final RxBool waterUtility;
  final RxBool electricityUtility;
  final RxBool gasUtility;
  final RxBool waterRight;
  final RxBool permit;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Ã˜Â§Ã˜Â·Ã™â€žÃ˜Â§Ã˜Â¹Ã˜Â§Ã˜Âª Ã˜Â®Ã˜Â§Ã™â€ Ã™â€¡',
          style: TextStyle(fontWeight: FontWeight.w900),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 12,
          runSpacing: 12,
          children: [
            SizedBox(
              width: 180,
              child: Obx(
                () => DropdownButtonFormField<String>(
                  initialValue: propertyType.value,
                  decoration: const InputDecoration(
                    labelText: 'Ã™â€ Ã™Ë†Ã˜Â¹ Ã™â€¦Ã™â€žÃšÂ©',
                  ),
                  items: _propertyTypeOptions
                      .map(
                        (option) => DropdownMenuItem(
                          value: option.value,
                          child: Text(option.label),
                        ),
                      )
                      .toList(),
                  onChanged: (value) =>
                      propertyType.value = value ?? 'apartment',
                ),
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: areaM2,
                label: 'Ã™â€¦Ã˜ÂªÃ˜Â±Ã˜Â§ÃšËœ',
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: bedrooms,
                label: 'Ã˜Â®Ã™Ë†Ã˜Â§Ã˜Â¨',
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: floor,
                label: 'Ã˜Â·Ã˜Â¨Ã™â€šÃ™â€¡',
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: totalFloors,
                label: 'ÃšÂ©Ã™â€ž Ã˜Â·Ã˜Â¨Ã™â€šÃ˜Â§Ã˜Âª',
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: ageYears,
                label: 'Ã˜Â³Ã™â€  Ã˜Â¨Ã™â€ Ã˜Â§',
              ),
            ),
            SizedBox(
              width: 140,
              child: _NumberField(
                controller: terraceCount,
                label: 'Ã˜ÂªÃ˜Â¹Ã˜Â¯Ã˜Â§Ã˜Â¯ Ã˜ÂªÃ˜Â±Ã˜Â§Ã˜Â³',
              ),
            ),
            Obx(
              () => SizedBox(
                width: 160,
                child: _NumberField(
                  controller: backyardAreaM2,
                  label:
                      'Ã™â€¦Ã˜ÂªÃ˜Â±Ã˜Â§ÃšËœ Ã˜Â­Ã›Å’Ã˜Â§Ã˜Â· Ã˜Â®Ã™â€žÃ™Ë†Ã˜Âª',
                ),
              ),
            ),
            SizedBox(
              width: 180,
              child: Obx(
                () => DropdownButtonFormField<String>(
                  initialValue: documentType.value,
                  decoration: const InputDecoration(
                    labelText: 'Ã™â€ Ã™Ë†Ã˜Â¹ Ã˜Â³Ã™â€ Ã˜Â¯',
                  ),
                  items: _documentTypeOptions
                      .map(
                        (option) => DropdownMenuItem(
                          value: option.value,
                          child: Text(option.label),
                        ),
                      )
                      .toList(),
                  onChanged: (value) =>
                      documentType.value = value ?? 'six_dang',
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Obx(
          () => Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              FilterChip(
                label: const Text('Ã™Â¾Ã˜Â§Ã˜Â±ÃšÂ©Ã›Å’Ã™â€ ÃšÂ¯'),
                selected: parking.value,
                onSelected: (v) => parking.value = v,
              ),
              FilterChip(
                label: const Text('Ã˜Â¢Ã˜Â³Ã˜Â§Ã™â€ Ã˜Â³Ã™Ë†Ã˜Â±'),
                selected: elevator.value,
                onSelected: (v) => elevator.value = v,
              ),
              FilterChip(
                label: const Text('Ã˜Â§Ã™â€ Ã˜Â¨Ã˜Â§Ã˜Â±Ã›Å’'),
                selected: storage.value,
                onSelected: (v) => storage.value = v,
              ),
              FilterChip(
                label: const Text('Ã˜Â¨Ã˜Â§Ã˜Â²Ã˜Â³Ã˜Â§Ã˜Â²Ã›Å’ Ã˜Â´Ã˜Â¯Ã™â€¡'),
                selected: renovated.value,
                onSelected: (v) => renovated.value = v,
              ),
              if (propertyType.value != 'villa')
                FilterChip(
                  label: const Text('Ã™Ë†Ã˜Â±Ã™Ë†Ã˜Â¯Ã›Å’ Ã˜Â¬Ã˜Â¯Ã˜Â§'),
                  selected: separateEntrance.value,
                  onSelected: (v) => separateEntrance.value = v,
                ),
            ],
          ),
        ),
        Obx(() {
          if (propertyType.value != 'garden' && propertyType.value != 'land') {
            return const SizedBox.shrink();
          }
          final isGarden = propertyType.value == 'garden';
          return Padding(
            padding: const EdgeInsets.only(top: 14),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isGarden
                      ? 'Ù…Ø´Ø®ØµØ§Øª Ø¨Ø§Øº'
                      : 'Ù…Ø´Ø®ØµØ§Øª Ø²Ù…ÛŒÙ† Ø®Ø§Ù„ÛŒ',
                  style: Theme.of(
                    context,
                  ).textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w900),
                ),
                const SizedBox(height: 10),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: [
                    FilterChip(
                      label: const Text('Ø§Ù…ØªÛŒØ§Ø² Ø¢Ø¨'),
                      selected: waterUtility.value,
                      onSelected: (v) => waterUtility.value = v,
                    ),
                    FilterChip(
                      label: const Text('Ø§Ù…ØªÛŒØ§Ø² Ø¨Ø±Ù‚'),
                      selected: electricityUtility.value,
                      onSelected: (v) => electricityUtility.value = v,
                    ),
                    FilterChip(
                      label: const Text('Ø§Ù…ØªÛŒØ§Ø² Ú¯Ø§Ø²'),
                      selected: gasUtility.value,
                      onSelected: (v) => gasUtility.value = v,
                    ),
                    FilterChip(
                      label: const Text('Ù…Ø¬ÙˆØ² Ø¯Ø§Ø±Ø¯'),
                      selected: permit.value,
                      onSelected: (v) => permit.value = v,
                    ),
                    if (isGarden)
                      FilterChip(
                        label: const Text('Ø§Ù…ØªÛŒØ§Ø² Ø¢Ø¨ Ù…Ø¯Ø§Ø±'),
                        selected: waterRight.value,
                        onSelected: (v) => waterRight.value = v,
                      ),
                    if (isGarden)
                      FilterChip(
                        label: const Text('Ø§Ø³ØªØ®Ø±'),
                        selected: pool.value,
                        onSelected: (v) => pool.value = v,
                      ),
                    if (isGarden)
                      FilterChip(
                        label: const Text('Ø¯Ø§Ø±Ø§ÛŒ Ø³Ø§Ø®ØªÙ…Ø§Ù†'),
                        selected: gardenBuilding.value,
                        onSelected: (v) => gardenBuilding.value = v,
                      ),
                  ],
                ),
                if (isGarden && gardenBuilding.value) ...[
                  const SizedBox(height: 12),
                  Wrap(
                    spacing: 12,
                    runSpacing: 12,
                    children: [
                      SizedBox(
                        width: 180,
                        child: _NumberField(
                          controller: gardenBuildingAreaM2,
                          label: 'Ù…ØªØ±Ø§Ú˜ Ø³Ø§Ø®ØªÙ…Ø§Ù†',
                        ),
                      ),
                      SizedBox(
                        width: 180,
                        child: _NumberField(
                          controller: gardenBuildingFloors,
                          label: 'Ø·Ø¨Ù‚Ø§Øª Ø³Ø§Ø®ØªÙ…Ø§Ù†',
                        ),
                      ),
                    ],
                  ),
                ],
              ],
            ),
          );
        }),
        const SizedBox(height: 16),
        Text(
          'Ã™Ë†Ã›Å’ÃšËœÃšÂ¯Ã›Å’Ã¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜ÂªÃšÂ©Ã™â€¦Ã›Å’Ã™â€žÃ›Å’',
          style: Theme.of(
            context,
          ).textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w900),
        ),
        const SizedBox(height: 12),
        Obx(
          () => Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _directionOptions
                .map(
                  (item) => FilterChip(
                    label: Text(item.label),
                    selected: directions.contains(item.value),
                    onSelected: (selected) {
                      if (selected) {
                        if (directions.length >= 2) {
                          Get.snackbar(
                            'Ã˜Â¬Ã™â€¡Ã˜Âª Ã™Ë†Ã˜Â§Ã˜Â­Ã˜Â¯',
                            'Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã™Ë†Ã˜Â§Ã˜Â­Ã˜Â¯Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â³Ã˜Â± Ã™â€ Ã˜Â¨Ã˜Â´ Ã˜Â­Ã˜Â¯Ã˜Â§ÃšÂ©Ã˜Â«Ã˜Â± Ã˜Â¯Ã™Ë† Ã˜Â¬Ã™â€¡Ã˜Âª Ã˜Â§Ã™â€ Ã˜ÂªÃ˜Â®Ã˜Â§Ã˜Â¨ ÃšÂ©Ã™â€ Ã›Å’Ã˜Â¯.',
                          );
                          return;
                        }
                        directions.add(item.value);
                      } else {
                        directions.remove(item.value);
                      }
                    },
                  ),
                )
                .toList(),
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 12,
          runSpacing: 12,
          children: [
            _OptionDropdown(
              width: 180,
              label: 'ÃšÂ©Ã™Â',
              value: flooring,
              options: _flooringOptions,
            ),
            _OptionDropdown(
              width: 180,
              label: 'ÃšÂ¯Ã˜Â±Ã™â€¦Ã˜Â§Ã›Å’Ã˜Â´',
              value: heating,
              options: _heatingOptions,
            ),
            _OptionDropdown(
              width: 180,
              label: 'Ã™â€ Ã™Ë†Ã˜Â¹ ÃšÂ©Ã˜Â§Ã˜Â¨Ã›Å’Ã™â€ Ã˜Âª',
              value: cabinetType,
              options: _cabinetOptions,
            ),
            _OptionDropdown(
              width: 180,
              label: 'Ã™Â¾Ã™Ë†Ã˜Â´Ã˜Â´ Ã˜Â¯Ã›Å’Ã™Ë†Ã˜Â§Ã˜Â±',
              value: wallCovering,
              options: _wallCoveringOptions,
            ),
            _OptionDropdown(
              width: 180,
              label: 'Ã™â€ Ã™Ë†Ã˜Â¹ Ã˜Â¢Ã˜Â´Ã™Â¾Ã˜Â²Ã˜Â®Ã˜Â§Ã™â€ Ã™â€¡',
              value: kitchenType,
              options: _kitchenOptions,
            ),
            SizedBox(
              width: 180,
              child: _NumberField(
                controller: masterServiceCount,
                label: 'Ã˜ÂªÃ˜Â¹Ã˜Â¯Ã˜Â§Ã˜Â¯ Ã™â€¦Ã˜Â³Ã˜ÂªÃ˜Â±',
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Obx(
          () => Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              ..._coolingOptions.map(
                (item) => FilterChip(
                  label: Text(item.label),
                  selected: cooling.contains(item.value),
                  onSelected: (selected) {
                    if (selected) {
                      cooling.add(item.value);
                    } else {
                      cooling.remove(item.value);
                    }
                  },
                ),
              ),
              FilterChip(
                label: const Text('ÃšÂ¯Ã˜Â§Ã˜Â² Ã˜Â±Ã™Ë†Ã™â€¦Ã›Å’Ã˜Â²Ã›Å’'),
                selected: stoveTop.value,
                onSelected: (v) => stoveTop.value = v,
              ),
              FilterChip(
                label: const Text(
                  'Ã™Â¾Ã˜Â§Ã˜Â±ÃšÂ©Ã›Å’Ã™â€ ÃšÂ¯ Ã˜Â¨Ã˜Â§ Ã™â€¦Ã˜Â²Ã˜Â§Ã˜Â­Ã™â€¦Ã˜Âª',
                ),
                selected: parkingDisturbance.value,
                onSelected: (v) => parkingDisturbance.value = v,
              ),
              FilterChip(
                label: const Text(
                  'Ã™Â¾Ã˜Â§Ã˜Â±ÃšÂ©Ã›Å’Ã™â€ ÃšÂ¯ Ã˜Â¯Ã˜Â± Ã˜Â³Ã™â€ Ã˜Â¯',
                ),
                selected: parkingInDeed.value,
                onSelected: (v) => parkingInDeed.value = v,
              ),
              FilterChip(
                label: const Text(
                  'Ã™Â¾Ã˜Â§Ã˜Â±ÃšÂ©Ã›Å’Ã™â€ ÃšÂ¯ Ã™â€¦Ã˜Â´Ã˜Â§Ã˜Â¹',
                ),
                selected: parkingCommon.value,
                onSelected: (v) => parkingCommon.value = v,
              ),
              FilterChip(
                label: const Text(
                  'Ã™Â¾Ã™â€ Ã˜Â¬Ã˜Â±Ã™â€¡ Ã˜Â¯Ã™Ë†Ã˜Â¬Ã˜Â¯Ã˜Â§Ã˜Â±Ã™â€¡',
                ),
                selected: doubleGlazedWindows.value,
                onSelected: (v) => doubleGlazedWindows.value = v,
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class _AddressDraft {
  final areaId = ''.obs;
  final streetId = ''.obs;
  final neighborhoodId = ''.obs;
  final manual = TextEditingController();

  void dispose() => manual.dispose();
}

class _AddressEditor extends StatelessWidget {
  const _AddressEditor({
    required this.index,
    required this.draft,
    required this.canRemove,
    required this.onRemove,
  });

  final int index;
  final _AddressDraft draft;
  final bool canRemove;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final locations = Get.find<LocationsController>();
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: DecoratedBox(
        decoration: BoxDecoration(
          border: Border.all(color: Theme.of(context).dividerColor),
          borderRadius: BorderRadius.circular(16),
        ),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Obx(() {
            final areas = locations.areas;
            final area = areas.firstWhereOrNull(
              (item) => item.id == draft.areaId.value,
            );
            final streets = area?.streets ?? const <StreetNode>[];
            final street = streets.firstWhereOrNull(
              (item) => item.id == draft.streetId.value,
            );
            final neighborhoods =
                street?.neighborhoods ?? const <NeighborhoodNode>[];
            return Column(
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        'Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ ${index + 1}',
                        style: const TextStyle(fontWeight: FontWeight.w800),
                      ),
                    ),
                    if (canRemove)
                      IconButton(
                        tooltip: 'Ã˜Â­Ã˜Â°Ã™Â Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³',
                        onPressed: onRemove,
                        icon: const Icon(Icons.close_rounded),
                      ),
                  ],
                ),
                LayoutBuilder(
                  builder: (context, constraints) {
                    final compact = constraints.maxWidth < 760;
                    final width = compact
                        ? constraints.maxWidth
                        : (constraints.maxWidth - 24) / 3;
                    return Wrap(
                      spacing: 12,
                      runSpacing: 12,
                      children: [
                        SizedBox(
                          width: width,
                          child: DropdownButtonFormField<String>(
                            initialValue: draft.areaId.value.isEmpty
                                ? null
                                : draft.areaId.value,
                            decoration: const InputDecoration(
                              labelText: 'Ã™â€¦Ã™â€ Ã˜Â·Ã™â€šÃ™â€¡',
                            ),
                            items: areas
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e.id,
                                    child: Text(e.name),
                                  ),
                                )
                                .toList(),
                            onChanged: (value) {
                              draft.areaId.value = value ?? '';
                              draft.streetId.value = '';
                              draft.neighborhoodId.value = '';
                            },
                          ),
                        ),
                        SizedBox(
                          width: width,
                          child: DropdownButtonFormField<String>(
                            initialValue: draft.streetId.value.isEmpty
                                ? null
                                : draft.streetId.value,
                            decoration: const InputDecoration(
                              labelText: 'Ã˜Â®Ã›Å’Ã˜Â§Ã˜Â¨Ã˜Â§Ã™â€ ',
                            ),
                            items: streets
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e.id,
                                    child: Text(e.name),
                                  ),
                                )
                                .toList(),
                            onChanged: (value) {
                              draft.streetId.value = value ?? '';
                              draft.neighborhoodId.value = '';
                            },
                          ),
                        ),
                        SizedBox(
                          width: width,
                          child: DropdownButtonFormField<String>(
                            initialValue: draft.neighborhoodId.value.isEmpty
                                ? null
                                : draft.neighborhoodId.value,
                            decoration: const InputDecoration(
                              labelText: 'Ã™â€¦Ã˜Â­Ã™â€žÃ™â€¡',
                            ),
                            items: neighborhoods
                                .map(
                                  (e) => DropdownMenuItem(
                                    value: e.id,
                                    child: Text(e.name),
                                  ),
                                )
                                .toList(),
                            onChanged: (value) =>
                                draft.neighborhoodId.value = value ?? '',
                          ),
                        ),
                      ],
                    );
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: draft.manual,
                  decoration: const InputDecoration(
                    labelText:
                        'Ã˜Â¢Ã˜Â¯Ã˜Â±Ã˜Â³ Ã˜Â¯Ã™â€šÃ›Å’Ã™â€š Ã˜Â¯Ã˜Â³Ã˜ÂªÃ›Å’',
                  ),
                  minLines: 1,
                  maxLines: 3,
                ),
              ],
            );
          }),
        ),
      ),
    );
  }
}

class _OptionDropdown extends StatelessWidget {
  const _OptionDropdown({
    required this.width,
    required this.label,
    required this.value,
    required this.options,
  });

  final double width;
  final String label;
  final RxString value;
  final List<_OptionItem> options;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width,
      child: Obx(
        () => DropdownButtonFormField<String>(
          initialValue: value.value,
          decoration: InputDecoration(labelText: label),
          items: options
              .map(
                (option) => DropdownMenuItem(
                  value: option.value,
                  child: Text(option.label),
                ),
              )
              .toList(),
          onChanged: (selected) =>
              value.value = selected ?? options.first.value,
        ),
      ),
    );
  }
}

class _NumberField extends StatelessWidget {
  const _NumberField({required this.controller, required this.label});

  final TextEditingController controller;
  final String label;

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      keyboardType: const TextInputType.numberWithOptions(
        signed: false,
        decimal: false,
      ),
      inputFormatters: [FilteringTextInputFormatter.digitsOnly],
      decoration: InputDecoration(labelText: label),
      onChanged: (value) {
        final cleaned = value.replaceAll(RegExp(r'[^0-9]'), '');
        if (cleaned != value) {
          controller.value = TextEditingValue(
            text: cleaned,
            selection: TextSelection.collapsed(offset: cleaned.length),
          );
        }
      },
    );
  }
}

String _typeLabel(String type) => switch (type) {
  'sale' => 'Ã™ÂÃ˜Â±Ã™Ë†Ã˜Â´',
  'partnership' => 'Ã™â€¦Ã˜Â´Ã˜Â§Ã˜Â±ÃšÂ©Ã˜Âª',
  'rent_lease' => 'Ã˜Â±Ã™â€¡Ã™â€  Ã™Ë† Ã˜Â§Ã˜Â¬Ã˜Â§Ã˜Â±Ã™â€¡',
  _ => type,
};

IconData _typeIcon(String type) => switch (type) {
  'sale' => Icons.sell_outlined,
  'partnership' => Icons.handshake_outlined,
  'rent_lease' => Icons.key_outlined,
  _ => Icons.apartment_outlined,
};

IconData _uploadIcon(String name) {
  final lower = name.toLowerCase();
  if (lower.endsWith('.mp4') ||
      lower.endsWith('.mov') ||
      lower.endsWith('.webm')) {
    return Icons.movie_outlined;
  }
  return Icons.image_outlined;
}

class _PropertyTypeOption {
  const _PropertyTypeOption(this.value, this.label);

  final String value;
  final String label;
}

class _OptionItem {
  const _OptionItem(this.value, this.label);

  final String value;
  final String label;
}

const _propertyTypeOptions = [
  _PropertyTypeOption('apartment', 'Ã˜Â¢Ã™Â¾Ã˜Â§Ã˜Â±Ã˜ÂªÃ™â€¦Ã˜Â§Ã™â€ '),
  _PropertyTypeOption('villa', 'Ã™Ë†Ã›Å’Ã™â€žÃ˜Â§'),
  _PropertyTypeOption('land', 'Ã˜Â²Ã™â€¦Ã›Å’Ã™â€ '),
  _PropertyTypeOption('office', 'Ã˜Â§Ã˜Â¯Ã˜Â§Ã˜Â±Ã›Å’'),
  _PropertyTypeOption('shop', 'Ã˜ÂªÃ˜Â¬Ã˜Â§Ã˜Â±Ã›Å’ / Ã™â€¦Ã˜ÂºÃ˜Â§Ã˜Â²Ã™â€¡'),
  _PropertyTypeOption('garden', 'Ã˜Â¨Ã˜Â§Ã˜Âº'),
  _PropertyTypeOption('old_house', 'ÃšÂ©Ã™â€žÃ™â€ ÃšÂ¯Ã›Å’'),
];

const _directionOptions = [
  _OptionItem('north', 'Ã˜Â´Ã™â€¦Ã˜Â§Ã™â€žÃ›Å’'),
  _OptionItem('south', 'Ã˜Â¬Ã™â€ Ã™Ë†Ã˜Â¨Ã›Å’'),
  _OptionItem('east', 'Ã˜Â´Ã˜Â±Ã™â€šÃ›Å’'),
  _OptionItem('west', 'Ã˜ÂºÃ˜Â±Ã˜Â¨Ã›Å’'),
];

const _flooringOptions = [
  _OptionItem('ceramic', 'Ã˜Â³Ã˜Â±Ã˜Â§Ã™â€¦Ã›Å’ÃšÂ©'),
  _OptionItem('mosaic', 'Ã™â€¦Ã™Ë†Ã˜Â²Ã˜Â§Ã›Å’Ã›Å’ÃšÂ©'),
  _OptionItem('carpet', 'Ã™â€¦Ã™Ë†ÃšÂ©Ã˜Âª'),
  _OptionItem('parquet', 'Ã™Â¾Ã˜Â§Ã˜Â±ÃšÂ©Ã˜Âª'),
  _OptionItem('stone', 'Ã˜Â³Ã™â€ ÃšÂ¯'),
  _OptionItem('laminate', 'Ã™â€žÃ™â€¦Ã›Å’Ã™â€ Ã˜Âª'),
];

const _heatingOptions = [
  _OptionItem('radiator', 'Ã˜Â´Ã™Ë†Ã™ÂÃ˜Â§ÃšËœ'),
  _OptionItem('floor_heating', 'ÃšÂ¯Ã˜Â±Ã™â€¦Ã˜Â§Ã›Å’Ã˜Â´ Ã˜Â§Ã˜Â² ÃšÂ©Ã™Â'),
  _OptionItem('package', 'Ã™Â¾ÃšÂ©Ã›Å’Ã˜Â¬'),
  _OptionItem('heater', 'Ã˜Â¨Ã˜Â®Ã˜Â§Ã˜Â±Ã›Å’'),
  _OptionItem('fan_coil', 'Ã™ÂÃ™â€ Ã¢â‚¬Å’ÃšÂ©Ã™Ë†Ã›Å’Ã™â€ž'),
];

const _cabinetOptions = [
  _OptionItem('mdf', 'MDF'),
  _OptionItem('metal', 'Ã™ÂÃ™â€žÃ˜Â²Ã›Å’'),
  _OptionItem('high_gloss', 'Ã™â€¡Ã˜Â§Ã›Å’Ã¢â‚¬Å’ÃšÂ¯Ã™â€žÃ˜Â³'),
  _OptionItem('wood', 'Ãšâ€ Ã™Ë†Ã˜Â¨Ã›Å’'),
  _OptionItem('membrane', 'Ã™â€¦Ã™â€¦Ã˜Â¨Ã˜Â±Ã˜Â§Ã™â€ '),
];

const _coolingOptions = [
  _OptionItem('water_cooler', 'ÃšÂ©Ã™Ë†Ã™â€žÃ˜Â± Ã˜Â¢Ã˜Â¨Ã›Å’'),
  _OptionItem('split', 'ÃšÂ©Ã™Ë†Ã™â€žÃ˜Â± ÃšÂ¯Ã˜Â§Ã˜Â²Ã›Å’'),
  _OptionItem('chiller', 'Ãšâ€ Ã›Å’Ã™â€žÃ˜Â±'),
  _OptionItem('fan_coil', 'Ã™ÂÃ™â€ Ã¢â‚¬Å’ÃšÂ©Ã™Ë†Ã›Å’Ã™â€ž'),
];

const _wallCoveringOptions = [
  _OptionItem('paint', 'Ã˜Â±Ã™â€ ÃšÂ¯'),
  _OptionItem('wallpaper', 'ÃšÂ©Ã˜Â§Ã˜ÂºÃ˜Â° Ã˜Â¯Ã›Å’Ã™Ë†Ã˜Â§Ã˜Â±Ã›Å’'),
  _OptionItem('stone', 'Ã˜Â³Ã™â€ ÃšÂ¯'),
  _OptionItem('wood', 'Ãšâ€ Ã™Ë†Ã˜Â¨'),
  _OptionItem('decorative_panel', 'Ã˜Â¯Ã›Å’Ã™Ë†Ã˜Â§Ã˜Â±Ã™Â¾Ã™Ë†Ã˜Â´'),
];

const _kitchenOptions = [
  _OptionItem('open', 'Ã˜Â§Ã™Ë†Ã™Â¾Ã™â€ '),
  _OptionItem('island', 'Ã˜Â¬Ã˜Â²Ã›Å’Ã˜Â±Ã™â€¡'),
  _OptionItem('closed', 'Ã˜Â¨Ã˜Â¯Ã™Ë†Ã™â€  Ã˜Â§Ã™Ë†Ã™Â¾Ã™â€ '),
  _OptionItem('semi_open', 'Ã™â€ Ã›Å’Ã™â€¦Ã™â€¡ Ã˜Â§Ã™Ë†Ã™Â¾Ã™â€ '),
];

String _optionLabel(List<_OptionItem> options, String value) => options
    .firstWhere((option) => option.value == value, orElse: () => options.first)
    .label;

class _DocumentTypeOption {
  const _DocumentTypeOption(this.value, this.label);

  final String value;
  final String label;
}

const _documentTypeOptions = [
  _DocumentTypeOption('six_dang', 'Ã˜Â´Ã˜Â´Ã¢â‚¬Å’Ã˜Â¯Ã˜Â§Ã™â€ ÃšÂ¯'),
  _DocumentTypeOption('vekalaati', 'Ã™Ë†ÃšÂ©Ã˜Â§Ã™â€žÃ˜ÂªÃ›Å’'),
  _DocumentTypeOption('astaneh', 'Ã˜Â¢Ã˜Â³Ã˜ÂªÃ˜Â§Ã™â€ Ã™â€¡'),
  _DocumentTypeOption(
    'gholnameh',
    'Ã™â€šÃ™Ë†Ã™â€žÃ™â€ Ã˜Â§Ã™â€¦Ã™â€¡Ã¢â‚¬Å’Ã˜Â§Ã›Å’',
  ),
];

String _documentTypeLabel(String value) => _documentTypeOptions
    .firstWhere(
      (option) => option.value == value,
      orElse: () => _documentTypeOptions.first,
    )
    .label;

int _number(String value) {
  final normalized = RegExp(
    r'\d+',
  ).allMatches(value).map((match) => match.group(0)!).join();
  return int.tryParse(normalized) ?? 0;
}
