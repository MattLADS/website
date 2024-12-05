import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/chatbot_page.dart';
import 'package:matt_lads_app/pages/assignments_page.dart';
import 'package:matt_lads_app/pages/select_user_page.dart';
import 'package:matt_lads_app/services/auth/auth_gate.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';
import 'package:matt_lads_app/pages/feed.dart';
import 'package:matt_lads_app/pages/profile.dart';
import 'package:matt_lads_app/pages/register_page.dart';
import 'package:matt_lads_app/pages/settings.dart';
import 'package:matt_lads_app/services/auth/backend_auth.dart';


void main() async {
  WidgetsFlutterBinding.ensureInitialized();
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
    final authService = BackendAuthService();
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      theme: Provider.of<ThemeProvider>(context).themeData,
      routes: {
        '/':(context) => AuthGate(),
        //'/forum': (context) => const HomePage(), 
        '/forum/': (context) => const HomePage(),
        '/profile': (context) {
          final args = ModalRoute.of(context)?.settings.arguments as Map<String, dynamic>?;

          return ProfilePage(
            url: args?['url'] ?? 'https://via.placeholder.com/150', // Default or passed profile picture URL
            username: args?['username'] ?? 'DefaultUser', // Default or passed username
            email: args?['email'] ?? 'No email provided', // Default or passed email
            classes: args?['classes'] ?? [], // Default or passed list of classes
          );
        },
        
        '/register': (context) => RegisterPage(
          onRegister: (username, password) async {
            try {
              await authService.register(username, password);
              Navigator.of(context).pushReplacementNamed('/forum/');
            } catch (e) {
              print(e); // Handle error (e.g., show a dialog)
            }
          },
          onLogin: () {
            Navigator.of(context).pushReplacementNamed('/login');
          },
        ),
        '/settings': (context) => const Settings(),
        '/chatbot': (context) => const ChatbotPage(),
        '/assignments': (context) => const AssignmentsPage(),
        '/direct_messages': (context) => SelectUserPage(),

      },
    );
  }
}