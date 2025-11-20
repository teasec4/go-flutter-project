import 'package:frontend_flutter/data/datasource/account_remote_datasource.dart';
import 'package:frontend_flutter/domain/model/account.dart';
import 'package:frontend_flutter/domain/repository/account_repo.dart';

class AccountRepositoryImpl implements AccountRepo{
  // remote data source
  AccountRemoteDatasource accountRemoteDatasource;
  AccountRepositoryImpl({
      required this.accountRemoteDatasource
  });

  @override
  Future<Account> getAccountInfo(String id) async{
    try{
      final accountInfo = await accountRemoteDatasource.getAccountInfo(id);
      return accountInfo;
    } catch (e) {
      throw Exception(
        'Some error'
      );
    }
  }

  @override
  Future<Account> deposit(String id, int amount) async{
    try{
      final depositResponse = await accountRemoteDatasource.deposit(id, amount);
      return depositResponse;
    } catch (e){
      throw Exception("Deposit error");
    }
  }

  @override
  Future<Account> withdraw(String id, int amount) async{
    try{
      final withdrawResponse = await accountRemoteDatasource.withdraw(id, amount);
      return withdrawResponse;
    } catch (e){
      throw Exception("Withdraw error");
    }
  }
}