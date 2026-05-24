class AppUser {
  AppUser({
    required this.id,
    required this.phone,
    this.firstName = '',
    this.lastName = '',
    this.displayName = '',
    this.cityId = '',
  });
  final String id;
  final String phone;
  final String firstName;
  final String lastName;
  final String displayName;
  final String cityId;

  factory AppUser.fromJson(Map<String, dynamic> json) => AppUser(
    id: json['id']?.toString() ?? '',
    phone: json['phone']?.toString() ?? '',
    firstName: json['firstName']?.toString() ?? '',
    lastName: json['lastName']?.toString() ?? '',
    displayName: json['displayName']?.toString() ?? '',
    cityId: json['cityId']?.toString() ?? '',
  );
}

class PrivacySettings {
  PrivacySettings({
    this.showPhoneToTeam = true,
    this.showActivityStatus = true,
    this.allowInviteByPhone = true,
  });

  final bool showPhoneToTeam;
  final bool showActivityStatus;
  final bool allowInviteByPhone;

  factory PrivacySettings.fromJson(Map<String, dynamic> json) =>
      PrivacySettings(
        showPhoneToTeam: json['showPhoneToTeam'] as bool? ?? true,
        showActivityStatus: json['showActivityStatus'] as bool? ?? true,
        allowInviteByPhone: json['allowInviteByPhone'] as bool? ?? true,
      );

  Map<String, dynamic> toJson() => {
    'showPhoneToTeam': showPhoneToTeam,
    'showActivityStatus': showActivityStatus,
    'allowInviteByPhone': allowInviteByPhone,
  };

  PrivacySettings copyWith({
    bool? showPhoneToTeam,
    bool? showActivityStatus,
    bool? allowInviteByPhone,
  }) => PrivacySettings(
    showPhoneToTeam: showPhoneToTeam ?? this.showPhoneToTeam,
    showActivityStatus: showActivityStatus ?? this.showActivityStatus,
    allowInviteByPhone: allowInviteByPhone ?? this.allowInviteByPhone,
  );
}

class UserSession {
  UserSession({
    required this.id,
    required this.deviceName,
    required this.deviceType,
    required this.browser,
    required this.os,
    required this.ip,
    required this.lastSeenAt,
    required this.createdAt,
    required this.expiresAt,
    required this.current,
  });

  final String id;
  final String deviceName;
  final String deviceType;
  final String browser;
  final String os;
  final String ip;
  final DateTime? lastSeenAt;
  final DateTime? createdAt;
  final DateTime? expiresAt;
  final bool current;

  factory UserSession.fromJson(Map<String, dynamic> json) => UserSession(
    id: json['id']?.toString() ?? '',
    deviceName: json['deviceName']?.toString() ?? 'دستگاه ناشناس',
    deviceType: json['deviceType']?.toString() ?? 'desktop',
    browser: json['browser']?.toString() ?? 'Browser',
    os: json['os']?.toString() ?? 'Unknown',
    ip: json['ip']?.toString() ?? '',
    lastSeenAt: DateTime.tryParse(json['lastSeenAt']?.toString() ?? ''),
    createdAt: DateTime.tryParse(json['createdAt']?.toString() ?? ''),
    expiresAt: DateTime.tryParse(json['expiresAt']?.toString() ?? ''),
    current: json['current'] as bool? ?? false,
  );
}

class Business {
  Business({
    required this.id,
    required this.name,
    required this.phones,
    this.address = '',
    this.licenseStatus = 'trial',
  });
  final String id;
  final String name;
  final List<String> phones;
  final String address;
  final String licenseStatus;

  factory Business.fromJson(Map<String, dynamic> json) => Business(
    id: json['id']?.toString() ?? '',
    name: json['name']?.toString() ?? '',
    phones: (json['phones'] as List? ?? const [])
        .map((e) => e.toString())
        .toList(),
    address: json['address']?.toString() ?? '',
    licenseStatus: json['licenseStatus']?.toString() ?? 'trial',
  );
}

class Member {
  Member({
    required this.id,
    required this.userPhone,
    required this.role,
    required this.status,
    required this.commissionPercent,
    this.userDisplayName = '',
  });
  final String id;
  final String userPhone;
  final String userDisplayName;
  final String role;
  final String status;
  final double commissionPercent;

  factory Member.fromJson(Map<String, dynamic> json) => Member(
    id: json['id']?.toString() ?? '',
    userPhone: json['userPhone']?.toString() ?? '',
    userDisplayName: json['userDisplayName']?.toString() ?? '',
    role: json['role']?.toString() ?? '',
    status: json['status']?.toString() ?? '',
    commissionPercent: (json['commissionPercent'] as num? ?? 0).toDouble(),
  );
}

class Invitation {
  Invitation({
    required this.id,
    required this.businessName,
    required this.inviteePhone,
    required this.role,
    required this.status,
    required this.commissionPercent,
  });
  final String id;
  final String businessName;
  final String inviteePhone;
  final String role;
  final String status;
  final double commissionPercent;

