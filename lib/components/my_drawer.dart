import 'dart:developer';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer_tile.dart';
import 'package:matt_lads_app/pages/settings.dart';
import 'package:http/http.dart' as http;

// Drawer widget

class MyDrawer extends StatelessWidget {
  MyDrawer({super.key});

  //logout auth w/firebase
  //final _auth = AuthService();

  void logout(BuildContext context) async {
    try {
      final response = await http.post(
        Uri.parse('http://localhost:8080/logout'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        log("Logged out successfully");
        
        // Navigate back to the login page after logout
        Navigator.of(context).pushReplacementNamed('/login'); 
      } else {
        log("Failed to log out: ${response.reasonPhrase}");
      }
    } catch (e) {
     log("Error during logout: $e");
    }
  }
  @override
  Widget build (BuildContext context) {
    return Drawer(
      backgroundColor: Theme.of(context).colorScheme.surface,
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 25.0),
        child: Column(
          children: [
          // app logo
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 50.0),
            child: Icon(
              Icons.person,
              size: 72,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
          
          Divider(
            indent: 25,
            endIndent: 25,
            color: Theme.of(context).colorScheme.secondary,
          ),
        
          const SizedBox(height: 10),

          //home 
          MyDrawerTile(
            title: "H O M E",
            icon: Icons.home,
            onTap: () {
              Navigator.pop(context);
            },
          ),
          //profile
          MyDrawerTile(
            title: "P R O F I L E",
            icon: Icons.person,
            onTap: () {}, 
          ), 
          //search list
          MyDrawerTile(
            title: "S E A R C H",
            icon: Icons.search,
            onTap: () {}, 
          ), 
          //settings
          MyDrawerTile(
            title: "S E T T I N G S",
            icon: Icons.settings,
            onTap: () {
              Navigator.pop(context);
            // go to settings page
              Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => const Settings(),
            )
            );

            },
          ), 

          const Spacer(),
          
          //logout
          MyDrawerTile(
            title: "L O G O U T",
            icon: Icons.logout,
            onTap: () => logout(context),
          ), 
        ],
      )
    )
    ));
  }
}