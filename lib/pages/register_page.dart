import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_button.dart';
import 'package:matt_lads_app/components/my_loading_circle.dart';
import 'package:matt_lads_app/components/my_text_field.dart';
import 'package:matt_lads_app/services/auth/auth_service.dart';
 
//1:08:53, TUTORIAL TIMESTAMP

class RegisterPage extends StatefulWidget {
  final void Function()? onTap;

  const RegisterPage({super.key, required this.onTap});
  
  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {

  //get access to auth service
  final _auth = AuthService();


  // BUILD UI

  final TextEditingController nameController = TextEditingController();
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController confirmPasswordController = TextEditingController();

  //register tapped
  void register() async {
    //password match create user
    if (passwordController.text == confirmPasswordController.text) {
      //loading circle
      showLoadingCircle(context);

      try {
        await _auth.register(emailController.text, passwordController.text);
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
    } else {
      //passwords don't match
      showDialog(
        context: context,
        builder: (context) => const AlertDialog(
          title: Text("Passwords don't match"),
          ),
    );
  }
    //don't match show error
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
          
                Text("Let's create an account", style: TextStyle(color: Theme.of(context).colorScheme.primary, 
                fontSize: 20)
                ),
          
                const SizedBox(height: 25),
          
                MyTextField(
                  controller: nameController,
                  hintText: "Enter your name",
                  obscureText: false,
                ),
          
                const SizedBox(height: 10),

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

                //sign in
                MyButton(text: "Register", onTap: register,),

                const SizedBox(height: 50),

                // alr a member? login 
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text ("Already a member?", style: TextStyle(color: Theme.of(context).colorScheme.primary),
                    ),
                    const SizedBox(width: 5),

                    GestureDetector(
                      //tap to go to login page
                      onTap: widget.onTap,
                      child: Text("Login here", style: TextStyle(color: Theme.of(context).colorScheme.primary, fontWeight: FontWeight.bold),
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

//email, password, confirm password, then redirected to home page after