import "package:frontend_flutter/data/datasource/account_remote_datasource.dart";
import "package:frontend_flutter/data/datasource/account_remote_datasource_impl.dart";
import "package:frontend_flutter/data/datasource/auth_remote_datasource.dart";
import "package:frontend_flutter/data/datasource/auth_remote_datasource_impl.dart";
import "package:frontend_flutter/data/repository/account_repository_impl.dart";
import "package:frontend_flutter/data/repository/auth_repository_impl.dart";
import "package:frontend_flutter/data/service/token_service.dart";
import "package:frontend_flutter/domain/repository/account_repo.dart";
import "package:frontend_flutter/domain/repository/auth_repository.dart";
import "package:frontend_flutter/presentation/account_cubit.dart";
import "package:frontend_flutter/presentation/auth_cubit.dart";
import "package:get_it/get_it.dart";

final getIt = GetIt.instance;

void setupServiceLocator(){
  // TokenService
  getIt.registerSingleton<TokenService>(TokenService());
  
  // Auth
  getIt.registerSingleton<AuthRemoteDatasource>(
    AuthRemoteDatasourceImpl()
  );
  getIt.registerSingleton<AuthRepository>(
    AuthRepositoryImpl(
      remoteDatasource: getIt<AuthRemoteDatasource>(),
      tokenService: getIt<TokenService>(),
    )
  );
  
  // Account
  getIt.registerSingleton<AccountRemoteDatasource>(
    AccountRemoteDatasourceImpl(tokenService: getIt<TokenService>())
  );
  getIt.registerSingleton<AccountRepo>(
    AccountRepositoryImpl(
      accountRemoteDatasource: getIt<AccountRemoteDatasource>()
    )
  );
  
  // Cubits
  getIt.registerSingleton<AuthCubit>(
    AuthCubit(getIt<AuthRepository>())
  );
  getIt.registerSingleton<AccountCubit>(
    AccountCubit(getIt<AccountRepo>())
  );
}