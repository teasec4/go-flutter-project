import 'dart:async';
import 'dart:convert';
import 'dart:developer' as developer;
import 'dart:io';

import 'package:frontend_flutter/data/datasource/account_remote_datasource.dart';
import 'package:frontend_flutter/domain/model/account.dart';
import 'package:http/http.dart' as http;

const String _baseUrl = "http://192.168.5.10:8080";

class AccountRemoteDatasourceImpl implements AccountRemoteDatasource{
  String? _token;

  // Login and get auth token
  Future<void> login(String accountId) async {
    try {
      final url = Uri.parse("$_baseUrl/login");
      developer.log('Logging in with accountId: $accountId');

      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'accountId': accountId}),
      ).timeout(const Duration(seconds: 15));

      developer.log('Login Status code: ${response.statusCode}');
      developer.log('Login Response body: ${response.body}');

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        _token = data['token'] as String?;

        if (_token == null || _token!.isEmpty) {
          throw Exception('Empty token received from server');
        }

        developer.log('Login successful, token: $_token');
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Login failed: ${errorData['error']}');
      }
    } on SocketException catch(e){
      developer.log('Network error: $e');
      throw Exception('No internet connection');
    } on TimeoutException catch(e){
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
    } catch (e) {
      developer.log('Login Error: $e');
      rethrow;
    }
  }

  Map<String, String> _getAuthHeaders() {
    if (_token == null || _token!.isEmpty) {
      throw Exception('Not authenticated. Call login first.');
    }
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $_token',
    };
  }
  
  @override
  Future<Account> getAccountInfo(String id) async{
    try{
      // First, login to get token
      await login(id);

      final url = Uri.parse("$_baseUrl/account?id=$id");
      
      final response = await http
        .get(url, headers: _getAuthHeaders())
      .timeout(Duration(seconds: 15));

      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200){
         final data = jsonDecode(response.body);
         return Account.fromJSON(data);
       } else {
         throw Exception(
           '${jsonDecode(response.body)['error']}'
         );
       }
      } on SocketException catch(e){
       developer.log('Network error: $e');
       throw Exception('No internet connection');
      } on TimeoutException catch(e){
       developer.log('Timeout: $e');
       throw Exception('Request timeout. Try again');
      } catch (e) {
       developer.log('Error: $e');
       throw Exception('Failed to fetch: $e');
      }
  }

  @override
  Future<Account> deposit(String id, int amount) async{
    try{
      final url = Uri.parse("$_baseUrl/account/deposit");
      final response = await http.post(
        url,
        headers: _getAuthHeaders(),
        body: jsonEncode({
          'accountId' : id,
          'amount' : amount
        })
      ).timeout(const Duration(seconds: 15));
      
      developer.log('Deposit request sent: accountId=$id, amount=$amount');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');
      if (response.statusCode == 200){
        return Account.fromJSON(jsonDecode(response.body));
      } else {
         final errorData = jsonDecode(response.body);
         throw Exception('Deposit failed: ${errorData['error']}');
       }

      } on SocketException catch(e){
       developer.log('Network error: $e');
       throw Exception('No internet connection');
      } on TimeoutException catch(e){
       developer.log('Timeout: $e');
       throw Exception('Request timeout. Try again');
      } catch (e){
       developer.log('Deposit Error: $e');
       rethrow;
      }
  }

  @override
  Future<Account> withdraw(String id, int amount) async{
    try{
      final url = Uri.parse("$_baseUrl/account/withdraw");
      final response = await http.post(
        url,
        headers: _getAuthHeaders(),
        body: jsonEncode({
          'accountId' : id,
          'amount' : amount
        }),
      ).timeout(const Duration(seconds: 15));
      
      developer.log('Withdraw request sent: accountId=$id, amount=$amount');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');
      
      if (response.statusCode == 200){
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        final errorData = jsonDecode(response.body);
        throw Exception('Withdraw failed: ${errorData['error']}');
      }
      } on SocketException catch(e){
      developer.log('Network error: $e');
      throw Exception('No internet connection');
      } on TimeoutException catch(e){
      developer.log('Timeout: $e');
      throw Exception('Request timeout. Try again');
      } catch (e){
      developer.log('Withdraw Error: $e');
      rethrow;
      }
      }
      }