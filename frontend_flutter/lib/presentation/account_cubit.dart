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
      final account = await accountRepo.getAccountInfo(id);
      emit(AccountLoadedData(account));
    } catch (e) {
      emit(AccountHadError(e.toString()));
    }
  }

  Future<void> deposit(String id, int amount) async {
    try {
      emit(AccountLoadingData());
      final account = await accountRepo.deposit(id, amount);
      emit(AccountLoadedData(account));
    } catch (e) {
      emit(AccountHadError(e.toString()));
    }
  }

  Future<void> withdraw(String id, int amount) async {
    try {
      emit(AccountLoadingData());
      final account = await accountRepo.withdraw(id, amount);
      emit(AccountLoadedData(account));
    } catch (e) {
      emit(AccountHadError(e.toString()));
    }
  }
}
