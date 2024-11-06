import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:matt_lads_app/pages/login_page.dart';
import 'package:matt_lads_app/pages/register_page.dart';

class LoginOrRegister extends StatefulWidget {
  const LoginOrRegister({super.key});

  @override
  State<LoginOrRegister> createState() => _LoginOrRegisterState();
}

class _LoginOrRegisterState extends State<LoginOrRegister> {
  bool showLoginPage = true;

  void togglePages() {
    setState(() {
      showLoginPage = !showLoginPage;
    });
  }

  Future<void> login(String username, String password) async {
    final url = Uri.parse('http://localhost:8080/login');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'username': username, 'password': password}),
    );

    if (response.statusCode == 200) {
      // Handle successful login
    } else {
      // Handle login error
      throw Exception('Failed to login user');
    }
  }

  Future<void> register(String username, String password) async {
    final url = Uri.parse('http://localhost:8080/signup');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'username': username, 'password': password}),
    );

    if (response.statusCode == 200) {
      // Handle successful registration
    } else {
      // Handle registration error
      throw Exception('Failed to register user');
    }
  }
    @override
    Widget build(BuildContext context) {
      if (showLoginPage) {
        return LoginPage(
          onTap: togglePages,
        );
      } else {
        return RegisterPage(
          onTap: togglePages,
        );
      }
    }
  }