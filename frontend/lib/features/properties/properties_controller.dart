import 'package:dio/dio.dart' as dio;
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

import '../../core/api/api_client.dart';
import '../../core/media/media_picker.dart';
import '../../data/models.dart';
import '../business/business_controller.dart';
import '../locations/locations_controller.dart';

class PropertiesController extends GetxController {
  PropertiesController(this._api);

  final ApiClient _api;
  final loading = false.obs;
  final files = <PropertyFileModel>[].obs;
  final vaults = <ChannelVaultModel>[].obs;
  final ownerShareRequests = <PropertyShareRequestModel>[].obs;
  final requesterShareRequests = <PropertyShareRequestModel>[].obs;
  final incomingOffers = <PropertyOfferModel>[].obs;
  final outgoingOffers = <PropertyOfferModel>[].obs;
  final notifications = <NotificationModel>[].obs;
  final selectedUploads = <PickedMedia>[].obs;
  final latestFiles = <PropertyFileModel>[].obs;
  final latestFilter = ''.obs;
  final latestLoading = false.obs;
  final latestLoadingMore = false.obs;
  final latestHasMore = true.obs;
  final latestUnreadCount = 0.obs;

  final _box = GetStorage();

  String? get _businessId => Get.find<BusinessController>().selected.value?.id;

  String _latestReadKey(String businessId) =>
      'latest-files-read-at.$businessId.${latestFilter.value}';

