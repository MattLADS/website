import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/firebase_options.dart';
import 'package:matt_lads_app/services/auth/auth_gate.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';

void main() async {

  //FIREBASE SETUP HERE.

  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);

  runApp(
    ChangeNotifierProvider(
      create: (context) => ThemeProvider(),
      child: const PostApp(),
    )
  );
}

class PostApp extends StatelessWidget {
  const PostApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      home: const AuthGate(),
      //home: HomePage(),
      theme: Provider.of<ThemeProvider>(context). themeData,
    );
  }
}