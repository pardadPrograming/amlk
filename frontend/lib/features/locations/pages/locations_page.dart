import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../data/models.dart';
import '../../../shared/responsive.dart';
import '../../auth/auth_controller.dart';
import '../locations_controller.dart';

class LocationsPage extends StatefulWidget {
  const LocationsPage({super.key});

  @override
  State<LocationsPage> createState() => _LocationsPageState();
}

class _LocationsPageState extends State<LocationsPage> {
  final controller = Get.find<LocationsController>();

  @override
  void initState() {
    super.initState();
    controller.load();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('مناطق و محله‌ها'),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _nameDialog(
          context,
          title: 'منطقه جدید',
          label: 'نام منطقه',
          onSubmit: controller.addArea,
        ),
        icon: const Icon(Icons.add_location_alt_outlined),
        label: const Text('افزودن منطقه'),
      ),
      body: ResponsivePage(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const GradientHeader(
              title: 'مناطق، خیابان‌ها و محله‌ها',
              subtitle:
                  'ساختار جغرافیایی فایل‌های ملکی را اینجا تعریف کنید. حذف منطقه یا خیابان، زیرشاخه‌های آن را هم حذف می‌کند.',
              icon: Icons.map_outlined,
            ),
            const SizedBox(height: 18),
            AppCard(
              padding: 16,
              child: Row(
                children: [
                  const Expanded(
                    child: Text(
                      'اگر شهر، منطقه، خیابان یا محله در لیست سیستمی نبود، آن را برای تایید مدیریت کل پیشنهاد دهید.',
                    ),
                  ),
                  const SizedBox(width: 12),
                  FilledButton.icon(
                    onPressed: () => _suggestLocationDialog(context),
                    icon: const Icon(Icons.add_location_alt_outlined),
                    label: const Text('پیشنهاد لوکیشن'),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 18),
            Obx(() {
              if (controller.loading.value && controller.areas.isEmpty) {
                return const Center(
                  child: Padding(
                    padding: EdgeInsets.all(36),
                    child: CircularProgressIndicator(),
                  ),
                );
              }
              if (controller.areas.isEmpty) {
                return const AppCard(child: _EmptyLocations());
              }
              return Column(
                children: controller.areas
                    .map((area) => _AreaCard(area: area))
                    .toList(),
              );
            }),
          ],
        ),
      ),
    );
  }
}

