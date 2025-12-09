# JWT на Flutter клиенте - Гайд

## ✅ ЧТО НЕ МЕНЯЕТСЯ

Текущий код ПОЛНОСТЬЮ совместим с JWT. Никакие изменения не требуются:

```dart
// 1. Сохранение токена - работает с JWT так же
await tokenService.saveToken(token);  // JWT токен

// 2. Получение токена - работает с JWT так же
String? token = await tokenService.getToken();  // JWT токен

// 3. Отправка в Authorization - работает с JWT так же
headers: {
  'Authorization': 'Bearer $token'  // JWT или обычный токен - для фронта одинаково
}

// 4. Логика логина - не меняется
await login(userId, password);  // Получит JWT вместо старого токена
```

---

## 🎯 ЧТО МОЖНО ДОБАВИТЬ (ОПЦИОНАЛЬНО)

Хотя текущая реализация работает, есть несколько опциональных улучшений для большей надежности:

### Опция 1: Проверка истечения JWT перед запросом

Если вам нужно узнать, истек ли токен, можно распарсить JWT:

```dart
import 'dart:convert';

class JwtService {
  /// Проверить истек ли JWT токен
  static bool isTokenExpired(String token) {
    try {
      final parts = token.split('.');
      if (parts.length != 3) return true;

      // Декодируем payload (вторая часть)
      final payload = parts[1];
      final decoded = utf8.decode(base64Url.decode(payload.padRight(
        payload.length + (4 - payload.length % 4) % 4,
        '=',
      )));
      
      final json = jsonDecode(decoded);
      final exp = json['exp'] as int; // expiration timestamp
      
      // Если текущее время > exp, токен истек
      return DateTime.now().millisecondsSinceEpoch > exp * 1000;
    } catch (e) {
      print('Error checking token expiration: $e');
      return true; // Считаем истекшим если не можем распарсить
    }
  }

  /// Получить userID из JWT
  static String? getUserIdFromToken(String token) {
    try {
      final parts = token.split('.');
      if (parts.length != 3) return null;

      final payload = parts[1];
      final decoded = utf8.decode(base64Url.decode(payload.padRight(
        payload.length + (4 - payload.length % 4) % 4,
        '=',
      )));
      
      final json = jsonDecode(decoded);
      return json['userID'] as String?;
    } catch (e) {
      print('Error parsing JWT: $e');
      return null;
    }
  }
}
```

Использование:
```dart
// Проверить истек ли токен
if (JwtService.isTokenExpired(token)) {
  // Токен истек, нужно заново залогиниться
  await logout();
} else {
  // Токен еще действителен, можем использовать
  await makeApiRequest();
}

// Получить userID из токена
String? userId = JwtService.getUserIdFromToken(token);
print('Current user: $userId');
```

---

### Опция 2: Добавить Interceptor для автоматической переотправки при 401

Если сервер вернул 401, можно автоматически выйти и перенаправить на логин:

```dart
// В account_remote_datasource_impl.dart

Future<dynamic> _makeRequest(
  Future<http.Response> Function() request,
) async {
  try {
    final response = await request();
    
    // Если токен невалидный (401)
    if (response.statusCode == 401) {
      developer.log('Unauthorized - token expired or invalid');
      // Автоматически удалить токен
      final tokenService = TokenService();
      await tokenService.deleteToken();
      // Можно выбросить специальное исключение
      throw UnauthorizedException('Token expired, please login again');
    }
    
    return response;
  } catch (e) {
    rethrow;
  }
}
```

---

### Опция 3: Хранение токена более безопасно

Текущий код использует `SharedPreferences` (не зашифровано). Для более высокой безопасности:

```bash
flutter pub add flutter_secure_storage
```

```dart
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenService {
  static const String _tokenKey = 'auth_token';
  static late FlutterSecureStorage _secureStorage;
  String? _cachedToken;

  static Future<void> initialize() async {
    _secureStorage = const FlutterSecureStorage(
      aOptions: AndroidOptions(
        keyCipherAlgorithm: KeyCipherAlgorithm.RSA_ECB_OAEPwithSHA_256andMGF1Padding,
        storageCipherAlgorithm: StorageCipherAlgorithm.AES_GCM_NoPadding,
      ),
      iOptions: IOSOptions(
        accessibility: KeychainAccessibility.first_available_when_unlocked,
      ),
    );
  }

  Future<void> saveToken(String token) async {
    try {
      _cachedToken = token;
      await _secureStorage.write(key: _tokenKey, value: token);
    } catch (e) {
      print('TokenService.saveToken error: $e');
    }
  }

  Future<String?> getToken() async {
    try {
      if (_cachedToken != null) {
        return _cachedToken;
      }

      final token = await _secureStorage.read(key: _tokenKey);
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
      await _secureStorage.delete(key: _tokenKey);
    } catch (e) {
      print('TokenService.deleteToken error: $e');
    }
  }
}
```

---

## 📊 Сравнение: Клиент со Stateful vs Stateless токеном

| Аспект | Stateful Token | JWT (Stateless) |
|--------|---|---|
| **Сохранение** | Сохраняется как есть | Сохраняется как есть |
| **Отправка** | `Bearer token` header | `Bearer token` header |
| **Логика валидации** | На сервере (в БД) | На сервере (подпись) |
| **Клиентская логика** | **НЕ МЕНЯЕТСЯ** | **НЕ МЕНЯЕТСЯ** |
| **Опциональные улучшения** | Проверка в клиенте БЕЗ смысла | Можно проверить exp в клиенте |

---

## 🚀 ИТОГ ДЛЯ КЛИЕНТА

**Текущий код полностью работает с JWT без изменений!**

```dart
// Все это работает с JWT точно так же:

// 1. Сохранение
await tokenService.saveToken(jwtToken);

// 2. Получение
String? token = await tokenService.getToken();

// 3. Использование в запросах
headers: {'Authorization': 'Bearer $token'}

// 4. Удаление при логауте
await tokenService.deleteToken();
```

Никаких изменений кода не требуется! JWT просто более эффективен на сервере (0 DB calls вместо 1).

---

## 📝 Рекомендации

1. **Минимум:** Ничего не меняй. Текущий код работает идеально.

2. **Хорошо:** Добавь опцию 2 (автоматическая переотправка при 401).

3. **Лучше всего:** 
   - Добавь опцию 2 (обработка 401)
   - Используй `flutter_secure_storage` вместо `SharedPreferences` (опция 3)
   - Опционально добавь JWT парсинг для debug (опция 1)

4. **Production:**
   - ✅ Используй опции 2 + 3
   - ✅ Убедись что API на HTTPS
   - ✅ Установи разумное время жизни токена (24 часа - норма)
