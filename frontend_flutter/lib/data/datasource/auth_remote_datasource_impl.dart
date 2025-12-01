import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:developer' as developer;

import 'package:http/http.dart' as http;
import 'package:frontend_flutter/data/datasource/auth_remote_datasource.dart';

const String _baseUrl = "http://192.168.5.10:8080";

class LoginRequest {
  final String userId;
  final String password;

  LoginRequest({required this.userId, required this.password});

  Map<String, dynamic> toJson() => {
        'userId': userId,
        'password': password
      };
}

class RegisterRequest {
  final String userId;
  final String password;

  RegisterRequest({required this.userId, required this.password});

  Map<String, dynamic> toJson() => {
        'userId': userId,
        'password': password
      };
}


class AuthRemoteDatasourceImpl implements AuthRemoteDatasource{

  @override
  Future<String> login(String userId, String password) async {
    try {
      final url = Uri.parse("$_baseUrl/login");

      final request = LoginRequest(userId: userId, password: password);
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(request.toJson()),
      ).timeout(const Duration(seconds: 15));

      developer.log('Login request sent: userId=$userId');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final token = data['token'] as String;

        if (token.isEmpty) {
          throw Exception('Empty token received from server');
        }

        developer.log('Login successful, token received');
        return token;
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Login failed: ${errorData['error']}');
      }
    } on SocketException catch (e) {
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch (e) {
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('Login Error: $e');
      rethrow;
    }
  }

  @override
  Future<String> register(String userId, String password) async {
    try {
      final url = Uri.parse("$_baseUrl/register");

      final request = RegisterRequest(userId: userId, password: password);
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(request.toJson()),
      ).timeout(const Duration(seconds: 15));

      developer.log('Register request sent: userId=$userId');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 201) {
        // После регистрации залогиниться
        developer.log('Registration successful, logging in...');
        return await login(userId, password);
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Register failed: ${errorData['error']}');
      }
    } on SocketException catch (e) {
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch (e) {
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('Register Error: $e');
      rethrow;
    }
  }

}