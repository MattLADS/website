import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_text_field.dart';
import 'package:matt_lads_app/components/my_button.dart';

class RegisterPage extends StatefulWidget {
  final void Function(String username, String password) onRegister;
  final VoidCallback onLogin; // Callback for switching to the login page

  const RegisterPage({
    super.key,
    required this.onRegister,
    required this.onLogin,
  });

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmPasswordController = TextEditingController();

  void register() {
    final username = emailController.text;
    final password = passwordController.text;
    final confirmPassword = confirmPasswordController.text;

    if (password == confirmPassword) {
      widget.onRegister(username, password);
    } else {
      showDialog(
        context: context,
        builder: (context) => const AlertDialog(
          title: Text("Passwords don't match"),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
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
                Text("Let's create an account",
                    style: TextStyle(
                        color: Theme.of(context).colorScheme.primary,
                        fontSize: 20)),
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
                MyTextField(
                  controller: confirmPasswordController,
                  hintText: "Confirm password",
                  obscureText: true,
                ),
                const SizedBox(height: 25),
                MyButton(text: "Register", onTap: register),
                const SizedBox(height: 50),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text("Already a member?",
                        style:
                            TextStyle(color: Theme.of(context).colorScheme.primary)),
                    const SizedBox(width: 5),
                    GestureDetector(
                      onTap: widget.onLogin,
                      child: Text(
                        "Login here",
                        style: TextStyle(
                            color: Theme.of(context).colorScheme.primary,
                            fontWeight: FontWeight.bold),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
