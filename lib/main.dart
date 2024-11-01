import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/firebase_options.dart';
import 'package:matt_lads_app/services/auth/auth_gate.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';
import 'package:process_run/process_run.dart';
import 'dart:io';
import 'dart:developer';
import 'dart:convert';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  // Start the Go server
  await startGoServer();

  runApp(
    ChangeNotifierProvider(
      create: (context) => ThemeProvider(),
      child: const PostApp(),
    ),
  );
}

Future<void> startGoServer() async {
  // Adjust the path to your Go server executable
  const goServerPath = './website.go';
  final process = await Process.start(goServerPath, []);
  process.stdout.transform(utf8.decoder).listen((data) {
    log(data);
  });
  process.stderr.transform(utf8.decoder).listen((data) {
    log('Error: $data');
  });
}

class PostApp extends StatelessWidget {
  const PostApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      home: const AuthGate(),
      theme: Provider.of<ThemeProvider>(context).themeData,
    );
  }
}