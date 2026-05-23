import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../data/models.dart';
import '../business/business_controller.dart';

class DashboardController extends GetxController {
  DashboardController(this._api);
  final ApiClient _api;
  final loading = false.obs;
  final stats = <String, dynamic>{}.obs;

  Future<void> load() async {
    final businessController = Get.find<BusinessController>();
    if (businessController.selected.value == null) {
      await businessController.loadFirstBusiness();
    }
    final business = businessController.selected.value;
    if (business == null) return;
    loading.value = true;
    try {
      final res = await _api.dio.get('/businesses/${business.id}/dashboard');
      stats.value = Map<String, dynamic>.from(res.data['data']);
      businessController.selected.value = Business.fromJson(
        Map<String, dynamic>.from(stats['business']),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }
}
