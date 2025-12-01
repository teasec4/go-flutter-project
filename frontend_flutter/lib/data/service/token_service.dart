import 'package:shared_preferences/shared_preferences.dart';

class TokenService {
  static const String _tokenKey = 'auth_token';
  static late SharedPreferences _prefs;
  String? _cachedToken;

  static Future<void> initialize() async {
    _prefs = await SharedPreferences.getInstance();
  }

  Future<void> saveToken(String token) async {
    try {
      _cachedToken = token;
      await _prefs.setString(_tokenKey, token);
    } catch (e) {
      print('TokenService.saveToken error: $e');
    }
  }

  Future<String?> getToken() async {
    try {
      if (_cachedToken != null) {
        return _cachedToken;
      }

      final token = _prefs.getString(_tokenKey);
      if (token != null) {
        _cachedToken = token;
      }
      return token;
    } catch (e) {
      print('TokenService.getToken error: $e');
      return null;
    }
  }

  Future<void> deleteToken() async {
    try {
      _cachedToken = null;
      await _prefs.remove(_tokenKey);
    } catch (e) {
      print('TokenService.deleteToken error: $e');
    }
  }
}
