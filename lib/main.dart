import 'dart:ffi';
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
import 'dart:isolate';

  
typedef goServerType = Void Function();
typedef goServerFunc = void Function();
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
  final lib = DynamicLibrary.open('${Directory(Platform.resolvedExecutable).parent.path}/Frameworks/goServer.so');

  final goServerFunc goServer =
      lib.lookupFunction<Void Function(), goServerFunc>('goServer');
  goServer();
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