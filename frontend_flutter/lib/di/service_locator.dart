import "package:frontend_flutter/data/datasource/account_remote_datasource.dart";
import "package:frontend_flutter/data/datasource/account_remote_datasource_impl.dart";
import "package:frontend_flutter/data/repository/account_repository_impl.dart";
import "package:frontend_flutter/domain/repository/account_repo.dart";
import "package:get_it/get_it.dart";

final getIt = GetIt.instance;

void setupServiceLocator(){
  // reg AccountRemoteDatasource
  getIt.registerSingleton<AccountRemoteDatasource>(AccountRemoteDatasourceImpl());
  
  // reg AccountRepository
  getIt.registerSingleton<AccountRepo>(
    AccountRepositoryImpl(
      accountRemoteDatasource: getIt<AccountRemoteDatasource>()
    )
  );
}