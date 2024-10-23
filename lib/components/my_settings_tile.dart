import 'package:flutter/material.dart';

class MySettingsTile extends StatelessWidget{
  final String title;
  final Widget action;

  const MySettingsTile({super.key, required this.title, required this.action});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text("Dark Mode"),
      trailing: CupertinoSwitch(
        onChanged: (value) => 
          Provider.of<ThemeProvider>(context, listen: false).toggleTheme(),
        value: Provider.of<ThemeProvider>(context, listen: false).isDarkMode,
      ),
    );
  } 
}