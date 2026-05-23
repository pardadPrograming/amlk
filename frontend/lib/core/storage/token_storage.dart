import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenStorage {
  static const _access = 'access_token';
  static const _refresh = 'refresh_token';
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  Future<String?> get accessToken => _storage.read(key: _access);
  Future<String?> get refreshToken => _storage.read(key: _refresh);

  Future<void> save(String accessToken, String refreshToken) async {
    await _storage.write(key: _access, value: accessToken);
    await _storage.write(key: _refresh, value: refreshToken);
  }

  Future<void> clear() async {
    await _storage.delete(key: _access);
    await _storage.delete(key: _refresh);
  }
}
