//check if user is logged in or not, if yes go to home if not go to login/register

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/feed.dart';
import 'package:matt_lads_app/services/auth/login_or_register.dart';

class AuthGate extends StatelessWidget {
  const AuthGate({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: StreamBuilder(
        stream: FirebaseAuth.instance.authStateChanges(),
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return const HomePage();
          } else {
            return const LoginOrRegister();
          }
        },
      ),
    );
  }
}