import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../app/app.dart';
import '../../core/api/api_client.dart';
import '../../core/storage/token_storage.dart';
import '../../data/models.dart';
import '../consultants/consultants_controller.dart';

class AuthController extends GetxController {
  AuthController(this._api, this._storage);
  final ApiClient _api;
  final TokenStorage _storage;

  final loading = false.obs;
  final phone = ''.obs;
  final devCode = ''.obs;
  final user = Rxn<AppUser>();
  final sessions = <UserSession>[].obs;
  final cities = <CityModel>[].obs;
  final lastActivityAt = Rxn<DateTime>();

  Future<void> requestOtp(String value) async {
    loading.value = true;
    try {
      final res = await _api.dio.post(
        '/auth/request-otp',
        data: {'phone': value},
      );
      phone.value = res.data['data']['phone'];
      devCode.value = res.data['data']['developmentCode']?.toString() ?? '';
      Get.toNamed(AppRoutes.otp);
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadLatestTestOtp() async {
    loading.value = true;
    try {
      final res = await _api.dio.get('/test/latest-otp');
      final data = res.data['data'];
      phone.value = data['phone']?.toString() ?? phone.value;
      devCode.value = data['code']?.toString() ?? '';
      Get.snackbar('OTP تست', 'آخرین کد تست دریافت شد');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> verifyOtp(String code) async {
    if (phone.value.isEmpty) {
      Get.snackbar(
        'خطا',
        'شماره موبایل در این نشست ثبت نشده است. دوباره کد بگیرید.',
      );
      return;
    }
    loading.value = true;
    try {
      final res = await _api.dio.post(
        '/auth/verify-otp',
        data: {'phone': phone.value, 'code': code},
      );
      final data = res.data['data'];
      await _storage.save(data['accessToken'], data['refreshToken']);
      user.value = AppUser.fromJson(Map<String, dynamic>.from(data['user']));
      await bootstrap();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> bootstrap() async {
    final res = await _api.dio.get('/auth/me');
    final data = res.data['data'];
    user.value = AppUser.fromJson(Map<String, dynamic>.from(data['user']));
    final businesses = data['businesses'] as List? ?? const [];
    final consultants = Get.find<ConsultantsController>();
    await consultants.loadInbox();

    if ((user.value?.firstName ?? '').isEmpty ||
        (user.value?.lastName ?? '').isEmpty ||
        (user.value?.cityId ?? '').isEmpty) {
      Get.offAllNamed(AppRoutes.profile);
    } else if (businesses.isEmpty) {
      Get.offAllNamed(AppRoutes.createBusiness);
    } else {
      Get.offAllNamed(
        AppRoutes.dashboard,
        arguments: Map<String, dynamic>.from(businesses.first),
      );
    }
    _notifyPendingInvitations(consultants.pendingInboxCount);
  }

  void _notifyPendingInvitations(int count) {
    if (count == 0) return;
    Get.snackbar(
      'دعوت‌نامه در انتظار',
      count == 1 ? 'یک دعوت‌نامه جدید دارید.' : '$count دعوت‌نامه جدید دارید.',
      mainButton: TextButton(
        onPressed: () => Get.toNamed(AppRoutes.inbox),
        child: const Text('مشاهده'),
      ),
    );
  }

  Future<void> loadCities() async {
    try {
      final res = await _api.dio.get('/catalog/cities');
      final list = res.data['data'] as List? ?? const [];
      cities.assignAll(
        list.map((e) => CityModel.fromJson(Map<String, dynamic>.from(e))),
      );
    } catch (e) {
      Get.snackbar('Ø®Ø·Ø§', _api.message(e));
    }
  }

  Future<void> completeProfile(
    String firstName,
    String lastName,
    String cityId,
  ) async {
    loading.value = true;
    try {
      final res = await _api.dio.patch(
        '/auth/profile',
        data: {'firstName': firstName, 'lastName': lastName, 'cityId': cityId},
      );
      user.value = AppUser.fromJson(
        Map<String, dynamic>.from(res.data['data']),
      );
      await bootstrap();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> loadSecurityProfile() async {
    loading.value = true;
    try {
      final res = await _api.dio.get('/auth/security');
      final data = Map<String, dynamic>.from(res.data['data']);
      sessions.assignAll(
        (data['sessions'] as List? ?? const [])
            .map((e) => UserSession.fromJson(Map<String, dynamic>.from(e)))
            .toList(),
      );
      lastActivityAt.value = DateTime.tryParse(
        data['lastActivityAt']?.toString() ?? '',
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> revokeSession(String sessionId) async {
    loading.value = true;
    try {
      await _api.dio.delete('/auth/sessions/$sessionId');
      sessions.removeWhere((session) => session.id == sessionId);
      Get.snackbar('نشست غیرفعال شد', 'دسترسی این دستگاه قطع شد.');
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }
}
