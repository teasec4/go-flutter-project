import 'package:flutter/material.dart';
import 'package:frontend_flutter/presentation/home_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Go Flutter App',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.red.shade600),
      ),
      home: HomePage(),
      routes: {
        '/home': (context) => HomePage(),
      },
    );
  }
}