  factory Invitation.fromJson(Map<String, dynamic> json) => Invitation(
    id: json['id']?.toString() ?? '',
    businessName: json['businessName']?.toString() ?? '',
    inviteePhone: json['inviteePhone']?.toString() ?? '',
    role: json['role']?.toString() ?? '',
    status: json['status']?.toString() ?? '',
    commissionPercent: (json['commissionPercent'] as num? ?? 0).toDouble(),
  );
}

class AreaNode {
  AreaNode({
    required this.id,
    required this.name,
    required this.streets,
    this.cityId = '',
  });

  final String id;
  final String name;
  final String cityId;
  final List<StreetNode> streets;

  factory AreaNode.fromJson(Map<String, dynamic> json) => AreaNode(
    id: json['id']?.toString() ?? '',
    name: json['name']?.toString() ?? '',
    cityId: json['cityId']?.toString() ?? '',
    streets: (json['streets'] as List? ?? const [])
        .map((e) => StreetNode.fromJson(Map<String, dynamic>.from(e as Map)))
        .toList(),
  );
}

class StreetNode {
  StreetNode({
    required this.id,
    required this.areaId,
    required this.name,
    required this.neighborhoods,
    this.cityId = '',
  });

  final String id;
  final String cityId;
  final String areaId;
  final String name;
  final List<NeighborhoodNode> neighborhoods;

  factory StreetNode.fromJson(Map<String, dynamic> json) => StreetNode(
    id: json['id']?.toString() ?? '',
    cityId: json['cityId']?.toString() ?? '',
    areaId: json['areaId']?.toString() ?? '',
    name: json['name']?.toString() ?? '',
    neighborhoods: (json['neighborhoods'] as List? ?? const [])
        .map(
          (e) => NeighborhoodNode.fromJson(Map<String, dynamic>.from(e as Map)),
        )
        .toList(),
  );
}

class NeighborhoodNode {
  NeighborhoodNode({
    required this.id,
    required this.areaId,
    required this.streetId,
    required this.name,
    this.cityId = '',
  });

  final String id;
  final String cityId;
  final String areaId;
  final String streetId;
  final String name;

  factory NeighborhoodNode.fromJson(Map<String, dynamic> json) =>
      NeighborhoodNode(
        id: json['id']?.toString() ?? '',
        cityId: json['cityId']?.toString() ?? '',
        areaId: json['areaId']?.toString() ?? '',
        streetId: json['streetId']?.toString() ?? '',
        name: json['name']?.toString() ?? '',
      );
}

class PropertyFileModel {
  PropertyFileModel({
    required this.id,
    required this.type,
    required this.types,
    required this.title,
    required this.addresses,
    required this.media,
    this.houseInfo = const {},
    this.vaultIds = const [],
    this.vaultPlacements = const [],
    this.sharingHistory = const [],
    this.status = 'active',
    this.isPartnershipCopy = false,
    this.partnershipCommissionPercent = 0,
    this.businessCommissionPercent = 0,
    this.ownerCommissionPercent = 0,
    this.salePrice = 0,
    this.finalPrice = 0,
    this.depositPrice = 0,
    this.rentPrice = 0,
    this.description = '',
    this.internalDescription = '',
    this.createdAt = '',
    this.updatedAt = '',
  });

  final String id;
  final String type;
  final List<String> types;
  final String title;
  final String status;
  final bool isPartnershipCopy;
  final String description;
  final String internalDescription;
  final int salePrice;
  final int finalPrice;
  final int depositPrice;
  final int rentPrice;
  final Map<String, dynamic> houseInfo;
  final List<PropertyAddressModel> addresses;
  final List<PropertyMediaModel> media;
  final List<String> vaultIds;
  final List<PropertyVaultPlacementModel> vaultPlacements;
  final List<PropertySharingHistoryModel> sharingHistory;
  final double partnershipCommissionPercent;
  final double businessCommissionPercent;
  final double ownerCommissionPercent;
  final String createdAt;
  final String updatedAt;

