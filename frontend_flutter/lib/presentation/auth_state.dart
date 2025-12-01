sealed class AuthState {}

final class AuthInitial extends AuthState {}

final class AuthLoading extends AuthState {}

final class AuthAuthenticated extends AuthState {}

final class AuthUnauthenticated extends AuthState {}

final class AuthError extends AuthState {
  final String message;
  AuthError(this.message);
}