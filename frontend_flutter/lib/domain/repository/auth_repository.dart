import 'package:frontend_flutter/domain/model/auth.dart';

abstract class AuthRepository {
  Future<Auth> login(String userId, String password);
  Future<Auth> register(String userId, String password);
  Future<void> logout();
  Future<Auth?> getStoredAuthToken();
}