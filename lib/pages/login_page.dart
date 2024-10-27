
import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_loading_circle.dart';
import 'package:matt_lads_app/components/my_text_field.dart';
import 'package:matt_lads_app/components/my_button.dart';
import 'package:matt_lads_app/services/auth/auth_service.dart';

class LoginPage extends StatefulWidget {
  final void Function()? onTap;

  const LoginPage({super.key, required this.onTap});
  
  @override
  State<LoginPage> createState() => _LoginPageState();
}
class _LoginPageState extends State<LoginPage> {

  //access firebase auth service
  final _auth = AuthService();

  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();

  //login method
  void login() async {
    //loading circle lol
    showLoadingCircle(context);

    try {
      await _auth.loginEmailPassword(emailController.text, passwordController.text);
      //finish loading circle
      if (mounted) hideLoadingCircle(context);
    } catch (e) {
      if (mounted) hideLoadingCircle(context);
        if (mounted) {
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: Text(e.toString()),
            )
          );
        }
    }
  }

  @override
  Widget build (BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,

      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 25.0),
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const SizedBox(height: 50),
          
                Icon(
                  Icons.lock_open_rounded,
                  size: 72,
                  color: Theme.of(context).colorScheme.primary,
                ),
                const SizedBox(height: 50),
          
                Text("Welcome back!", style: TextStyle(color: Theme.of(context).colorScheme.primary, 
                fontSize: 20)
                ),
          
                const SizedBox(height: 25),
          
                MyTextField(
                  controller: emailController,
                  hintText: "Enter email or username",
                  obscureText: false,
                ),
          
                const SizedBox(height: 10),
          
                MyTextField(
                  controller: passwordController,
                  hintText: "Enter password",
                  obscureText: true,
                ),
                
                const SizedBox(height: 10),

                //forgot 
                Align(
                  alignment: Alignment.centerRight,
                  child: Text("Forgot password?",
                  style: TextStyle(color: Theme.of(context).colorScheme.primary, fontWeight: FontWeight.bold),
                  ),
                ),

                const SizedBox(height: 10),

                //sign in
                MyButton(text: "Login", onTap: login,),

                const SizedBox(height: 50),

                // not a member? register now
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text ("Not a member?", style: TextStyle(color: Theme.of(context).colorScheme.primary),
                    ),
                    const SizedBox(width: 5),

                    GestureDetector(
                      //tap to go to register page
                      onTap: widget.onTap,
                      child: Text("Register now", style: TextStyle(color: Theme.of(context).colorScheme.primary, fontWeight: FontWeight.bold),
                      ),  
                    ),
                  ],
                )
              ],
              ),
          ),
        )
      ),
    );
    }
}