  factory PropertyFileModel.fromJson(
    Map<String, dynamic> json,
  ) => PropertyFileModel(
    id: json['id']?.toString() ?? '',
    type: json['type']?.toString() ?? '',
    status: json['status']?.toString() ?? 'active',
    isPartnershipCopy: json['isPartnershipCopy'] == true,
    types: (json['types'] as List? ?? const [])
        .map((e) => e.toString())
        .toList(),
    title: json['title']?.toString() ?? '',
    description: json['description']?.toString() ?? '',
    internalDescription: json['internalDescription']?.toString() ?? '',
    salePrice: (json['salePrice'] as num? ?? 0).toInt(),
    finalPrice: (json['finalPrice'] as num? ?? 0).toInt(),
    depositPrice: (json['depositPrice'] as num? ?? 0).toInt(),
    rentPrice: (json['rentPrice'] as num? ?? 0).toInt(),
    businessCommissionPercent: (json['businessCommissionPercent'] as num? ?? 0)
        .toDouble(),
    ownerCommissionPercent: (json['ownerCommissionPercent'] as num? ?? 0)
        .toDouble(),
    partnershipCommissionPercent:
        (json['partnershipCommissionPercent'] as num? ?? 0).toDouble(),
    createdAt: json['createdAt']?.toString() ?? '',
    updatedAt: json['updatedAt']?.toString() ?? '',
    houseInfo: Map<String, dynamic>.from(json['houseInfo'] as Map? ?? const {}),
    addresses: (json['addresses'] as List? ?? const [])
        .map((e) => PropertyAddressModel.fromJson(Map<String, dynamic>.from(e)))
        .toList(),
    media: (json['media'] as List? ?? const [])
        .map((e) => PropertyMediaModel.fromJson(Map<String, dynamic>.from(e)))
        .toList(),
    vaultIds: (json['vaultIds'] as List? ?? const [])
        .map((e) => e.toString())
        .toList(),
    vaultPlacements: (json['vaultPlacements'] as List? ?? const [])
        .map(
          (e) => PropertyVaultPlacementModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    sharingHistory: (json['sharingHistory'] as List? ?? const [])
        .map(
          (e) => PropertySharingHistoryModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
  );
}

class PropertySharingHistoryModel {
  PropertySharingHistoryModel({
    required this.requestId,
    required this.userId,
    this.userName = '',
    this.userPhone = '',
    this.status = '',
    this.commissionPercent = 0,
    this.sharedCopyFileId = '',
    this.createdAt = '',
  });

  final String requestId;
  final String userId;
  final String userName;
  final String userPhone;
  final String status;
  final double commissionPercent;
  final String sharedCopyFileId;
  final String createdAt;

  factory PropertySharingHistoryModel.fromJson(Map<String, dynamic> json) =>
      PropertySharingHistoryModel(
        requestId: json['requestId']?.toString() ?? '',
        userId: json['userId']?.toString() ?? '',
        userName: json['userName']?.toString() ?? '',
        userPhone: json['userPhone']?.toString() ?? '',
        status: json['status']?.toString() ?? '',
        commissionPercent: (json['commissionPercent'] as num? ?? 0).toDouble(),
        sharedCopyFileId: json['sharedCopyFileId']?.toString() ?? '',
        createdAt: json['createdAt']?.toString() ?? '',
      );
}

class PropertyVaultPlacementModel {
  PropertyVaultPlacementModel({
    required this.vaultId,
    this.commissionPercent = 25,
  });

  final String vaultId;
  final double commissionPercent;

  factory PropertyVaultPlacementModel.fromJson(Map<String, dynamic> json) =>
      PropertyVaultPlacementModel(
        vaultId: json['vaultId']?.toString() ?? '',
        commissionPercent: (json['commissionPercent'] as num? ?? 25).toDouble(),
      );

  Map<String, dynamic> toJson() => {
    'vaultId': vaultId,
    'commissionPercent': commissionPercent,
  };
}

class PropertyShareRequestModel {
  PropertyShareRequestModel({
    required this.id,
    required this.propertyFileId,
    required this.propertyTitle,
    this.requesterName = '',
    this.requesterPhone = '',
    this.commissionPercent = 0,
    this.status = '',
    this.sharedCopyFileId = '',
  });

  final String id;
  final String propertyFileId;
  final String propertyTitle;
  final String requesterName;
  final String requesterPhone;
  final double commissionPercent;
  final String status;
  final String sharedCopyFileId;

  factory PropertyShareRequestModel.fromJson(Map<String, dynamic> json) =>
      PropertyShareRequestModel(
        id: json['id']?.toString() ?? '',
        propertyFileId: json['propertyFileId']?.toString() ?? '',
        propertyTitle: json['propertyTitle']?.toString() ?? '',
        requesterName: json['requesterName']?.toString() ?? '',
        requesterPhone: json['requesterPhone']?.toString() ?? '',
        commissionPercent: (json['commissionPercent'] as num? ?? 0).toDouble(),
        status: json['status']?.toString() ?? '',
        sharedCopyFileId: json['sharedCopyFileId']?.toString() ?? '',
      );
}

class PropertyOfferHistoryModel {
  PropertyOfferHistoryModel({
    required this.id,
    this.action = '',
    this.fromStatus = '',
    this.toStatus = '',
    this.note = '',
    this.createdAt = '',
  });

  final String id;
  final String action;
  final String fromStatus;
  final String toStatus;
  final String note;
  final String createdAt;

  factory PropertyOfferHistoryModel.fromJson(Map<String, dynamic> json) =>
      PropertyOfferHistoryModel(
        id: json['id']?.toString() ?? '',
        action: json['action']?.toString() ?? '',
        fromStatus: json['fromStatus']?.toString() ?? '',
        toStatus: json['toStatus']?.toString() ?? '',
        note: json['note']?.toString() ?? '',
        createdAt: json['createdAt']?.toString() ?? '',
      );
}

class PropertyOfferModel {
  PropertyOfferModel({
    required this.id,
    required this.propertyFileId,
    required this.propertyTitle,
    this.ownerUserId = '',
    this.ownerName = '',
    this.owner,
    this.requesterUserId = '',
    this.requesterName = '',
    this.contactName = '',
    this.requestTitle = '',
    this.commissionPercent = 25,
    this.score = 0,
    this.tier = '',
    this.status = '',
    this.chatChannelId = '',
    this.sharedCopyFileId = '',
    this.propertyFile,
    this.history = const [],
  });

  final String id;
  final String propertyFileId;
  final String propertyTitle;
  final String ownerUserId;
  final String ownerName;
  final AppUser? owner;
  final String requesterUserId;
  final String requesterName;
  final String contactName;
  final String requestTitle;
  final double commissionPercent;
  final int score;
  final String tier;
  final String status;
  final String chatChannelId;
  final String sharedCopyFileId;
  final PropertyFileModel? propertyFile;
  final List<PropertyOfferHistoryModel> history;

  factory PropertyOfferModel.fromJson(Map<String, dynamic> json) =>
      PropertyOfferModel(
        id: json['id']?.toString() ?? '',
        propertyFileId: json['propertyFileId']?.toString() ?? '',
        propertyTitle: json['propertyTitle']?.toString() ?? '',
        ownerUserId: json['ownerUserId']?.toString() ?? '',
        ownerName: json['ownerName']?.toString() ?? '',
        owner: json['owner'] is Map
            ? AppUser.fromJson(Map<String, dynamic>.from(json['owner'] as Map))
            : null,
        requesterUserId: json['requesterUserId']?.toString() ?? '',
        requesterName: json['requesterName']?.toString() ?? '',
        contactName: json['contactName']?.toString() ?? '',
        requestTitle: json['requestTitle']?.toString() ?? '',
        commissionPercent: (json['commissionPercent'] as num? ?? 25).toDouble(),
        score: (json['score'] as num? ?? 0).toInt(),
        tier: json['tier']?.toString() ?? '',
        status: json['status']?.toString() ?? '',
        chatChannelId: json['chatChannelId']?.toString() ?? '',
        sharedCopyFileId: json['sharedCopyFileId']?.toString() ?? '',
        propertyFile: json['propertyFile'] is Map
            ? PropertyFileModel.fromJson(
                Map<String, dynamic>.from(json['propertyFile'] as Map),
              )
            : null,
        history: (json['history'] as List? ?? const [])
            .map(
              (e) => PropertyOfferHistoryModel.fromJson(
                Map<String, dynamic>.from(e),
              ),
            )
            .toList(),
      );
}

class NotificationModel {
  NotificationModel({
    required this.id,
    required this.type,
    required this.title,
    required this.body,
    this.businessId = '',
    this.propertyId = '',
    this.requestId = '',
    this.readAt = '',
  });

  final String id;
  final String type;
  final String title;
  final String body;
  final String businessId;
  final String propertyId;
  final String requestId;
  final String readAt;

  factory NotificationModel.fromJson(Map<String, dynamic> json) =>
      NotificationModel(
        id: json['id']?.toString() ?? '',
        type: json['type']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        body: json['body']?.toString() ?? '',
        businessId: json['businessId']?.toString() ?? '',
        propertyId: json['propertyId']?.toString() ?? '',
        requestId: json['requestId']?.toString() ?? '',
        readAt: json['readAt']?.toString() ?? '',
      );
}

class ChannelVaultModel {
  ChannelVaultModel({
    required this.id,
    required this.channelId,
    this.ownerUserId = '',
    this.businessId = '',
    this.title = '',
    this.isMain = false,
  });

  final String id;
  final String channelId;
  final String ownerUserId;
  final String businessId;
  final String title;
  final bool isMain;

  bool get isBusinessVault => businessId.isNotEmpty;

  factory ChannelVaultModel.fromJson(Map<String, dynamic> json) =>
      ChannelVaultModel(
        id: json['id']?.toString() ?? '',
        channelId: json['channelId']?.toString() ?? '',
        ownerUserId: json['ownerUserId']?.toString() ?? '',
        businessId: json['businessId']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        isMain: json['isMain'] == true,
      );
}

class PropertyAddressModel {
  PropertyAddressModel({
    required this.areaId,
    required this.streetId,
    required this.neighborhoodId,
    this.areaName = '',
    this.streetName = '',
    this.neighborhoodName = '',
    this.manualExactAddress = '',
  });

  final String areaId;
  final String streetId;
  final String neighborhoodId;
  final String areaName;
  final String streetName;
  final String neighborhoodName;
  final String manualExactAddress;

  factory PropertyAddressModel.fromJson(Map<String, dynamic> json) =>
      PropertyAddressModel(
        areaId: json['areaId']?.toString() ?? '',
        streetId: json['streetId']?.toString() ?? '',
        neighborhoodId: json['neighborhoodId']?.toString() ?? '',
        areaName: json['areaName']?.toString() ?? '',
        streetName: json['streetName']?.toString() ?? '',
        neighborhoodName: json['neighborhoodName']?.toString() ?? '',
        manualExactAddress: json['manualExactAddress']?.toString() ?? '',
      );

  Map<String, dynamic> toJson() => {
    'areaId': areaId,
    'streetId': streetId,
    'neighborhoodId': neighborhoodId,
    'manualExactAddress': manualExactAddress,
  };
}

class PropertyMediaModel {
  PropertyMediaModel({
    required this.id,
    required this.kind,
    required this.url,
    required this.contentType,
    required this.size,
  });

  final String id;
  final String kind;
  final String url;
  final String contentType;
  final int size;

  factory PropertyMediaModel.fromJson(Map<String, dynamic> json) =>
      PropertyMediaModel(
        id: json['id']?.toString() ?? '',
        kind: json['kind']?.toString() ?? '',
        url: json['url']?.toString() ?? '',
        contentType: json['contentType']?.toString() ?? '',
        size: (json['size'] as num? ?? 0).toInt(),
      );
}

class AdminAccountModel {
  AdminAccountModel({
    required this.id,
    required this.userId,
    required this.roles,
    required this.permissions,
    required this.status,
  });

  final String id;
  final String userId;
  final List<String> roles;
  final List<String> permissions;
  final String status;

  factory AdminAccountModel.fromJson(Map<String, dynamic> json) =>
      AdminAccountModel(
        id: json['id']?.toString() ?? '',
        userId: json['userId']?.toString() ?? '',
        roles: (json['roles'] as List? ?? const [])
            .map((e) => e.toString())
            .toList(),
        permissions: (json['permissions'] as List? ?? const [])
            .map((e) => e.toString())
            .toList(),
        status: json['status']?.toString() ?? '',
      );
}

class CityModel {
  CityModel({required this.id, required this.name});

  final String id;
  final String name;

  factory CityModel.fromJson(Map<String, dynamic> json) => CityModel(
    id: json['id']?.toString() ?? '',
    name: json['name']?.toString() ?? '',
  );
}

class LocationSuggestionModel {
  LocationSuggestionModel({
    required this.id,
    required this.cityId,
    required this.type,
    required this.name,
    required this.status,
    this.manualParentName = '',
    this.reviewNote = '',
  });

  final String id;
  final String cityId;
  final String type;
  final String name;
  final String status;
  final String manualParentName;
  final String reviewNote;

  factory LocationSuggestionModel.fromJson(Map<String, dynamic> json) =>
      LocationSuggestionModel(
        id: json['id']?.toString() ?? '',
        cityId: json['cityId']?.toString() ?? '',
        type: json['type']?.toString() ?? '',
        name: json['name']?.toString() ?? '',
        status: json['status']?.toString() ?? '',
        manualParentName: json['manualParentName']?.toString() ?? '',
        reviewNote: json['reviewNote']?.toString() ?? '',
      );
}

class PlatformSettingsModel {
  PlatformSettingsModel({
    this.otpApiKeyMasked = '',
    this.serviceSmsApiKeyMasked = '',
  });

  final String otpApiKeyMasked;
  final String serviceSmsApiKeyMasked;

  factory PlatformSettingsModel.fromJson(Map<String, dynamic> json) =>
      PlatformSettingsModel(
        otpApiKeyMasked: json['otpApiKeyMasked']?.toString() ?? '',
        serviceSmsApiKeyMasked:
            json['serviceSmsApiKeyMasked']?.toString() ?? '',
      );
}

class PropertyMatchResultModel {
  PropertyMatchResultModel({
    required this.propertyFile,
    required this.score,
    required this.tier,
    required this.matchedReasons,
    required this.missedReasons,
    this.access = const [],
  });

  final PropertyFileModel propertyFile;
  final int score;
  final String tier;
  final List<String> matchedReasons;
  final List<String> missedReasons;
  final List<PropertyMatchAccessModel> access;

  factory PropertyMatchResultModel.fromJson(Map<String, dynamic> json) =>
      PropertyMatchResultModel(
        propertyFile: PropertyFileModel.fromJson(
          Map<String, dynamic>.from(json['propertyFile'] as Map? ?? const {}),
        ),
        score: (json['score'] as num? ?? 0).toInt(),
        tier: json['tier']?.toString() ?? '',
        matchedReasons: (json['matchedReasons'] as List? ?? const [])
            .map((e) => e.toString())
            .toList(),
        missedReasons: (json['missedReasons'] as List? ?? const [])
            .map((e) => e.toString())
            .toList(),
        access: (json['access'] as List? ?? const [])
            .map(
              (e) => PropertyMatchAccessModel.fromJson(
                Map<String, dynamic>.from(e),
              ),
            )
            .toList(),
      );
}

class PropertyMatchAccessModel {
  PropertyMatchAccessModel({
    this.source = '',
    this.vaultId = '',
    this.vaultTitle = '',
    this.commissionPercent = 0,
    this.collaboration = false,
  });

  final String source;
  final String vaultId;
  final String vaultTitle;
  final double commissionPercent;
  final bool collaboration;

  bool get isVault => source == 'vault';

  factory PropertyMatchAccessModel.fromJson(Map<String, dynamic> json) =>
      PropertyMatchAccessModel(
        source: json['source']?.toString() ?? '',
        vaultId: json['vaultId']?.toString() ?? '',
        vaultTitle: json['vaultTitle']?.toString() ?? '',
        commissionPercent:
            double.tryParse(json['commissionPercent']?.toString() ?? '') ?? 0,
        collaboration: json['collaboration'] == true,
      );
}

class ContactPhoneModel {
  ContactPhoneModel({this.label = 'موبایل', required this.value});

  final String label;
  final String value;

  factory ContactPhoneModel.fromJson(Map<String, dynamic> json) =>
      ContactPhoneModel(
        label: json['label']?.toString() ?? 'موبایل',
        value: json['value']?.toString() ?? '',
      );

  Map<String, dynamic> toJson() => {'label': label, 'value': value};
}

class ContactPropertyRefModel {
  ContactPropertyRefModel({
    this.propertyId = '',
    required this.title,
    this.note = '',
  });

  final String propertyId;
  final String title;
  final String note;

  factory ContactPropertyRefModel.fromJson(Map<String, dynamic> json) =>
      ContactPropertyRefModel(
        propertyId: json['propertyId']?.toString() ?? '',
        title: json['title']?.toString() ?? '',
        note: json['note']?.toString() ?? '',
      );

  Map<String, dynamic> toJson() => {
    'propertyId': propertyId,
    'title': title,
    'note': note,
  };
}

class ContactRequestLocationModel {
  ContactRequestLocationModel({
    this.level = 'area',
    this.name = '',
    this.areaId = '',
    this.streetId = '',
    this.includeAll = true,
    this.preference = 'preferred',
    this.description = '',
  });

  final String level;
  final String name;
  final String areaId;
  final String streetId;
  final bool includeAll;
  final String preference;
  final String description;

  factory ContactRequestLocationModel.fromJson(Map<String, dynamic> json) =>
      ContactRequestLocationModel(
        level: json['level']?.toString() ?? 'area',
        name: json['name']?.toString() ?? '',
        areaId: json['areaId']?.toString() ?? '',
        streetId: json['streetId']?.toString() ?? '',
        includeAll: json['includeAll'] as bool? ?? true,
        preference: json['preference']?.toString() ?? 'preferred',
        description: json['description']?.toString() ?? '',
      );

  Map<String, dynamic> toJson() => {
    'level': level,
    'name': name,
    'areaId': areaId,
    'streetId': streetId,
    'includeAll': includeAll,
    'preference': preference,
    'description': description,
  };
}

class ContactRequestFloorRuleModel {
  ContactRequestFloorRuleModel({
    this.floorMin = 0,
    this.floorMax = 0,
    this.elevator = false,
    this.preference = 'preferred',
  });

  final int floorMin;
  final int floorMax;
  final bool elevator;
  final String preference;

  factory ContactRequestFloorRuleModel.fromJson(Map<String, dynamic> json) =>
      ContactRequestFloorRuleModel(
        floorMin: (json['floorMin'] as num? ?? 0).toInt(),
        floorMax: (json['floorMax'] as num? ?? 0).toInt(),
        elevator: json['elevator'] as bool? ?? false,
        preference: json['preference']?.toString() ?? 'preferred',
      );

  Map<String, dynamic> toJson() => {
    'floorMin': floorMin,
    'floorMax': floorMax,
    'elevator': elevator,
    'preference': preference,
  };
}

class ContactRequestOptionFilterModel {
  ContactRequestOptionFilterModel({
    required this.key,
    required this.values,
    this.preference = 'preferred',
  });

  final String key;
  final List<String> values;
  final String preference;

  factory ContactRequestOptionFilterModel.fromJson(Map<String, dynamic> json) =>
      ContactRequestOptionFilterModel(
        key: json['key']?.toString() ?? '',
        values: (json['values'] as List? ?? const [])
            .map((e) => e.toString())
            .toList(),
        preference: json['preference']?.toString() ?? 'preferred',
      );

  Map<String, dynamic> toJson() => {
    'key': key,
    'values': values,
    'preference': preference,
  };
}

class ContactRequestBooleanFilterModel {
  ContactRequestBooleanFilterModel({
    required this.key,
    this.value = true,
    this.preference = 'preferred',
  });

  final String key;
  final bool value;
  final String preference;

  factory ContactRequestBooleanFilterModel.fromJson(
    Map<String, dynamic> json,
  ) => ContactRequestBooleanFilterModel(
    key: json['key']?.toString() ?? '',
    value: json['value'] as bool? ?? true,
    preference: json['preference']?.toString() ?? 'preferred',
  );

  Map<String, dynamic> toJson() => {
    'key': key,
    'value': value,
    'preference': preference,
  };
}

class ContactRequestNumberFilterModel {
  ContactRequestNumberFilterModel({
    required this.key,
    this.min = 0,
    this.preference = 'preferred',
  });

  final String key;
  final int min;
  final String preference;

  factory ContactRequestNumberFilterModel.fromJson(Map<String, dynamic> json) =>
      ContactRequestNumberFilterModel(
        key: json['key']?.toString() ?? '',
        min: (json['min'] as num? ?? 0).toInt(),
        preference: json['preference']?.toString() ?? 'preferred',
      );

  Map<String, dynamic> toJson() => {
    'key': key,
    'min': min,
    'preference': preference,
  };
}

class ContactRequestHistoryChangeModel {
  ContactRequestHistoryChangeModel({
    required this.field,
    this.from = '',
    this.to = '',
  });

  final String field;
  final String from;
  final String to;

  factory ContactRequestHistoryChangeModel.fromJson(
    Map<String, dynamic> json,
  ) => ContactRequestHistoryChangeModel(
    field: json['field']?.toString() ?? '',
    from: json['from']?.toString() ?? '',
    to: json['to']?.toString() ?? '',
  );
}

class ContactRequestHistoryEntryModel {
  ContactRequestHistoryEntryModel({
    required this.id,
    required this.description,
    required this.changedAt,
    required this.changes,
    this.changedById = '',
  });

  final String id;
  final String changedById;
  final DateTime? changedAt;
  final String description;
  final List<ContactRequestHistoryChangeModel> changes;

  factory ContactRequestHistoryEntryModel.fromJson(Map<String, dynamic> json) =>
      ContactRequestHistoryEntryModel(
        id: json['id']?.toString() ?? '',
        changedById: json['changedById']?.toString() ?? '',
        changedAt: DateTime.tryParse(json['changedAt']?.toString() ?? ''),
        description: json['description']?.toString() ?? '',
        changes: (json['changes'] as List? ?? const [])
            .map(
              (e) => ContactRequestHistoryChangeModel.fromJson(
                Map<String, dynamic>.from(e),
              ),
            )
            .toList(),
      );
}

class ContactRequestModel {
  ContactRequestModel({
    this.id = '',
    required this.title,
    this.type = '',
    this.status = 'active',
    this.budgetMin = 0,
    this.budgetMax = 0,
    this.purchaseMin = 0,
    this.purchaseMax = 0,
    this.suggestedPurchaseMin = 0,
    this.suggestedPurchaseMax = 0,
    this.partnershipMin = 0,
    this.partnershipMax = 0,
    this.shareMin = 0,
    this.shareMax = 0,
    this.depositMin = 0,
    this.depositMax = 0,
    this.suggestedDepositMin = 0,
    this.suggestedDepositMax = 0,
    this.rentMin = 0,
    this.rentMax = 0,
    this.suggestedRentMin = 0,
    this.suggestedRentMax = 0,
    this.minAreaM2 = 0,
    this.maxAgeYears = 0,
    this.convertible = false,
    this.maxConvertibleDeposit = 0,
    this.rentWithOwner = false,
    this.landMinAreaM2 = 0,
    this.buildingMinAreaM2 = 0,
    this.permitFloorsMin = 0,
    this.locations = const [],
    this.floorRules = const [],
    this.optionFilters = const [],
    this.booleanFilters = const [],
    this.numberFilters = const [],
    this.history = const [],
    this.changeDescription = '',
    this.note = '',
  });

  final String id;
  final String title;
  final String type;
  final String status;
  final int budgetMin;
  final int budgetMax;
  final int purchaseMin;
  final int purchaseMax;
  final int suggestedPurchaseMin;
  final int suggestedPurchaseMax;
  final int partnershipMin;
  final int partnershipMax;
  final int shareMin;
  final int shareMax;
  final int depositMin;
  final int depositMax;
  final int suggestedDepositMin;
  final int suggestedDepositMax;
  final int rentMin;
  final int rentMax;
  final int suggestedRentMin;
  final int suggestedRentMax;
  final int minAreaM2;
  final int maxAgeYears;
  final bool convertible;
  final int maxConvertibleDeposit;
  final bool rentWithOwner;
  final int landMinAreaM2;
  final int buildingMinAreaM2;
  final int permitFloorsMin;
  final List<ContactRequestLocationModel> locations;
  final List<ContactRequestFloorRuleModel> floorRules;
  final List<ContactRequestOptionFilterModel> optionFilters;
  final List<ContactRequestBooleanFilterModel> booleanFilters;
  final List<ContactRequestNumberFilterModel> numberFilters;
  final List<ContactRequestHistoryEntryModel> history;
  final String changeDescription;
  final String note;

  factory ContactRequestModel.fromJson(
    Map<String, dynamic> json,
  ) => ContactRequestModel(
    id: json['id']?.toString() ?? '',
    title: json['title']?.toString() ?? '',
    type: json['type']?.toString() ?? '',
    status: json['status']?.toString() ?? 'active',
    budgetMin: (json['budgetMin'] as num? ?? 0).toInt(),
    budgetMax: (json['budgetMax'] as num? ?? 0).toInt(),
    purchaseMin: (json['purchaseMin'] as num? ?? 0).toInt(),
    purchaseMax: (json['purchaseMax'] as num? ?? 0).toInt(),
    suggestedPurchaseMin: (json['suggestedPurchaseMin'] as num? ?? 0).toInt(),
    suggestedPurchaseMax: (json['suggestedPurchaseMax'] as num? ?? 0).toInt(),
    partnershipMin: (json['partnershipMin'] as num? ?? 0).toInt(),
    partnershipMax: (json['partnershipMax'] as num? ?? 0).toInt(),
    shareMin: (json['shareMin'] as num? ?? 0).toInt(),
    shareMax: (json['shareMax'] as num? ?? 0).toInt(),
    depositMin: (json['depositMin'] as num? ?? 0).toInt(),
    depositMax: (json['depositMax'] as num? ?? 0).toInt(),
    suggestedDepositMin: (json['suggestedDepositMin'] as num? ?? 0).toInt(),
    suggestedDepositMax: (json['suggestedDepositMax'] as num? ?? 0).toInt(),
    rentMin: (json['rentMin'] as num? ?? 0).toInt(),
    rentMax: (json['rentMax'] as num? ?? 0).toInt(),
    suggestedRentMin: (json['suggestedRentMin'] as num? ?? 0).toInt(),
    suggestedRentMax: (json['suggestedRentMax'] as num? ?? 0).toInt(),
    minAreaM2: (json['minAreaM2'] as num? ?? 0).toInt(),
    maxAgeYears: (json['maxAgeYears'] as num? ?? 0).toInt(),
    convertible: json['convertible'] as bool? ?? false,
    maxConvertibleDeposit: (json['maxConvertibleDeposit'] as num? ?? 0).toInt(),
    rentWithOwner: json['rentWithOwner'] as bool? ?? false,
    landMinAreaM2: (json['landMinAreaM2'] as num? ?? 0).toInt(),
    buildingMinAreaM2: (json['buildingMinAreaM2'] as num? ?? 0).toInt(),
    permitFloorsMin: (json['permitFloorsMin'] as num? ?? 0).toInt(),
    locations: (json['locations'] as List? ?? const [])
        .map(
          (e) => ContactRequestLocationModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    floorRules: (json['floorRules'] as List? ?? const [])
        .map(
          (e) => ContactRequestFloorRuleModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    optionFilters: (json['optionFilters'] as List? ?? const [])
        .map(
          (e) => ContactRequestOptionFilterModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    booleanFilters: (json['booleanFilters'] as List? ?? const [])
        .map(
          (e) => ContactRequestBooleanFilterModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    numberFilters: (json['numberFilters'] as List? ?? const [])
        .map(
          (e) => ContactRequestNumberFilterModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    history: (json['history'] as List? ?? const [])
        .map(
          (e) => ContactRequestHistoryEntryModel.fromJson(
            Map<String, dynamic>.from(e),
          ),
        )
        .toList(),
    changeDescription: json['changeDescription']?.toString() ?? '',
    note: json['note']?.toString() ?? '',
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'title': title,
    'type': type,
    'status': status,
    'budgetMin': budgetMin,
    'budgetMax': budgetMax,
    'purchaseMin': purchaseMin,
    'purchaseMax': purchaseMax,
    'suggestedPurchaseMin': suggestedPurchaseMin,
    'suggestedPurchaseMax': suggestedPurchaseMax,
    'partnershipMin': partnershipMin,
    'partnershipMax': partnershipMax,
    'shareMin': shareMin,
    'shareMax': shareMax,
    'depositMin': depositMin,
    'depositMax': depositMax,
    'suggestedDepositMin': suggestedDepositMin,
    'suggestedDepositMax': suggestedDepositMax,
    'rentMin': rentMin,
    'rentMax': rentMax,
    'suggestedRentMin': suggestedRentMin,
    'suggestedRentMax': suggestedRentMax,
    'minAreaM2': minAreaM2,
    'maxAgeYears': maxAgeYears,
    'convertible': convertible,
    'maxConvertibleDeposit': maxConvertibleDeposit,
    'rentWithOwner': rentWithOwner,
    'landMinAreaM2': landMinAreaM2,
    'buildingMinAreaM2': buildingMinAreaM2,
    'permitFloorsMin': permitFloorsMin,
    'locations': locations.map((e) => e.toJson()).toList(),
    'floorRules': floorRules.map((e) => e.toJson()).toList(),
    'optionFilters': optionFilters.map((e) => e.toJson()).toList(),
    'booleanFilters': booleanFilters.map((e) => e.toJson()).toList(),
    'numberFilters': numberFilters.map((e) => e.toJson()).toList(),
    if (changeDescription.isNotEmpty) 'changeDescription': changeDescription,
    'note': note,
  };
}

class ContactModel {
  ContactModel({
    required this.id,
    required this.displayName,
    required this.phones,
    required this.tags,
    required this.properties,
    required this.requests,
    this.firstName = '',
    this.lastName = '',
    this.company = '',
    this.note = '',
  });

  final String id;
  final String firstName;
  final String lastName;
  final String displayName;
  final String company;
  final List<ContactPhoneModel> phones;
  final List<String> tags;
  final List<ContactPropertyRefModel> properties;
  final List<ContactRequestModel> requests;
  final String note;

  factory ContactModel.fromJson(Map<String, dynamic> json) => ContactModel(
    id: json['id']?.toString() ?? '',
    firstName: json['firstName']?.toString() ?? '',
    lastName: json['lastName']?.toString() ?? '',
    displayName: json['displayName']?.toString() ?? '',
    company: json['company']?.toString() ?? '',
    phones: (json['phones'] as List? ?? const [])
        .map((e) => ContactPhoneModel.fromJson(Map<String, dynamic>.from(e)))
        .toList(),
    tags: (json['tags'] as List? ?? const []).map((e) => e.toString()).toList(),
    properties: (json['properties'] as List? ?? const [])
        .map(
          (e) => ContactPropertyRefModel.fromJson(Map<String, dynamic>.from(e)),
        )
        .toList(),
    requests: (json['requests'] as List? ?? const [])
        .map((e) => ContactRequestModel.fromJson(Map<String, dynamic>.from(e)))
        .toList(),
    note: json['note']?.toString() ?? '',
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'firstName': firstName,
    'lastName': lastName,
    'displayName': displayName,
    'company': company,
    'phones': phones.map((e) => e.toJson()).toList(),
    'tags': tags,
    'properties': properties.map((e) => e.toJson()).toList(),
    'requests': requests.map((e) => e.toJson()).toList(),
    'note': note,
  };
}
