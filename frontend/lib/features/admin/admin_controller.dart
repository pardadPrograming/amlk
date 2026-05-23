import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../data/models.dart';

class AdminController extends GetxController {
  AdminController(this._api);

  final ApiClient _api;
  final loading = false.obs;
  final account = Rxn<AdminAccountModel>();
  final cities = <CityModel>[].obs;
  final suggestions = <LocationSuggestionModel>[].obs;
  final users = <AppUser>[].obs;
  final businesses = <Business>[].obs;
  final settings = PlatformSettingsModel().obs;

  bool get isAdmin => account.value != null;

  Future<void> load() async {
    loading.value = true;
    try {
      await Future.wait([loadMe(), loadCities()]);
      if (isAdmin) {
        await Future.wait([
          loadSuggestions(),
          loadUsers(),
          loadBusinesses(),
          loadSettings(),
        ]);
      }
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadMe() async {
    try {
      final res = await _api.dio.get('/admin/me');
      account.value = AdminAccountModel.fromJson(
        Map<String, dynamic>.from(res.data['data'] as Map),
      );
    } catch (_) {
      account.value = null;
    }
  }

  Future<void> loadCities() async {
    final res = await _api.dio.get('/catalog/cities');
    final list = res.data['data'] as List? ?? const [];
    cities.assignAll(
      list.map((e) => CityModel.fromJson(Map<String, dynamic>.from(e))),
    );
  }

  Future<void> loadSuggestions() async {
    final res = await _api.dio.get(
      '/admin/location-suggestions',
      queryParameters: {'status': 'pending'},
    );
    final list = res.data['data'] as List? ?? const [];
    suggestions.assignAll(
      list.map(
        (e) => LocationSuggestionModel.fromJson(Map<String, dynamic>.from(e)),
      ),
    );
  }

  Future<void> loadUsers() async {
    final res = await _api.dio.get('/admin/users');
    final list = res.data['data'] as List? ?? const [];
    users.assignAll(
      list.map((e) => AppUser.fromJson(Map<String, dynamic>.from(e))),
    );
  }

  Future<void> loadBusinesses() async {
    final res = await _api.dio.get('/admin/businesses');
    final list = res.data['data'] as List? ?? const [];
    businesses.assignAll(
      list.map((e) => Business.fromJson(Map<String, dynamic>.from(e))),
    );
  }

  Future<void> loadSettings() async {
    final res = await _api.dio.get('/admin/settings');
    settings.value = PlatformSettingsModel.fromJson(
      Map<String, dynamic>.from(res.data['data'] as Map),
    );
  }

  Future<void> saveSettings({
    required String otpApiKey,
    required String serviceSmsApiKey,
  }) async {
    loading.value = true;
    try {
      final res = await _api.dio.patch(
        '/admin/settings',
        data: {
          'otpApiKey': otpApiKey.trim(),
          'serviceSmsApiKey': serviceSmsApiKey.trim(),
        },
      );
      settings.value = PlatformSettingsModel.fromJson(
        Map<String, dynamic>.from(res.data['data'] as Map),
      );
      Get.snackbar('ذخیره شد', 'تنظیمات API Keys مجموعه به‌روز شد');
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> createCity(String name) async {
    if (name.trim().isEmpty) return;
    loading.value = true;
    try {
      await _api.dio.post('/admin/cities', data: {'name': name.trim()});
      await loadCities();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> approveSuggestion(String id) async {
    loading.value = true;
    try {
      await _api.dio.post('/admin/location-suggestions/$id/approve');
      await loadSuggestions();
      await loadCities();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> rejectSuggestion(String id, String note) async {
    loading.value = true;
    try {
      await _api.dio.post(
        '/admin/location-suggestions/$id/reject',
        data: {'note': note},
      );
      await loadSuggestions();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }
}
