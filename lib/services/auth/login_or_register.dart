import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/login_page.dart';
import 'package:matt_lads_app/pages/register_page.dart';
import 'package:matt_lads_app/services/auth/backend_auth.dart';
import 'dart:developer';
import 'package:shared_preferences/shared_preferences.dart';



class LoginOrRegister extends StatefulWidget {
  const LoginOrRegister({super.key});

  @override
  State<LoginOrRegister> createState() => _LoginOrRegisterState();
}


  class _LoginOrRegisterState extends State<LoginOrRegister> {
   bool showLoginPage = true;
   final BackendAuthService authService = BackendAuthService();
  
   void togglePages() {
    setState(() 
    {
      showLoginPage = !showLoginPage;
    });
   }
  
  Future<void> login(String username, String password) async {
    log('login function triggered with username: $username');
    try {
      log('Calling authService.login...');
      bool success = await authService.login(username, password);
      log('authService.login returned: $success');
      if (success) {
        SharedPreferences prefs = await SharedPreferences.getInstance();
        log('SharedPreferences instance created');
        await prefs.setString('username', username);
        log('Username saved in SharedPreferences: $username');
        Navigator.of(context).pushReplacementNamed('/forum/');
      }
    } catch (e) {
      log('Login failed with error: $e');
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Login Failed'),
          content: Text(e.toString()),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: Text('Try Again'),
            ),
          ],
        ),
      );
    }
  }
  Future<void> register(String username, String password) async {
    try {
      bool success = await authService.register(username, password);
      if (success) {
        SharedPreferences prefs = await SharedPreferences.getInstance();
        await prefs.setString('username', username);
        log('Username saved after registration: $username');
        Navigator.of(context).pushReplacementNamed('/forum/');
      }
    } catch (e) {
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Registration Failed'),
          content: Text(e.toString()),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: Text('Try Again'),
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return showLoginPage
        ? LoginPage(onRegister: togglePages, onLogin: login)
        : RegisterPage(onLogin: togglePages, onRegister: register);
  }

}