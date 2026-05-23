import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../data/models.dart';
import '../business/business_controller.dart';

class LocationsController extends GetxController {
  LocationsController(this._api);

  final ApiClient _api;
  final loading = false.obs;
  final areas = <AreaNode>[].obs;
  final systemCities = <CityModel>[].obs;

  String? get _businessId => Get.find<BusinessController>().selected.value?.id;

  Future<void> load() async {
    final businessId = _businessId;
    if (businessId == null || businessId.isEmpty) return;
    loading.value = true;
    try {
      await loadSystemCities();
      final res = await _api.dio.get('/businesses/$businessId/locations');
      final list = res.data['data'] as List? ?? const [];
      areas.assignAll(
        list.map((e) => AreaNode.fromJson(Map<String, dynamic>.from(e))),
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadSystemCities() async {
    final res = await _api.dio.get('/catalog/cities');
    final list = res.data['data'] as List? ?? const [];
    systemCities.assignAll(
      list.map((e) => CityModel.fromJson(Map<String, dynamic>.from(e))),
    );
  }

  Future<void> suggestLocation({
    required String cityId,
    required String type,
    required String name,
    String parentAreaId = '',
    String parentStreetId = '',
    String manualParentName = '',
  }) async {
    loading.value = true;
    try {
      await _api.dio.post(
        '/location-suggestions',
        data: {
          'cityId': cityId,
          'type': type,
          'name': name,
          'parentAreaId': parentAreaId,
          'parentStreetId': parentStreetId,
          'manualParentName': manualParentName,
        },
      );
      Get.snackbar('ثبت شد', 'پیشنهاد لوکیشن برای بررسی مدیریت ارسال شد');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> addArea(String name) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () =>
          _api.dio.post('/businesses/$businessId/areas', data: {'name': name}),
    );
  }

  Future<void> addStreet(String areaId, String name) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () => _api.dio.post(
        '/businesses/$businessId/areas/$areaId/streets',
        data: {'name': name},
      ),
    );
  }

  Future<void> addNeighborhood(
    String areaId,
    String streetId,
    String name,
  ) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () => _api.dio.post(
        '/businesses/$businessId/areas/$areaId/streets/$streetId/neighborhoods',
        data: {'name': name},
      ),
    );
  }

  Future<void> deleteArea(String areaId) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () => _api.dio.delete('/businesses/$businessId/areas/$areaId'),
    );
  }

  Future<void> deleteStreet(String areaId, String streetId) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () => _api.dio.delete(
        '/businesses/$businessId/areas/$areaId/streets/$streetId',
      ),
    );
  }

  Future<void> deleteNeighborhood(
    String areaId,
    String streetId,
    String neighborhoodId,
  ) async {
    final businessId = _businessId;
    if (businessId == null) return;
    await _mutate(
      () => _api.dio.delete(
        '/businesses/$businessId/areas/$areaId/streets/$streetId/neighborhoods/$neighborhoodId',
      ),
    );
  }

  Future<void> _mutate(Future<dynamic> Function() action) async {
    loading.value = true;
    try {
      await action();
      await load();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }
}
