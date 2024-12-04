import 'package:flutter/material.dart';
import 'package:matt_lads_app/pages/direct_messages_page.dart';
import 'package:provider/provider.dart';
import 'package:matt_lads_app/themes/theme_provider.dart';

class SelectUserPage extends StatelessWidget {
  SelectUserPage({super.key});

  final List<String> _users = [
    'User1',
    'User2',
    'User3',
    'User4',
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Provider.of<ThemeProvider>(context).themeData;

    return Scaffold(
      backgroundColor: theme.colorScheme.surface,
      appBar: AppBar(
        title: const Text("S E L E C T  U S E R"),
        foregroundColor: theme.colorScheme.primary,
        backgroundColor: theme.colorScheme.secondary,
      ),
      body: Container(
        decoration: BoxDecoration(
          color: theme.colorScheme.background,
        ),
        child: ListView.builder(
          itemCount: _users.length,
          itemBuilder: (context, index) {
            final user = _users[index];
            return ListTile(
              title: Text(
                user,
                style: TextStyle(
                  color: theme.colorScheme.onSurface,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => DirectMessagesPage(selectedUser: user),
                  ),
                );
              },
            );
          },
        ),
      ),
    );
  }
}