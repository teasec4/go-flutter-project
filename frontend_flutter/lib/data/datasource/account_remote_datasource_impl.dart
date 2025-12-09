import 'dart:async';
import 'dart:convert';
import 'dart:developer' as developer;
import 'dart:io';

import 'package:frontend_flutter/data/datasource/account_remote_datasource.dart';
import 'package:frontend_flutter/data/service/token_service.dart';
import 'package:frontend_flutter/domain/model/account.dart';
import 'package:http/http.dart' as http;

const String _baseUrl = "http://192.168.5.10:8080";

class AccountRemoteDatasourceImpl implements AccountRemoteDatasource {
  final TokenService _tokenService;

  AccountRemoteDatasourceImpl({required TokenService tokenService})
    : _tokenService = tokenService;

  Future<Map<String, String>> _getAuthHeaders() async {
    final token = await _tokenService.getToken();
    if (token == null || token.isEmpty) {
      throw Exception('No auth token found');
    }
    
    // Extract userId from JWT token (payload is the second part)
    String userId = '';
    try {
      final parts = token.split('.');
      if (parts.length == 3) {
        final payload = parts[1];
        // Add padding if needed
        final normalized = payload.replaceAll('-', '+').replaceAll('_', '/');
        final paddingNeeded = (4 - (normalized.length % 4)) % 4;
        final padded = normalized + ('=' * paddingNeeded);
        
        final decoded = utf8.decode(base64Decode(padded));
        final Map<String, dynamic> decodedMap = jsonDecode(decoded);
        userId = decodedMap['sub'] as String? ?? '';
      }
    } catch (e) {
      developer.log('Failed to extract userId from token: $e');
    }
    
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
      if (userId.isNotEmpty) 'X-User-ID': userId,
    };
  }

  @override
  Future<Account> deposit(String id, int amount) async {
    try {
      final url = Uri.parse("$_baseUrl/account/deposit");
      final headers = await _getAuthHeaders();
      final response = await http.post(
          url,
          headers: headers,
          body: jsonEncode({
            'amount': amount
          })
      ).timeout(const Duration(seconds: 15));

      developer.log('Deposit request sent: accountId=$id, amount=$amount');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');
      if (response.statusCode == 200) {
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Deposit failed: ${errorData['error']}');
      }
    } on SocketException catch (e) {
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch (e) {
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('Deposit Error: $e');
      rethrow;
    }
  }

  @override
  Future<Account> withdraw(String id, int amount) async {
    try {
      final url = Uri.parse("$_baseUrl/account/withdraw");
      final headers = await _getAuthHeaders();
      final response = await http.post(
        url,
        headers: headers,
        body: jsonEncode({
          'amount': amount
        }),
      ).timeout(const Duration(seconds: 15));

      developer.log('Withdraw request sent: accountId=$id, amount=$amount');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200) {
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Withdraw failed: ${errorData['error']}');
      }
    } on SocketException catch (e) {
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch (e) {
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('Withdraw Error: $e');
      rethrow;
    }
  }


  @override
  Future<Account> getAccountInfo(String id) async {
    try {
      final url = Uri.parse("$_baseUrl/account");
      final headers = await _getAuthHeaders();

      final response = await http.get(
        url,
        headers: headers,
      ).timeout(const Duration(seconds: 15));

      developer.log('GetAccountInfo request sent');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200) {
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Failed to get account: ${errorData['error']}');
      }
    } on SocketException catch (e) {
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch (e) {
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('GetAccountInfo Error: $e');
      rethrow;
    }
  }
}