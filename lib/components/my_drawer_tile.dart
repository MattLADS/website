import 'package:flutter/material.dart';

class MyDrawerTile extends StatelessWidget {
  final String title;
  final IconData icon;
  final void Function ()? onTap;
  
  const MyDrawerTile({
    super.key,
    required this.title,
    required this.icon,
    required this.onTap,
  });

  // BUILD UI
  @override
  Widget build (BuildContext context) {
  // List tile
    return ListTile(
      title: Text(
        title,
        style: TextStyle(color:Theme.of(context).colorScheme.inversePrimary),
      ), // Text
    leading: Icon(
      icon,
      color: Theme.of(context).colorScheme.primary,
      ),
    onTap: onTap,
    ); // ListTile
  }
}