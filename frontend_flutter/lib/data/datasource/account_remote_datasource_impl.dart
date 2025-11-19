import 'dart:convert';
import 'dart:developer' as developer;

import 'package:frontend_flutter/data/datasource/account_remote_datasource.dart';
import 'package:frontend_flutter/domain/model/account.dart';
import 'package:http/http.dart' as http;

class AccountRemoteDatasourceImpl implements AccountRemoteDatasource{
  
  @override
  Future<Account> getAccountInfo(String id) async{
    try{
      final url = Uri.parse("http://192.168.5.10:8080/balance?id=$id");
      developer.log('Fetching from: $url');

      final response = await http
        .get(url)
      .timeout(Duration(seconds: 15));

      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200){
        final data = json.decode(response.body);
        developer.log('Parsed JSON: $data');
        return Account.fromJSON(data);
      } else {
        throw Exception(
          'HTTP response failed with status ${response.statusCode}'
        );
      }
    } catch (e) {
      developer.log('Error: $e');
      throw Exception('Failed to fetch: $e');
    }
  }
}