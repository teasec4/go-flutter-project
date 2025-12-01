import 'package:bloc/bloc.dart';
import 'package:frontend_flutter/domain/model/account.dart';
import 'package:frontend_flutter/domain/repository/account_repo.dart';
import 'package:meta/meta.dart';

part 'account_state.dart';

class AccountCubit extends Cubit<AccountState> {
  final AccountRepo accountRepo;

  AccountCubit(this.accountRepo) : super(AccountInitial());

  Future<void> getAccountInfo(String id) async {
    try {
      emit(AccountLoadingData());
      print('AccountCubit.getAccountInfo() called');
      final account = await accountRepo.getAccountInfo(id);
      print('AccountCubit.getAccountInfo() success');
      emit(AccountLoadedData(account));
    } catch (e) {
      print('AccountCubit.getAccountInfo() error: $e');
      emit(AccountHadError(_extractErrorMessage(e)));
    }
  }

  Future<void> deposit(String id, int amount) async {
    try {
      emit(AccountLoadingData());
      final account = await accountRepo.deposit(id, amount);
      emit(AccountLoadedData(account));
    } catch (e) {
      emit(AccountHadError(_extractErrorMessage(e)));
    }
  }

  Future<void> withdraw(String id, int amount) async {
    try {
      emit(AccountLoadingData());
      final account = await accountRepo.withdraw(id, amount);
      emit(AccountLoadedData(account));
    } catch (e) {
      emit(AccountHadError(_extractErrorMessage(e)));
    }
  }

  String _extractErrorMessage(dynamic error) {
    if (error is Exception) {
      // Extract message from "Exception: message" format
      final str = error.toString();
      if (str.startsWith('Exception: ')) {
        return str.substring(11); // Remove "Exception: " prefix
      }
      return str;
    }
    return error.toString();
  }
}
