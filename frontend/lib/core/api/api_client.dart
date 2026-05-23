import 'package:dio/dio.dart';

import '../logging/app_logger.dart';
import '../storage/token_storage.dart';

class ApiClient {
  ApiClient(this._storage)
    : dio = Dio(
        BaseOptions(
          baseUrl: const String.fromEnvironment(
            'API_BASE_URL',
            defaultValue: 'http://127.0.0.1:8080/api/v1',
          ),
          connectTimeout: const Duration(seconds: 8),
          receiveTimeout: const Duration(seconds: 18),
          headers: {'Content-Type': 'application/json'},
        ),
      ) {
    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await _storage.accessToken;
          if (token != null && token.isNotEmpty) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          handler.next(options);
        },
        onError: (error, handler) async {
          AppLogger.error(
            'api',
            error,
            stackTrace: error.stackTrace,
            details: {
              'method': error.requestOptions.method,
              'url': error.requestOptions.uri.toString(),
              'statusCode': error.response?.statusCode,
              'response': error.response?.data,
            },
          );
          if (_shouldRetry(error)) {
            try {
              final response = await dio.fetch<dynamic>(
                error.requestOptions
                  ..extra['_retry'] = true
                  ..connectTimeout = const Duration(seconds: 8)
                  ..receiveTimeout = const Duration(seconds: 18),
              );
              handler.resolve(response);
              return;
            } catch (retryError, stackTrace) {
              AppLogger.error(
                'api-retry',
                retryError,
                stackTrace: stackTrace,
                details: {
                  'method': error.requestOptions.method,
                  'url': error.requestOptions.uri.toString(),
                },
              );
            }
          }
          handler.next(error);
        },
      ),
    );
  }

  final TokenStorage _storage;
  final Dio dio;

  bool _shouldRetry(DioException error) {
    if (error.requestOptions.extra['_retry'] == true) {
      return false;
    }
    return error.type == DioExceptionType.connectionTimeout ||
        error.type == DioExceptionType.connectionError ||
        error.type == DioExceptionType.receiveTimeout;
  }

  String message(Object error) {
    AppLogger.error('ui-catch', error);
    if (error is DioException) {
      final data = error.response?.data;
      if (data is Map && data['error'] is Map) {
        return data['error']['message']?.toString() ?? 'خطای ارتباط با سرور';
      }
      if (error.type == DioExceptionType.connectionTimeout) {
        return 'اتصال به سرور بیش از حد طول کشید. لطفا چند لحظه بعد دوباره تلاش کنید.';
      }
      if (error.type == DioExceptionType.receiveTimeout) {
        return 'پاسخ سرور بیش از حد طول کشید. لطفا دوباره تلاش کنید.';
      }
      if (error.type == DioExceptionType.connectionError) {
        return 'ارتباط با سرور برقرار نشد. مطمئن شوید backend اجراست.';
      }
      if (error.response?.statusCode != null) {
        return 'درخواست با خطای ${error.response!.statusCode} مواجه شد';
      }
      return 'خطای ارتباط با سرور رخ داد';
    }
    return error.toString().isEmpty
        ? 'خطای پیش‌بینی‌نشده رخ داد'
        : error.toString();
  }
}
