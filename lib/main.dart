import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/login_page.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/pages/feed.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';

void main() {
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
      home: LoginPage(),
      //home: HomePage(),
      theme: Provider.of<ThemeProvider>(context). themeData,
    );
  }
}