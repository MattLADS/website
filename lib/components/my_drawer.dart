import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer_tile.dart';
import 'package:matt_lads_app/pages/settings.dart';

// Drawer widget

class MyDrawer extends StatelessWidget {
  const MyDrawer({super.key});

  @override
  Widget build (BuildContext context) {
    return Drawer(
      backgroundColor: Theme.of(context).colorScheme.surface,
      child: SafeArea(
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
                builder: (context) => Settings(),
            )
            );

            },
          ), 
          //logout
          MyDrawerTile(
            title: "L O G O U T",
            icon: Icons.logout,
            onTap: () {},
          ), 
        ],
      )
    )
    );
  }
}