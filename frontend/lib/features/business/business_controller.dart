import 'package:get/get.dart';

import '../../app/app.dart';
import '../../core/api/api_client.dart';
import '../../data/models.dart';

class BusinessController extends GetxController {
  BusinessController(this._api);
  final ApiClient _api;
  final loading = false.obs;
  final selected = Rxn<Business>();

  Future<void> create({
    required String name,
    required String phone,
    required String address,
    required String workingHours,
    required String licenseNumber,
  }) async {
    loading.value = true;
    try {
      final res = await _api.dio.post(
        '/businesses',
        data: {
          'name': name,
          'phones': [phone],
          'address': address,
          'workingHours': workingHours,
          'licenseNumber': licenseNumber,
        },
      );
      selected.value = Business.fromJson(
        Map<String, dynamic>.from(res.data['data']),
      );
      Get.offAllNamed(AppRoutes.dashboard);
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadFirstBusiness() async {
    final res = await _api.dio.get('/businesses');
    final list = res.data['data'] as List? ?? const [];
    selected.value = list.isEmpty
        ? null
        : Business.fromJson(Map<String, dynamic>.from(list.first));
  }

  Future<void> leaveSelectedBusiness() async {
    final business = selected.value;
    if (business == null) return;
    loading.value = true;
    try {
      await _api.dio.delete('/businesses/${business.id}/membership');
      selected.value = null;
      await loadFirstBusiness();
      Get.back();
      if (selected.value == null) {
        Get.offAllNamed(AppRoutes.createBusiness);
      } else {
        Get.offAllNamed(AppRoutes.dashboard);
      }
      Get.snackbar('خروج انجام شد', 'دسترسی شما به این کسب‌وکار قطع شد.');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }
}
