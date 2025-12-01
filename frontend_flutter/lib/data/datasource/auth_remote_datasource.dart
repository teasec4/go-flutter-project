abstract class AuthRemoteDatasource {
  Future<String> login(String id, String password);
  Future<String> register(String id, String password);
}