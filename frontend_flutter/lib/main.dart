import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend_flutter/core/colors.dart';
import 'package:frontend_flutter/data/service/token_service.dart';
import 'package:frontend_flutter/di/service_locator.dart';
import 'package:frontend_flutter/presentation/auth_cubit.dart';
import 'package:frontend_flutter/presentation/auth_state.dart';
import 'package:frontend_flutter/presentation/home_page.dart';
import 'package:frontend_flutter/presentation/login_page.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await TokenService.initialize();
  setupServiceLocator();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Banking App',
      theme: ThemeData(
        useMaterial3: true,
        colorScheme: ColorScheme.fromSeed(
          seedColor: AppColors.primary,
          brightness: Brightness.light,
        ),
        scaffoldBackgroundColor: AppColors.background,
        appBarTheme: AppBarTheme(
          backgroundColor: AppColors.surface,
          elevation: 0,
          titleTextStyle: const TextStyle(
            color: AppColors.textPrimary,
            fontSize: 18,
            fontWeight: FontWeight.w600,
          ),
          iconTheme: const IconThemeData(
            color: AppColors.textPrimary,
          ),
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.primary,
            foregroundColor: AppColors.surface,
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
        outlinedButtonTheme: OutlinedButtonThemeData(
          style: OutlinedButton.styleFrom(
            foregroundColor: AppColors.primary,
            side: const BorderSide(color: AppColors.border),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
          ),
        ),
        inputDecorationTheme: InputDecorationTheme(
          filled: true,
          fillColor: AppColors.surface,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: AppColors.border),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: AppColors.border),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
            borderSide: const BorderSide(color: AppColors.primary, width: 2),
          ),
          prefixIconColor: AppColors.textSecondary,
          labelStyle: const TextStyle(color: AppColors.textSecondary),
        ),
        snackBarTheme: SnackBarThemeData(
          backgroundColor: AppColors.textPrimary,
          contentTextStyle: const TextStyle(
            color: AppColors.surface,
            fontSize: 14,
          ),
          behavior: SnackBarBehavior.floating,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        dividerTheme: const DividerThemeData(
          color: AppColors.divider,
        ),
      ),
      home: BlocProvider(
        create: (_) => getIt<AuthCubit>()..checkAuth(),
        child: BlocListener<AuthCubit, AuthState>(
          listener: (context, state) {
            if (state is AuthError) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(state.message),
                  backgroundColor: Colors.red,
                ),
              );
            }
          },
          child: BlocBuilder<AuthCubit, AuthState>(
            builder: (context, state) {
              return switch (state) {
                AuthInitial() || AuthLoading() => const Scaffold(
                    body: Center(child: CircularProgressIndicator()),
                  ),
                AuthAuthenticated() => const HomePage(),
                _ => const LoginPage(), // AuthError и AuthUnauthenticated
              };
            },
          ),
        ),
      ),
    );
  }
}
