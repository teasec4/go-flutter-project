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
      
      final response = await http
        .get(url)
      .timeout(Duration(seconds: 15));

      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');

      if (response.statusCode == 200){
        final data = json.decode(response.body);
        return Account.fromJSON(data);
      } else {
        throw Exception(
          '${jsonDecode(response.body)['error']}'
        );
      }
    } catch (e) {
      developer.log('Error: $e');
      throw Exception('Failed to fetch: $e');
    }
  }

  @override
  Future<Account> deposit(String id, int amount) async{
    try{
      final url = Uri.parse("http://192.168.5.10:8080/deposit");
      final response = await http.post(
        url,
        headers:  {'Content-Type': 'application/json'},
        body: jsonEncode({
          'accountId' : id,
          'amount' : amount
        })
      );
      
      developer.log('Deposit request sent: accountId=$id, amount=$amount');
      developer.log('Status code: ${response.statusCode}');
      developer.log('Response body: ${response.body}');
      if (response.statusCode == 200){
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        throw Exception(
            '${jsonDecode(response.body)['error']}'
        );
      }

    } catch (e){
      throw Exception("Failed to deposit $e");
    }
  }

  @override
  Future<Account> withdraw(String id, int amount) async{
    try{
      final url = Uri.parse("http://192.168.5.10:8080/withdraw");
      final response = await http.post(
        url,
          headers:  {'Content-Type': 'application/json'},
          body: jsonEncode({
            'accountId' : id,
            'amount' : amount
          }),
      );
      if (response.statusCode == 200){
        return Account.fromJSON(jsonDecode(response.body));
      } else {
        throw Exception(
          "Some json err"
        );
      }
    } catch (e){
      throw Exception("Bed req ${e}");
    }
  }
}