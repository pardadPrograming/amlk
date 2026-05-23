import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../data/models.dart';
import '../business/business_controller.dart';

class ConsultantsController extends GetxController {
  ConsultantsController(this._api);
  final ApiClient _api;
  final loading = false.obs;
  final members = <Member>[].obs;
  final invitations = <Invitation>[].obs;
  final inbox = <Invitation>[].obs;
  final roleFilter = 'all'.obs;
  final query = ''.obs;

  List<Member> get filteredMembers => members.where((m) {
    final roleOk = roleFilter.value == 'all' || m.role == roleFilter.value;
    final q = query.value.trim();
    final queryOk =
        q.isEmpty || m.userPhone.contains(q) || m.userDisplayName.contains(q);
    return roleOk && queryOk;
  }).toList();

  List<Invitation> get pendingInbox =>
      inbox.where((invite) => invite.status == 'pending').toList();

  int get pendingInboxCount => pendingInbox.length;

  Future<String?> _businessId() async {
    final businessController = Get.find<BusinessController>();
    if (businessController.selected.value == null) {
      await businessController.loadFirstBusiness();
    }
    return businessController.selected.value?.id;
  }

  Future<void> load() async {
    final id = await _businessId();
    if (id == null) return;
    loading.value = true;
    try {
      final memberRes = await _api.dio.get('/businesses/$id/members');
      final inviteRes = await _api.dio.get('/businesses/$id/invitations');
      members.value = (memberRes.data['data'] as List? ?? const [])
          .map((e) => Member.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      invitations.value = (inviteRes.data['data'] as List? ?? const [])
          .map((e) => Invitation.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> invite(String phone, double commission, String role) async {
    final id = await _businessId();
    if (id == null) return;
    try {
      await _api.dio.post(
        '/businesses/$id/invitations',
        data: {'phone': phone, 'role': role, 'commissionPercent': commission},
      );
      await load();
      Get.back();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> updateMember(
    Member member, {
    String? role,
    double? commission,
    String? status,
  }) async {
    final id = await _businessId();
    if (id == null) return;
    try {
      await _api.dio.patch(
        '/businesses/$id/members/${member.id}',
        data: {
          'role': ?role,
          'commissionPercent': ?commission,
          'status': ?status,
        },
      );
      await load();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> loadInbox() async {
    try {
      final res = await _api.dio.get('/invitations/inbox');
      inbox.value = (res.data['data'] as List? ?? const [])
          .map((e) => Invitation.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    }
  }

  Future<void> accept(String id) async {
    await _api.dio.post('/invitations/$id/accept');
    await loadInbox();
  }

  Future<void> reject(String id) async {
    await _api.dio.post('/invitations/$id/reject');
    await loadInbox();
  }
}
