import 'package:frontend_flutter/data/datasource/auth_remote_datasource.dart';
import 'package:frontend_flutter/data/service/token_service.dart';
import 'package:frontend_flutter/domain/model/auth.dart';
import 'package:frontend_flutter/domain/repository/auth_repository.dart';

class AuthRepositoryImpl implements AuthRepository{
  final AuthRemoteDatasource _remoteDatasource;
  final TokenService _tokenService;

  AuthRepositoryImpl({
    required AuthRemoteDatasource remoteDatasource,
    required TokenService tokenService}) : _remoteDatasource = remoteDatasource,
                                          _tokenService = tokenService;

  @override
  Future<Auth?> getStoredAuthToken() async {
    final token = await _tokenService.getToken();
    if (token != null) {
      // Извлечь userId из токена (format: userId:uuid)
      final parts = token.split(':');
      if (parts.isNotEmpty) {
        return Auth(id: parts[0], token: token);
      }
    }
    return null;
  }

  @override
  Future<Auth> login(String userId, String password) async{
    final token = await _remoteDatasource.login(userId, password);
    final auth = Auth(id: userId, token: token);
    await _tokenService.saveToken(token);
    return auth;
  }

  @override
  Future<void> logout() async{
    await _tokenService.deleteToken();
  }

  @override
  Future<Auth> register(String userId, String password) async {
    final token = await _remoteDatasource.register(userId, password);
    final auth = Auth(id: userId, token: token);
    await _tokenService.saveToken(token);
    return auth;
  }

}