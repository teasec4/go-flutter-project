part of 'account_cubit.dart';

@immutable
sealed class AccountState {}

final class AccountInitial extends AccountState {}

final class AccountLoadingData extends AccountState {}

final class AccountLoadedData extends AccountState {
  final Account account;
  AccountLoadedData(this.account);
}

final class AccountHadError extends AccountState {
  final String message;
  AccountHadError(this.message);
}
