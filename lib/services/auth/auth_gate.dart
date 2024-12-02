//check if user is logged in or not, if yes go to home if not go to login/register

import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/feed.dart';
import 'package:matt_lads_app/services/auth/login_or_register.dart';
import 'package:http/http.dart' as http;

class AuthGate extends StatelessWidget {
  const AuthGate({super.key});


  //checking login status
    Future<bool> isLoggedIn() async {
    final url = Uri.parse('http://localhost:8080/auth/status');
    final response = await http.get(url);
     // Returns true if status is 200 (logged in)
    return response.statusCode == 200;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: FutureBuilder<bool>(
        future: isLoggedIn(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            //loading icon while waiting to login
            return const Center(child: CircularProgressIndicator());
          } else if(snapshot.hasData && snapshot.data == true){
            //successfully logged in!
            return const HomePage();
          }else{
            //not logged in
            return const LoginOrRegister();
          }
        },
      ),
    );
  }
}