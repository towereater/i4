import 'package:flutter/material.dart';
import 'package:visualizer/features/screens/home_page.dart';

void main() {
  runApp(const VisualizerApp());
}

class VisualizerApp extends StatelessWidget {
  const VisualizerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Visualizer App',
      theme: ThemeData(
        primarySwatch: Colors.indigo,
        useMaterial3: true,
        brightness: Brightness.dark,
        scaffoldBackgroundColor: const Color(0xFF1E1E2C),
        fontFamily: 'Inter',
      ),
      debugShowCheckedModeBanner: false,
      home: const HomePage(),
    );
  }
}
