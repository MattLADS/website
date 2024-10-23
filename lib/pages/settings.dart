import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';
import 'package:provider/provider.dart';

class Settings extends StatelessWidget {
  const Settings({super.key});
  
  @override
  Widget build (BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface, 
      appBar: AppBar(
        title: const Text("S E T T I N G S"),
        foregroundColor: Theme.of(context).colorScheme.primary,
      ),
      body: Column(
        children: [
          ListTile(
            title: Text("Dark Mode"),
            trailing: CupertinoSwitch(
              onChanged: (value) => 
                Provider.of<ThemeProvider>(context, listen: false).toggleTheme(),
              value: Provider.of<ThemeProvider>(context, listen: false).isDarkMode,
            ),
          ),
      ],)

    );
  }
}