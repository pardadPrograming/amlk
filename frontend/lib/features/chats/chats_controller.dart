import 'package:dio/dio.dart' as dio;
import 'package:get/get.dart';

import '../../core/api/api_client.dart';
import '../../core/media/media_picker.dart';
import '../../data/models.dart';
import '../auth/auth_controller.dart';
import '../business/business_controller.dart';

enum ChatSection {
  privateChats,
  fileChats,
  businessFileChats,
  personalVaults,
  joinedVaults,
  businessVaults,
}

class ChatsController extends GetxController {
  ChatsController(this._api);

  final ApiClient _api;
  final loading = false.obs;
  final channels = <ChannelListModel>[].obs;
  final personalVaults = <ChannelVaultModel>[].obs;
  final businessVaults = <ChannelVaultModel>[].obs;
  final canManageBusinessVaults = false.obs;
  final threadMessages = <ChannelMessageModel>[].obs;
  final threadMembers = <ChannelMemberModel>[].obs;
  final vaultFiles = <ChannelVaultFileModel>[].obs;
  final contactCategories = <String>[].obs;
  final sending = false.obs;
  final threadLoading = false.obs;
  final loadingOlderMessages = false.obs;
  final loadingNewerMessages = false.obs;
  final threadUnreadCount = 0.obs;
  final firstUnreadMessageId = ''.obs;
  final hasOlderMessages = false.obs;
  final hasNewerMessages = false.obs;
  int _threadOffset = 0;
  int _threadTotal = 0;
  static const int _messagePageLimit = 30;

  String? get _businessId => Get.find<BusinessController>().selected.value?.id;
  String get currentUserId => Get.find<AuthController>().user.value?.id ?? '';

  ChannelListModel? channelById(String channelId) {
    for (final channel in channels) {
      if (channel.id == channelId) {
        return channel;
      }
    }
    return null;
  }

