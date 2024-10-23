import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_text_field.dart';
import 'package:matt_lads_app/components/my_button.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  
  @override
  State<LoginPage> createState() => _LoginPageState();
}
class _LoginPageState extends State<LoginPage> {

  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();

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
                MyButton(text: "Sign in", onTap: () {}),

                const SizedBox(height: 50),

                // not a member? register now
                Row(
                  children: [
                    Text ("Not a member?", style: TextStyle(color: Theme.of(context).colorScheme.primary),
                    ),
                    const SizedBox(width: 5),

                    GestureDetector(
                      onTap: () {},
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