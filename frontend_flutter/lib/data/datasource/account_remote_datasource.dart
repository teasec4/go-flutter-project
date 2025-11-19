import 'package:frontend_flutter/domain/model/account.dart';

abstract class AccountRemoteDatasource{
  Future<Account> getAccountInfo(String id);
  Future<Account> deposit(String id, int amount);
}