  Future<void> load() async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      final res = await _api.dio.get('/businesses/$businessId/properties');
      final list = res.data['data'] as List? ?? const [];
      files.assignAll(
        list.map(
          (e) => PropertyFileModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadLatestFiles({
    bool reset = true,
    bool markRead = false,
  }) async {
    final businessId = _businessId;
    if (businessId == null) return;
    if (reset) {
      latestLoading.value = true;
      latestHasMore.value = true;
    } else {
      if (latestLoadingMore.value || !latestHasMore.value) return;
      latestLoadingMore.value = true;
    }
    try {
      final query = <String, dynamic>{
        'limit': 30,
        'offset': reset ? 0 : latestFiles.length,
        if (latestFilter.value.isNotEmpty) 'type': latestFilter.value,
      };
      final res = await _api.dio.get(
        '/businesses/$businessId/properties/latest',
        queryParameters: query,
      );
      final page = (res.data['data'] as List? ?? const [])
          .map((e) => PropertyFileModel.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      if (reset) {
        latestFiles.assignAll(page);
      } else {
        final existing = latestFiles.map((item) => item.id).toSet();
        latestFiles.addAll(page.where((item) => !existing.contains(item.id)));
      }
      final total =
          int.tryParse(res.headers.value('x-total-count') ?? '') ??
          latestFiles.length;
      latestHasMore.value = latestFiles.length < total;
      if (markRead) {
        markLatestFilesRead();
      } else {
        _updateLatestUnreadFromLoaded(businessId);
      }
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      latestLoading.value = false;
      latestLoadingMore.value = false;
    }
  }

  Future<void> setLatestFilter(String type) async {
    latestFilter.value = type;
    await loadLatestFiles(reset: true, markRead: true);
  }

  Future<void> checkLatestUnread() async {
    final businessId = _businessId;
    if (businessId == null) return;
    try {
      final res = await _api.dio.get(
        '/businesses/$businessId/properties/latest',
        queryParameters: {
          'limit': 30,
          'offset': 0,
          if (latestFilter.value.isNotEmpty) 'type': latestFilter.value,
        },
      );
      final page = (res.data['data'] as List? ?? const [])
          .map((e) => PropertyFileModel.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      final readAt = DateTime.tryParse(
        _box.read(_latestReadKey(businessId))?.toString() ?? '',
      );
      if (readAt == null) {
        latestUnreadCount.value = 0;
        return;
      }
      latestUnreadCount.value = page.where((item) {
        final createdAt = DateTime.tryParse(item.createdAt);
        return createdAt != null && createdAt.isAfter(readAt);
      }).length;
    } catch (_) {
      // Silent polling should not interrupt the user.
    }
  }

  void markLatestFilesRead() {
    final businessId = _businessId;
    if (businessId == null) return;
    final newest = latestFiles.isEmpty ? null : latestFiles.first.createdAt;
    _box.write(
      _latestReadKey(businessId),
      newest == null || newest.isEmpty
          ? DateTime.now().toUtc().toIso8601String()
          : newest,
    );
    latestUnreadCount.value = 0;
  }

  bool isLatestFileUnread(PropertyFileModel file) {
    final businessId = _businessId;
    if (businessId == null) return false;
    final readAt = DateTime.tryParse(
      _box.read(_latestReadKey(businessId))?.toString() ?? '',
    );
    final createdAt = DateTime.tryParse(file.createdAt);
    return readAt != null && createdAt != null && createdAt.isAfter(readAt);
  }

  void _updateLatestUnreadFromLoaded(String businessId) {
    final readAt = DateTime.tryParse(
      _box.read(_latestReadKey(businessId))?.toString() ?? '',
    );
    if (readAt == null) {
      latestUnreadCount.value = 0;
      return;
    }
    latestUnreadCount.value = latestFiles.where((item) {
      final createdAt = DateTime.tryParse(item.createdAt);
      return createdAt != null && createdAt.isAfter(readAt);
    }).length;
  }

  Future<void> loadShareRequests() async {
    final businessId = _businessId;
    if (businessId == null) return;
    try {
      final ownerRes = await _api.dio.get(
        '/businesses/$businessId/property-share-requests',
      );
      ownerShareRequests.assignAll(
        (ownerRes.data['data'] as List? ?? const []).map(
          (e) =>
              PropertyShareRequestModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
      final requesterRes = await _api.dio.get(
        '/businesses/$businessId/property-share-requests?scope=requester',
      );
      requesterShareRequests.assignAll(
        (requesterRes.data['data'] as List? ?? const []).map(
          (e) =>
              PropertyShareRequestModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> loadOffers() async {
    final businessId = _businessId;
    if (businessId == null) return;
    try {
      final incomingRes = await _api.dio.get(
        '/businesses/$businessId/property-offers?scope=incoming',
      );
      incomingOffers.assignAll(
        (incomingRes.data['data'] as List? ?? const []).map(
          (e) => PropertyOfferModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
      final outgoingRes = await _api.dio.get(
        '/businesses/$businessId/property-offers?scope=outgoing',
      );
      outgoingOffers.assignAll(
        (outgoingRes.data['data'] as List? ?? const []).map(
          (e) => PropertyOfferModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> sendOffer(
    PropertyOfferModel offer,
    double commissionPercent,
  ) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/property-offers/${offer.id}/send',
        data: {'commissionPercent': commissionPercent},
      );
      await loadOffers();
      Get.snackbar('پیشنهاد ارسال شد', 'پیشنهاد فایل برای گیرنده ارسال شد.');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> respondOffer(PropertyOfferModel offer, bool approve) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/property-offers/${offer.id}/${approve ? 'accept' : 'reject'}',
      );
      await loadOffers();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> finalizeOffer(PropertyOfferModel offer, bool approve) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/property-offers/${offer.id}/${approve ? 'final-approve' : 'final-reject'}',
      );
      await loadOffers();
      await load();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadNotifications() async {
    try {
      final res = await _api.dio.get('/notifications');
      notifications.assignAll(
        (res.data['data'] as List? ?? const []).map(
          (e) => NotificationModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> markNotificationRead(NotificationModel notification) async {
    try {
      await _api.dio.post('/notifications/${notification.id}/read');
      final index = notifications.indexWhere(
        (item) => item.id == notification.id,
      );
      if (index >= 0) {
        notifications[index] = NotificationModel(
          id: notification.id,
          type: notification.type,
          title: notification.title,
          body: notification.body,
          businessId: notification.businessId,
          propertyId: notification.propertyId,
          requestId: notification.requestId,
          readAt: DateTime.now().toUtc().toIso8601String(),
        );
      }
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> loadVaults() async {
    final businessId = _businessId;
    final combined = <String, ChannelVaultModel>{};
    try {
      final userRes = await _api.dio.get('/vaults');
      for (final item in userRes.data['data'] as List? ?? const []) {
        final vault = ChannelVaultModel.fromJson(
          Map<String, dynamic>.from(item),
        );
        combined[vault.id] = vault;
      }
      if (businessId != null) {
        final businessRes = await _api.dio.get(
          '/businesses/$businessId/vaults',
        );
        for (final item in businessRes.data['data'] as List? ?? const []) {
          final vault = ChannelVaultModel.fromJson(
            Map<String, dynamic>.from(item),
          );
          combined[vault.id] = vault;
        }
      }
      vaults.assignAll(combined.values);
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    }
  }

  Future<void> createUserVault(String title) async {
    final name = title.trim();
    if (name.isEmpty) return;
    try {
      await _api.dio.post('/vaults', data: {'title': name});
      await loadVaults();
      Get.snackbar('صندوقچه ساخته شد', 'صندوقچه شخصی جدید آماده است.');
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    }
  }

  Future<void> createBusinessVault(String title) async {
    final businessId = _businessId;
    final name = title.trim();
    if (businessId == null || name.isEmpty) return;
    try {
      await _api.dio.post(
        '/businesses/$businessId/vaults',
        data: {'title': name},
      );
      await loadVaults();
      Get.snackbar('صندوقچه ساخته شد', 'صندوقچه املاک جدید آماده است.');
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    }
  }

  Future<void> pickMedia() async {
    final result = await pickMediaFiles();
    if (result.isEmpty) {
      Get.snackbar(
        'انتخاب فایل',
        'در نسخه فعلی فقط انتخاب فایل روی وب فعال است.',
      );
      return;
    }
    final merged = [...selectedUploads, ...result];
    final videos = merged.where((f) => _isVideo(f)).length;
    if (merged.length > 20) {
      Get.snackbar('محدودیت مدیا', 'برای هر فایل حداکثر ۲۰ مدیا مجاز است.');
      return;
    }
    if (videos > 2) {
      Get.snackbar('محدودیت ویدئو', 'برای هر فایل حداکثر ۲ ویدئو مجاز است.');
      return;
    }
    selectedUploads.assignAll(merged);
  }

  void removeUpload(PickedMedia file) {
    selectedUploads.remove(file);
  }

  Future<void> create({
    required String type,
    required List<String> types,
    required String title,
    required String description,
    required String internalDescription,
    required int salePrice,
    required int finalPrice,
    required int depositPrice,
    required int rentPrice,
    required bool convertible,
    required int maxConvertibleDeposit,
    required bool rentWithOwner,
    required Map<String, dynamic> houseInfo,
    required List<PropertyAddressModel> addresses,
    String status = 'active',
    List<PropertyVaultPlacementModel> vaultPlacements = const [],
  }) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      final res = await _api.dio.post(
        '/businesses/$businessId/properties',
        data: {
          'type': type,
          'status': status,
          'types': types,
          'title': title,
          'description': description,
          'internalDescription': internalDescription,
          'salePrice': salePrice,
          'finalPrice': finalPrice,
          'depositPrice': depositPrice,
          'rentPrice': rentPrice,
          'convertible': convertible,
          'maxConvertibleDeposit': maxConvertibleDeposit,
          'rentWithOwner': rentWithOwner,
          'houseInfo': houseInfo,
          'addresses': addresses.map((e) => e.toJson()).toList(),
          'vaultPlacements': vaultPlacements.map((e) => e.toJson()).toList(),
        },
      );
      var created = PropertyFileModel.fromJson(
        Map<String, dynamic>.from(res.data['data']),
      );
      for (final file in selectedUploads) {
        final form = dio.FormData.fromMap({
          'file': dio.MultipartFile.fromBytes(file.bytes, filename: file.name),
          'purpose': 'property_media',
          'targetType': 'property',
          'targetId': created.id,
          'businessId': businessId,
        });
        final upload = await _api.dio.post(
          '/uploads',
          data: form,
          options: dio.Options(contentType: 'multipart/form-data'),
        );
        final uploaded = Map<String, dynamic>.from(
          Map<String, dynamic>.from(upload.data['data'])['file'] as Map,
        );
        final confirmed = await _api.dio.post(
          '/businesses/$businessId/properties/${created.id}/media/confirm',
          data: {'fileId': uploaded['id']?.toString() ?? ''},
        );
        created = PropertyFileModel.fromJson(
          Map<String, dynamic>.from(confirmed.data['data']),
        );
      }
      selectedUploads.clear();
      await load();
      Get.back();
      Get.snackbar('فایل ثبت شد', 'فایل ملکی با موفقیت ذخیره شد.');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> updateVaults(
    PropertyFileModel file,
    List<PropertyVaultPlacementModel> vaultPlacements,
    String status,
  ) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      final res = await _api.dio.patch(
        '/businesses/$businessId/properties/${file.id}',
        data: {
          'type': file.type,
          'status': status,
          'types': file.types.isEmpty ? [file.type] : file.types,
          'title': file.title,
          'description': file.description,
          'internalDescription': file.internalDescription,
          'salePrice': file.salePrice,
          'finalPrice': file.finalPrice,
          'depositPrice': file.depositPrice,
          'rentPrice': file.rentPrice,
          'houseInfo': file.houseInfo,
          'addresses': file.addresses.map((e) => e.toJson()).toList(),
          'vaultPlacements': vaultPlacements.map((e) => e.toJson()).toList(),
        },
      );
      final updated = PropertyFileModel.fromJson(
        Map<String, dynamic>.from(res.data['data']),
      );
      final index = files.indexWhere((item) => item.id == updated.id);
      if (index >= 0) {
        files[index] = updated;
      } else {
        await load();
      }
      Get.snackbar('صندوقچه‌ها ذخیره شد', 'مسیرهای فایل به‌روزرسانی شد.');
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> requestShare(
    PropertyFileModel file, {
    double percent = 25,
  }) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/properties/${file.id}/share-requests',
        data: {'commissionPercent': percent},
      );
      await loadShareRequests();
      Get.snackbar('درخواست ثبت شد', 'درخواست مشارکت برای مالک فایل ارسال شد.');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> decideShare(
    PropertyShareRequestModel request,
    bool approve,
  ) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/property-share-requests/${request.id}/${approve ? 'approve' : 'reject'}',
      );
      await loadShareRequests();
      await load();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> receiveSharedFile(PropertyShareRequestModel request) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/property-share-requests/${request.id}/receive',
      );
      await loadShareRequests();
      await load();
      Get.snackbar(
        'فایل دریافت شد',
        'فایل مشارکتی به لیست فایل‌های شما اضافه شد.',
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> ensureLocations() async {
    await Get.find<LocationsController>().load();
  }

  bool _isVideo(PickedMedia file) {
    final ext = file.extension.toLowerCase();
    return ['mp4', 'mov', 'm4v', 'webm'].contains(ext);
  }
}
