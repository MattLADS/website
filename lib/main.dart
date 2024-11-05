import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/firebase_options.dart';
import 'package:matt_lads_app/services/auth/auth_gate.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';

// Conditional imports
// import 'package:matt_lads_app/go_server_stub.dart'
 //   if (dart.library.ffi) 'package:matt_lads_app/go_server_macos.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  // Start the Go server
  // startGoServer();

  runApp(
    ChangeNotifierProvider(
      create: (context) => ThemeProvider(),
      child: const PostApp(),
    ),
  );
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