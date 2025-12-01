import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend_flutter/domain/repository/auth_repository.dart';
import 'package:frontend_flutter/presentation/auth_state.dart';

class AuthCubit extends Cubit<AuthState> {
  final AuthRepository authRepository;

  AuthCubit(this.authRepository) : super(AuthInitial());

  Future<void> checkAuth() async {
    try {
      final auth = await authRepository.getStoredAuthToken();
      if (auth != null) {
        emit(AuthAuthenticated());
      } else {
        emit(AuthUnauthenticated());
      }
    } catch (e) {
      emit(AuthUnauthenticated());
    }
  }

  Future<void> login(String userId, String password) async {
    try {
      emit(AuthLoading());
      print('AuthCubit.login() called');
      final auth = await authRepository.login(userId, password);
      print('AuthCubit.login() success: ${auth.id}');
      emit(AuthAuthenticated());
      print('AuthCubit emitted AuthAuthenticated');
    } catch (e) {
      print('AuthCubit.login() error: $e');
      emit(AuthError(e.toString()));
    }
  }

  Future<void> register(String userId, String password) async {
    try {
      emit(AuthLoading());
      await authRepository.register(userId, password);
      emit(AuthAuthenticated());
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> logout() async {
    try {
      await authRepository.logout();
      emit(AuthUnauthenticated());
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }
}