  Future<void> load() async {
    loading.value = true;
    try {
      final channelRes = await _api.dio.get('/channels');
      channels.assignAll(
        (channelRes.data['data'] as List? ?? const []).map(
          (e) => ChannelListModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );

      final userVaultRes = await _api.dio.get('/vaults');
      personalVaults.assignAll(
        (userVaultRes.data['data'] as List? ?? const []).map(
          (e) => ChannelVaultModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );

      await _loadBusinessVaults();
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      loading.value = false;
    }
  }

  Future<void> _loadBusinessVaults() async {
    final businessId = _businessId;
    if (businessId == null) {
      canManageBusinessVaults.value = false;
      businessVaults.clear();
      return;
    }
    try {
      final businessVaultRes = await _api.dio.get(
        '/businesses/$businessId/vaults',
      );
      canManageBusinessVaults.value = true;
      businessVaults.assignAll(
        (businessVaultRes.data['data'] as List? ?? const []).map(
          (e) => ChannelVaultModel.fromJson(Map<String, dynamic>.from(e)),
        ),
      );
    } catch (e) {
      canManageBusinessVaults.value = false;
      businessVaults.clear();
      final statusCode = e is dio.DioException ? e.response?.statusCode : null;
      if (statusCode != 401 && statusCode != 403) {
        Get.snackbar('خطا', _api.message(e));
      }
    }
  }

  List<ChatSection> get visibleSections => [
    ChatSection.privateChats,
    ChatSection.fileChats,
    if (canManageBusinessVaults.value) ChatSection.businessFileChats,
    ChatSection.personalVaults,
    ChatSection.joinedVaults,
    if (canManageBusinessVaults.value) ChatSection.businessVaults,
  ];

  List<ChannelListModel> get privateChats => channels
      .where((item) => item.type == 'private' || item.type == 'user_main')
      .toList();

  List<ChannelListModel> get fileChats =>
      channels.where((item) => item.type == 'user_vault').toList();

  List<ChannelListModel> get businessFileChats => canManageBusinessVaults.value
      ? channels
            .where(
              (item) =>
                  item.type == 'business_vault' &&
                  businessVaults.any((vault) => vault.id == item.vaultId),
            )
            .toList()
      : const [];

  List<ChannelListModel> get joinedVaultChannels {
    final ownedIds = personalVaults.map((item) => item.id).toSet();
    final businessIds = businessVaults.map((item) => item.id).toSet();
    return channels
        .where(
          (item) =>
              (item.type == 'user_vault' || item.type == 'business_vault') &&
              !ownedIds.contains(item.vaultId) &&
              !businessIds.contains(item.vaultId),
        )
        .toList();
  }

  List<AccessibleVaultTargetModel> get accessibleVaultTargets {
    final seen = <String>{};
    final result = <AccessibleVaultTargetModel>[];
    void add(String channelId, String title, String subtitle) {
      if (channelId.isEmpty || seen.contains(channelId)) {
        return;
      }
      seen.add(channelId);
      result.add(
        AccessibleVaultTargetModel(
          channelId: channelId,
          title: title.isEmpty ? 'صندوقچه' : title,
          subtitle: subtitle,
        ),
      );
    }

    for (final vault in personalVaults) {
      add(vault.channelId, vault.title, 'صندوقچه شخصی');
    }
    for (final vault in businessVaults) {
      add(vault.channelId, vault.title, 'صندوقچه املاک');
    }
    for (final channel in joinedVaultChannels) {
      add(
        channel.id,
        channel.title,
        channel.isBusinessVault ? 'صندوقچه عضو شده املاک' : 'صندوقچه عضو شده',
      );
    }
    return result;
  }

  ChannelMemberModel? memberForUserId(String userId) {
    if (userId.isEmpty) {
      return null;
    }
    for (final member in threadMembers) {
      if (member.userId == userId) {
        return member;
      }
    }
    return null;
  }

  String memberTitle(ChannelMemberModel member) =>
      member.displayName.isEmpty ? member.phone : member.displayName;

  Future<void> loadMessages(
    String channelId, {
    bool fromUnread = false,
    int? offset,
    bool appendOlder = false,
    bool prependNewer = false,
  }) async {
    final res = await _api.dio.get(
      '/channels/$channelId/messages',
      queryParameters: {
        'limit': _messagePageLimit,
        'window': true,
        if (fromUnread) 'fromUnread': true,
        'offset': ?offset,
      },
    );
    final data = res.data['data'];
    final page = data is Map ? Map<String, dynamic>.from(data) : {};
    final rawItems = page['items'] as List? ?? const [];
    final items = rawItems
        .map((e) => ChannelMessageModel.fromJson(Map<String, dynamic>.from(e)))
        .toList();
    _threadTotal = (page['total'] as num? ?? items.length).toInt();
    if (!appendOlder && !prependNewer) {
      threadUnreadCount.value = (page['unreadCount'] as num? ?? 0).toInt();
      firstUnreadMessageId.value =
          page['firstUnreadMessageId']?.toString() ?? '';
    }
    if (appendOlder) {
      _appendMessages(items);
    } else if (prependNewer) {
      _prependMessages(items);
      _threadOffset = (page['offset'] as num? ?? 0).toInt();
    } else {
      threadMessages.assignAll(items);
      _threadOffset = (page['offset'] as num? ?? 0).toInt();
    }
    hasOlderMessages.value =
        (page['hasOlder'] as bool?) ??
        _threadOffset + threadMessages.length < _threadTotal;
    hasNewerMessages.value = (page['hasNewer'] as bool?) ?? _threadOffset > 0;
  }

  Future<void> loadMembers(String channelId) async {
    final res = await _api.dio.get('/channels/$channelId/members');
    threadMembers.assignAll(
      (res.data['data'] as List? ?? const []).map(
        (e) => ChannelMemberModel.fromJson(Map<String, dynamic>.from(e)),
      ),
    );
  }

  Future<void> loadThread(String channelId) async {
    threadLoading.value = true;
    try {
      await Future.wait([
        loadMessages(channelId, fromUnread: true),
        loadMembers(channelId),
      ]);
    } finally {
      threadLoading.value = false;
    }
  }

  Future<void> loadOlderMessages(String channelId) async {
    if (loadingOlderMessages.value || !hasOlderMessages.value) {
      return;
    }
    loadingOlderMessages.value = true;
    try {
      await loadMessages(
        channelId,
        offset: _threadOffset + threadMessages.length,
        appendOlder: true,
      );
    } finally {
      loadingOlderMessages.value = false;
    }
  }

  Future<void> loadNewerMessages(String channelId) async {
    if (loadingNewerMessages.value || !hasNewerMessages.value) {
      return;
    }
    loadingNewerMessages.value = true;
    try {
      final nextOffset = _threadOffset - _messagePageLimit;
      await loadMessages(
        channelId,
        offset: nextOffset < 0 ? 0 : nextOffset,
        prependNewer: true,
      );
    } finally {
      loadingNewerMessages.value = false;
    }
  }

  Future<void> loadVaultFiles(String channelId) async {
    final res = await _api.dio.get('/channels/$channelId/vault/files');
    vaultFiles.assignAll(
      (res.data['data'] as List? ?? const []).map(
        (e) => ChannelVaultFileModel.fromJson(Map<String, dynamic>.from(e)),
      ),
    );
  }

  Future<void> loadContactCategories() async {
    if (contactCategories.isNotEmpty) {
      return;
    }
    final res = await _api.dio.get('/contact-tags');
    final list = res.data['data'] as List? ?? const [];
    contactCategories.assignAll(list.map((item) => item.toString()));
  }

  Future<ProfileContactCategoryResultModel> addProfileToContactCategories(
    ChannelMemberModel member,
    List<String> tags,
  ) async {
    final businessId = _businessId;
    if (businessId == null) {
      throw Exception('برای ثبت مخاطب باید املاک فعال انتخاب شده باشد');
    }
    if (member.phone.isEmpty) {
      throw Exception('شماره تماس برای ثبت مخاطب وجود ندارد');
    }
    final res = await _api.dio.post(
      '/businesses/$businessId/contacts/profile-categories',
      data: {
        'phone': member.phone,
        'displayName': memberTitle(member),
        'tags': tags,
      },
    );
    return ProfileContactCategoryResultModel.fromJson(
      Map<String, dynamic>.from(res.data['data']),
    );
  }

  Future<ChannelListModel> startPrivateChat(ChannelMemberModel member) async {
    if (member.phone.isEmpty) {
      throw Exception('شماره تماس برای شروع چت وجود ندارد');
    }
    final res = await _api.dio.post(
      '/channels/private',
      data: {'phone': member.phone},
    );
    final channel = ChannelListModel.fromJson(
      Map<String, dynamic>.from(res.data['data']),
    );
    final existing = channels.indexWhere((item) => item.id == channel.id);
    if (existing >= 0) {
      channels[existing] = channel;
    } else {
      channels.insert(0, channel);
    }
    return channel;
  }

  Future<void> addMemberToVault(
    AccessibleVaultTargetModel vault,
    ChannelMemberModel member,
  ) async {
    if (member.phone.isEmpty) {
      throw Exception('شماره تماس برای افزودن به صندوقچه وجود ندارد');
    }
    await _api.dio.post(
      '/channels/${vault.channelId}/invites',
      data: {'phone': member.phone},
    );
  }

  Future<ChannelMediaModel> uploadChannelMedia(
    String channelId,
    PickedMedia file,
  ) async {
    final form = dio.FormData.fromMap({
      'file': dio.MultipartFile.fromBytes(file.bytes, filename: file.name),
      'purpose': 'channel_media',
      'targetType': 'channel',
      'targetId': channelId,
    });
    final res = await _api.dio.post(
      '/uploads',
      data: form,
      options: dio.Options(contentType: 'multipart/form-data'),
    );
    final data = Map<String, dynamic>.from(res.data['data']);
    final uploaded = Map<String, dynamic>.from(data['file'] as Map);
    return ChannelMediaModel(
      id: uploaded['id']?.toString() ?? '',
      fileId: uploaded['id']?.toString() ?? '',
      kind: uploaded['kind']?.toString() ?? 'file',
      url: uploaded['url']?.toString() ?? '',
      contentType: uploaded['contentType']?.toString() ?? '',
      size: int.tryParse(uploaded['size']?.toString() ?? '') ?? 0,
      createdAt: uploaded['createdAt']?.toString() ?? '',
      expiresAt: uploaded['expiresAt']?.toString() ?? '',
    );
  }

  Future<void> sendMessage(
    String channelId, {
    String text = '',
    String caption = '',
    List<ChannelMediaModel> media = const [],
    String vaultFileRefId = '',
  }) async {
    sending.value = true;
    try {
      await _api.dio.post(
        '/channels/$channelId/messages',
        data: {
          'text': text.trim(),
          'caption': caption.trim(),
          'media': media.map((item) => item.toJson()).toList(),
          if (vaultFileRefId.isNotEmpty) 'vaultFileRefId': vaultFileRefId,
        },
      );
      await loadMessages(channelId, offset: 0);
    } finally {
      sending.value = false;
    }
  }

  Future<void> editMessage(
    String channelId,
    String messageId, {
    String text = '',
    String caption = '',
  }) async {
    sending.value = true;
    try {
      final res = await _api.dio.patch(
        '/channels/$channelId/messages/$messageId',
        data: {'text': text.trim(), 'caption': caption.trim()},
      );
      final updated = ChannelMessageModel.fromJson(
        Map<String, dynamic>.from(res.data['data']),
      );
      final index = threadMessages.indexWhere((item) => item.id == messageId);
      if (index >= 0) {
        threadMessages[index] = updated;
      } else {
        await loadMessages(channelId);
      }
    } finally {
      sending.value = false;
    }
  }

  Future<void> deleteMessage(String channelId, String messageId) async {
    sending.value = true;
    try {
      await _api.dio.delete('/channels/$channelId/messages/$messageId');
      threadMessages.removeWhere((item) => item.id == messageId);
    } finally {
      sending.value = false;
    }
  }

  Future<void> pasteFiles(String channelId, List<PickedMedia> files) async {
    if (files.isEmpty) {
      return;
    }
    sending.value = true;
    try {
      final media = <ChannelMediaModel>[];
      for (final file in files) {
        media.add(await uploadChannelMedia(channelId, file));
      }
      await _api.dio.post(
        '/channels/$channelId/messages',
        data: {'media': media.map((item) => item.toJson()).toList()},
      );
      await loadMessages(channelId, offset: 0);
      Get.snackbar(
        'فایل ارسال شد',
        '${files.length} فایل از کلیپ‌بورد ارسال شد.',
      );
    } catch (e) {
      Get.snackbar('خطا', _api.message(e));
    } finally {
      sending.value = false;
    }
  }

  Future<void> startFileChat(
    String channelId,
    ChannelVaultFileModel file,
  ) async {
    await sendMessage(channelId, vaultFileRefId: file.id);
  }

  void _appendMessages(List<ChannelMessageModel> items) {
    final seen = threadMessages.map((item) => item.id).toSet();
    threadMessages.addAll(items.where((item) => seen.add(item.id)));
  }

  void _prependMessages(List<ChannelMessageModel> items) {
    final seen = threadMessages.map((item) => item.id).toSet();
    threadMessages.insertAll(0, items.where((item) => seen.add(item.id)));
  }
}

class ChannelListModel {
  ChannelListModel({
    required this.id,
    required this.type,
    this.title = '',
    this.ownerUserId = '',
    this.businessId = '',
    this.vaultId = '',
    this.updatedAt = '',
  });

  final String id;
  final String type;
  final String title;
  final String ownerUserId;
  final String businessId;
  final String vaultId;
  final String updatedAt;

  bool get isBusinessVault => type == 'business_vault';

  factory ChannelListModel.fromJson(Map<String, dynamic> json) =>
      ChannelListModel(
        id: json['id']?.toString() ?? '',
        type: json['type']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        ownerUserId: json['ownerUserId']?.toString() ?? '',
        businessId: json['businessId']?.toString() ?? '',
        vaultId: json['vaultId']?.toString() ?? '',
        updatedAt: json['updatedAt']?.toString() ?? '',
      );
}

class ChannelMediaModel {
  ChannelMediaModel({
    required this.id,
    required this.fileId,
    required this.kind,
    required this.url,
    this.contentType = '',
    this.size = 0,
    this.createdAt = '',
    this.expiresAt = '',
  });

  final String id;
  final String fileId;
  final String kind;
  final String url;
  final String contentType;
  final int size;
  final String createdAt;
  final String expiresAt;

  factory ChannelMediaModel.fromJson(Map<String, dynamic> json) =>
      ChannelMediaModel(
        id: json['id']?.toString() ?? '',
        fileId: json['fileId']?.toString() ?? '',
        kind: json['kind']?.toString() ?? '',
        url: json['url']?.toString() ?? '',
        contentType: json['contentType']?.toString() ?? '',
        size: int.tryParse(json['size']?.toString() ?? '') ?? 0,
        createdAt: json['createdAt']?.toString() ?? '',
        expiresAt: json['expiresAt']?.toString() ?? '',
      );

  Map<String, dynamic> toJson() => {
    'id': id,
    'fileId': fileId,
    'kind': kind,
    'url': url,
    'contentType': contentType,
    'size': size,
    'createdAt': createdAt,
    if (expiresAt.isNotEmpty) 'expiresAt': expiresAt,
  };
}

class ChannelMessageModel {
  ChannelMessageModel({
    required this.id,
    this.authorId = '',
    this.authorName = '',
    this.text = '',
    this.caption = '',
    this.media = const [],
    this.vaultFileRef,
    this.seenBy = const [],
    this.createdAt = '',
  });

  final String id;
  final String authorId;
  final String authorName;
  final String text;
  final String caption;
  final List<ChannelMediaModel> media;
  final ChannelVaultFilePreviewModel? vaultFileRef;
  final List<ChannelMessageSeenModel> seenBy;
  final String createdAt;

  bool isMine(String currentUserId) =>
      authorId.isNotEmpty && authorId == currentUserId;

  bool seenByOther(String currentUserId) => seenBy.any(
    (item) => item.userId.isNotEmpty && item.userId != currentUserId,
  );

  factory ChannelMessageModel.fromJson(
    Map<String, dynamic> json,
  ) => ChannelMessageModel(
    id: json['id']?.toString() ?? '',
    authorId: json['authorId']?.toString() ?? '',
    authorName: json['authorName']?.toString() ?? '',
    text: json['text']?.toString() ?? '',
    caption: json['caption']?.toString() ?? '',
    media: (json['media'] as List? ?? const [])
        .map((e) => ChannelMediaModel.fromJson(Map<String, dynamic>.from(e)))
        .toList(),
    vaultFileRef: json['vaultFileRef'] is Map
        ? ChannelVaultFilePreviewModel.fromJson(
            Map<String, dynamic>.from(json['vaultFileRef']),
          )
        : null,
    seenBy: (json['seenBy'] as List? ?? const [])
        .map(
          (e) => ChannelMessageSeenModel.fromJson(Map<String, dynamic>.from(e)),
        )
        .toList(),
    createdAt: json['createdAt']?.toString() ?? '',
  );
}

class ChannelMessageSeenModel {
  ChannelMessageSeenModel({required this.userId, this.seenAt = ''});

  final String userId;
  final String seenAt;

  factory ChannelMessageSeenModel.fromJson(Map<String, dynamic> json) =>
      ChannelMessageSeenModel(
        userId: json['userId']?.toString() ?? '',
        seenAt: json['seenAt']?.toString() ?? '',
      );
}

class ChannelVaultFileModel {
  ChannelVaultFileModel({
    required this.id,
    this.title = '',
    this.kind = '',
    this.url = '',
    this.contentType = '',
    this.size = 0,
    this.propertyStatus = '',
    this.commissionPercent = 0,
    this.createdAt = '',
  });

  final String id;
  final String title;
  final String kind;
  final String url;
  final String contentType;
  final int size;
  final String propertyStatus;
  final double commissionPercent;
  final String createdAt;

  factory ChannelVaultFileModel.fromJson(Map<String, dynamic> json) =>
      ChannelVaultFileModel(
        id: json['id']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        kind: json['kind']?.toString() ?? '',
        url: json['url']?.toString() ?? '',
        contentType: json['contentType']?.toString() ?? '',
        size: int.tryParse(json['size']?.toString() ?? '') ?? 0,
        propertyStatus: json['propertyStatus']?.toString() ?? '',
        commissionPercent:
            double.tryParse(json['commissionPercent']?.toString() ?? '') ?? 0,
        createdAt: json['createdAt']?.toString() ?? '',
      );
}

class ChannelVaultFilePreviewModel {
  ChannelVaultFilePreviewModel({
    required this.id,
    this.title = '',
    this.kind = '',
    this.url = '',
    this.contentType = '',
    this.size = 0,
    this.propertyStatus = '',
    this.commissionPercent = 0,
  });

  final String id;
  final String title;
  final String kind;
  final String url;
  final String contentType;
  final int size;
  final String propertyStatus;
  final double commissionPercent;

  factory ChannelVaultFilePreviewModel.fromJson(Map<String, dynamic> json) =>
      ChannelVaultFilePreviewModel(
        id: json['id']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        kind: json['kind']?.toString() ?? '',
        url: json['url']?.toString() ?? '',
        contentType: json['contentType']?.toString() ?? '',
        size: int.tryParse(json['size']?.toString() ?? '') ?? 0,
        propertyStatus: json['propertyStatus']?.toString() ?? '',
        commissionPercent:
            double.tryParse(json['commissionPercent']?.toString() ?? '') ?? 0,
      );
}

class AccessibleVaultTargetModel {
  AccessibleVaultTargetModel({
    required this.channelId,
    required this.title,
    this.subtitle = '',
  });

  final String channelId;
  final String title;
  final String subtitle;
}

class ProfileContactCategoryResultModel {
  ProfileContactCategoryResultModel({
    this.autoConsultant = false,
    this.existing = false,
  });

  final bool autoConsultant;
  final bool existing;

  factory ProfileContactCategoryResultModel.fromJson(
    Map<String, dynamic> json,
  ) => ProfileContactCategoryResultModel(
    autoConsultant: json['autoConsultant'] == true,
    existing: json['existing'] == true,
  );
}

class ChannelMemberModel {
  ChannelMemberModel({
    required this.id,
    this.userId = '',
    this.phone = '',
    this.displayName = '',
    this.isOnline = false,
    this.lastSeenAt = '',
  });

  final String id;
  final String userId;
  final String phone;
  final String displayName;
  final bool isOnline;
  final String lastSeenAt;

  factory ChannelMemberModel.fromJson(Map<String, dynamic> json) =>
      ChannelMemberModel(
        id: json['id']?.toString() ?? '',
        userId: json['userId']?.toString() ?? '',
        phone: json['phone']?.toString() ?? '',
        displayName: json['displayName']?.toString() ?? '',
        isOnline: json['isOnline'] == true,
        lastSeenAt: json['lastSeenAt']?.toString() ?? '',
      );
}
