import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../data/models.dart';
import '../business/business_controller.dart';

class CustomersController extends GetxController {
  CustomersController(this._api);

  final ApiClient _api;
  final loading = false.obs;
  final contacts = <ContactModel>[].obs;
  final systemTags = <String>[].obs;
  final areas = <AreaNode>[].obs;
  final requestMatches = <String, List<PropertyMatchResultModel>>{}.obs;
  final loadingMatches = <String>{}.obs;
  final query = ''.obs;
  final selectedTag = ''.obs;

  String? get _businessId => Get.find<BusinessController>().selected.value?.id;

  Future<void> load() async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await loadTags();
      await loadLocations();
      final res = await _api.dio.get(
        '/businesses/$businessId/contacts',
        queryParameters: {
          if (query.value.trim().isNotEmpty) 'q': query.value.trim(),
          if (selectedTag.value.isNotEmpty) 'tag': selectedTag.value,
        },
      );
      final list = res.data['data'] as List? ?? const [];
      contacts.assignAll(
        list.map((e) => ContactModel.fromJson(Map<String, dynamic>.from(e))),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadTags() async {
    final res = await _api.dio.get('/contact-tags');
    final list = res.data['data'] as List? ?? const [];
    systemTags.assignAll(list.map((e) => e.toString()));
  }

  Future<void> loadLocations() async {
    final businessId = _businessId;
    if (businessId == null) return;
    final res = await _api.dio.get('/businesses/$businessId/locations');
    final list = res.data['data'] as List? ?? const [];
    areas.assignAll(
      list.map((e) => AreaNode.fromJson(Map<String, dynamic>.from(e))),
    );
  }

  Future<void> create({
    required String firstName,
    required String lastName,
    required String company,
    required List<ContactPhoneModel> phones,
    required List<String> tags,
    required List<ContactPropertyRefModel> properties,
    required List<ContactRequestModel> requests,
    required String note,
  }) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.post(
        '/businesses/$businessId/contacts',
        data: {
          'firstName': firstName,
          'lastName': lastName,
          'company': company,
          'phones': phones.map((e) => e.toJson()).toList(),
          'tags': tags,
          'properties': properties.map((e) => e.toJson()).toList(),
          'requests': requests.map((e) => e.toJson()).toList(),
          'note': note,
        },
      );
      await load();
      Get.back();
      Get.snackbar('ثبت شد', 'مخاطب در دفترچه تلفن ثبت شد');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> updateContact({
    required String contactId,
    required String firstName,
    required String lastName,
    required String company,
    required List<ContactPhoneModel> phones,
    required List<String> tags,
    required List<ContactPropertyRefModel> properties,
    required List<ContactRequestModel> requests,
    required String note,
  }) async {
    final businessId = _businessId;
    if (businessId == null) return;
    loading.value = true;
    try {
      await _api.dio.patch(
        '/businesses/$businessId/contacts/$contactId',
        data: {
          'firstName': firstName,
          'lastName': lastName,
          'company': company,
          'phones': phones.map((e) => e.toJson()).toList(),
          'tags': tags,
          'properties': properties.map((e) => e.toJson()).toList(),
          'requests': requests.map((e) => e.toJson()).toList(),
          'note': note,
        },
      );
      await load();
      Get.back();
      Get.snackbar('ذخیره شد', 'تغییرات مخاطب و درخواست‌ها ثبت شد');
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadRequestMatches(String contactId, String requestId) async {
    final businessId = _businessId;
    if (businessId == null || contactId.isEmpty || requestId.isEmpty) return;
    final key = '$contactId:$requestId';
    loadingMatches.add(key);
    try {
      final res = await _api.dio.get(
        '/businesses/$businessId/contacts/$contactId/requests/$requestId/matches',
      );
      final list = res.data['data'] as List? ?? const [];
      requestMatches[key] = list
          .map(
            (e) =>
                PropertyMatchResultModel.fromJson(Map<String, dynamic>.from(e)),
          )
          .toList();
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    } finally {
      loadingMatches.remove(key);
    }
  }
}