Future<void> _suggestLocationDialog(BuildContext context) async {
  final controller = Get.find<LocationsController>();
  await controller.loadSystemCities();
  if (controller.systemCities.isEmpty) {
    Get.snackbar(
      'شهر سیستمی وجود ندارد',
      'ابتدا مدیر کل باید حداقل یک شهر سیستمی ایجاد کند.',
    );
    return;
  }
  final name = TextEditingController();
  final parent = TextEditingController();
  final formKey = GlobalKey<FormState>();
  final profileCityId = Get.find<AuthController>().user.value?.cityId ?? '';
  final visibleCities = profileCityId.isEmpty
      ? controller.systemCities
      : controller.systemCities
            .where((city) => city.id == profileCityId)
            .toList(growable: false);
  if (visibleCities.isEmpty) {
    Get.snackbar(
      'شهر پروفایل پیدا نشد',
      'از بخش پروفایل شهر فعال را انتخاب کنید.',
    );
    return;
  }
  var cityId = visibleCities.first.id;
  var type = 'neighborhood';
  await Get.dialog<void>(
    StatefulBuilder(
      builder: (context, setState) => AlertDialog(
        title: const Text('پیشنهاد لوکیشن سیستمی'),
        content: Form(
          key: formKey,
          child: SizedBox(
            width: 420,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                DropdownButtonFormField<String>(
                  initialValue: cityId,
                  decoration: const InputDecoration(labelText: 'شهر'),
                  items: visibleCities
                      .map(
                        (city) => DropdownMenuItem(
                          value: city.id,
                          child: Text(city.name),
                        ),
                      )
                      .toList(),
                  onChanged: (value) {
                    if (value != null) setState(() => cityId = value);
                  },
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: type,
                  decoration: const InputDecoration(labelText: 'نوع'),
                  items: const [
                    DropdownMenuItem(value: 'area', child: Text('منطقه')),
                    DropdownMenuItem(value: 'street', child: Text('خیابان')),
                    DropdownMenuItem(
                      value: 'neighborhood',
                      child: Text('محله'),
                    ),
                  ],
                  onChanged: (value) {
                    if (value != null) setState(() => type = value);
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: name,
                  decoration: const InputDecoration(labelText: 'نام پیشنهادی'),
                  validator: (value) => value == null || value.trim().isEmpty
                      ? 'نام پیشنهادی الزامی است'
                      : null,
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: parent,
                  decoration: const InputDecoration(
                    labelText: 'نام والد یا توضیح مسیر',
                    hintText: 'مثلا منطقه یا خیابان مرتبط',
                  ),
                ),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(onPressed: Get.back, child: const Text('انصراف')),
          FilledButton(
            onPressed: () async {
              if (!(formKey.currentState?.validate() ?? false)) return;
              Get.back<void>();
              await controller.suggestLocation(
                cityId: cityId,
                type: type,
                name: name.text,
                manualParentName: parent.text,
              );
            },
            child: const Text('ارسال برای تایید'),
          ),
        ],
      ),
    ),
  );
  name.dispose();
  parent.dispose();
}

class _AreaCard extends StatelessWidget {
  const _AreaCard({required this.area});

  final AreaNode area;

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<LocationsController>();
    return Padding(
      padding: const EdgeInsets.only(bottom: 14),
      child: AppCard(
        padding: 16,
        child: ExpansionTile(
          initiallyExpanded: true,
          leading: const Icon(Icons.location_city_outlined),
          title: Text(
            area.name,
            style: const TextStyle(fontWeight: FontWeight.w900),
          ),
          subtitle: Text('${area.streets.length} خیابان'),
          trailing: Wrap(
            spacing: 4,
            children: [
              IconButton(
                tooltip: 'افزودن خیابان',
                onPressed: () => _nameDialog(
                  context,
                  title: 'خیابان جدید',
                  label: 'نام خیابان',
                  onSubmit: (name) => controller.addStreet(area.id, name),
                ),
                icon: const Icon(Icons.add_road_outlined),
              ),
              IconButton(
                tooltip: 'حذف منطقه و زیرشاخه‌ها',
                onPressed: () => _confirmDelete(
                  context,
                  title: 'حذف منطقه',
                  message:
                      'با حذف این منطقه، تمام خیابان‌ها و محله‌های زیرمجموعه هم حذف می‌شوند.',
                  onConfirm: () => controller.deleteArea(area.id),
                ),
                icon: const Icon(Icons.delete_outline),
              ),
            ],
          ),
          children: area.streets.isEmpty
              ? const [
                  Padding(
                    padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
                    child: _MutedText(
                      'هنوز خیابانی برای این منطقه ثبت نشده است.',
                    ),
                  ),
                ]
              : area.streets
                    .map(
                      (street) => _StreetTile(areaId: area.id, street: street),
                    )
                    .toList(),
        ),
      ),
    );
  }
}

class _StreetTile extends StatelessWidget {
  const _StreetTile({required this.areaId, required this.street});

  final String areaId;
  final StreetNode street;

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<LocationsController>();
    return Padding(
      padding: const EdgeInsets.fromLTRB(8, 0, 8, 8),
      child: DecoratedBox(
        decoration: BoxDecoration(
          border: Border.all(color: Theme.of(context).dividerColor),
          borderRadius: BorderRadius.circular(16),
        ),
        child: ExpansionTile(
          leading: const Icon(Icons.alt_route_outlined),
          title: Text(
            street.name,
            style: const TextStyle(fontWeight: FontWeight.w800),
          ),
          subtitle: Text('${street.neighborhoods.length} محله'),
          trailing: Wrap(
            spacing: 4,
            children: [
              IconButton(
                tooltip: 'افزودن محله',
                onPressed: () => _nameDialog(
                  context,
                  title: 'محله جدید',
                  label: 'نام محله',
                  onSubmit: (name) =>
                      controller.addNeighborhood(areaId, street.id, name),
                ),
                icon: const Icon(Icons.add_home_work_outlined),
              ),
              IconButton(
                tooltip: 'حذف خیابان و محله‌ها',
                onPressed: () => _confirmDelete(
                  context,
                  title: 'حذف خیابان',
                  message:
                      'با حذف این خیابان، همه محله‌های زیرمجموعه حذف می‌شوند.',
                  onConfirm: () => controller.deleteStreet(areaId, street.id),
                ),
                icon: const Icon(Icons.delete_outline),
              ),
            ],
          ),
          children: street.neighborhoods.isEmpty
              ? const [
                  Padding(
                    padding: EdgeInsets.fromLTRB(16, 0, 16, 16),
                    child: _MutedText(
                      'هنوز محله‌ای برای این خیابان ثبت نشده است.',
                    ),
                  ),
                ]
              : street.neighborhoods
                    .map(
                      (neighborhood) => ListTile(
                        leading: const Icon(Icons.home_work_outlined),
                        title: Text(neighborhood.name),
                        trailing: IconButton(
                          tooltip: 'حذف محله',
                          onPressed: () => _confirmDelete(
                            context,
                            title: 'حذف محله',
                            message: 'این محله حذف شود؟',
                            onConfirm: () => controller.deleteNeighborhood(
                              areaId,
                              street.id,
                              neighborhood.id,
                            ),
                          ),
                          icon: const Icon(Icons.close_rounded),
                        ),
                      ),
                    )
                    .toList(),
        ),
      ),
    );
  }
}

class _EmptyLocations extends StatelessWidget {
  const _EmptyLocations();

  @override
  Widget build(BuildContext context) {
    return const Column(
      children: [
        Icon(Icons.map_outlined, size: 46),
        SizedBox(height: 12),
        Text(
          'هنوز منطقه‌ای ثبت نشده است.',
          style: TextStyle(fontWeight: FontWeight.w800),
        ),
        SizedBox(height: 6),
        Text('از دکمه افزودن منطقه شروع کنید.'),
      ],
    );
  }
}

class _MutedText extends StatelessWidget {
  const _MutedText(this.text);

  final String text;

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: Theme.of(
        context,
      ).textTheme.bodyMedium?.copyWith(color: Theme.of(context).hintColor),
    );
  }
}

Future<void> _nameDialog(
  BuildContext context, {
  required String title,
  required String label,
  required Future<void> Function(String name) onSubmit,
}) async {
  final textController = TextEditingController();
  final formKey = GlobalKey<FormState>();
  await Get.dialog<void>(
    AlertDialog(
      title: Text(title),
      content: Form(
        key: formKey,
        child: TextFormField(
          controller: textController,
          autofocus: true,
          textInputAction: TextInputAction.done,
          decoration: InputDecoration(labelText: label),
          validator: (value) => (value == null || value.trim().isEmpty)
              ? '$label الزامی است'
              : null,
          onFieldSubmitted: (_) async {
            if (formKey.currentState?.validate() ?? false) {
              Get.back<void>();
              await onSubmit(textController.text.trim());
            }
          },
        ),
      ),
      actions: [
        TextButton(onPressed: Get.back, child: const Text('انصراف')),
        FilledButton(
          onPressed: () async {
            if (formKey.currentState?.validate() ?? false) {
              Get.back<void>();
              await onSubmit(textController.text.trim());
            }
          },
          child: const Text('ثبت'),
        ),
      ],
    ),
  );
  textController.dispose();
}

Future<void> _confirmDelete(
  BuildContext context, {
  required String title,
  required String message,
  required Future<void> Function() onConfirm,
}) async {
  await Get.dialog<void>(
    AlertDialog(
      title: Text(title),
      content: Text(message),
      actions: [
        TextButton(onPressed: Get.back, child: const Text('انصراف')),
        FilledButton.tonal(
          onPressed: () async {
            Get.back<void>();
            await onConfirm();
          },
          child: const Text('حذف'),
        ),
      ],
    ),
  );
}
