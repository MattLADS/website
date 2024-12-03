import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class AssignmentsPage extends StatefulWidget {
  const AssignmentsPage({super.key});

  @override
  State<AssignmentsPage> createState() => _AssignmentsPageState();
}

class _AssignmentsPageState extends State<AssignmentsPage> {
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _filePathController = TextEditingController();
  final List<Map<String, String>> _assignments = [];

  Future<void> _uploadAssignment() async {
    final title = _titleController.text;
    final description = _descriptionController.text;
    final filePath = _filePathController.text;

    if (title.isEmpty || description.isEmpty || filePath.isEmpty) return;

    try {
      final response = await http.post(
        Uri.parse('http://localhost:8080/upload-assignment'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'title': title,
          'description': description,
          'filePath': filePath,
        }),
      );

      if (response.statusCode == 200) {
        setState(() {
          _assignments.add({
            'title': title,
            'description': description,
            'filePath': filePath,
          });
        });
        _titleController.clear();
        _descriptionController.clear();
        _filePathController.clear();
      } else {
        throw Exception('Failed to upload assignment');
      }
    } catch (e) {
      setState(() {
        _assignments.add({'title': 'Error', 'description': 'Failed to upload assignment', 'filePath': ''});
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      appBar: AppBar(
        title: const Text("A S S I G N M E N T S"),
        foregroundColor: Theme.of(context).colorScheme.primary,
        backgroundColor: Theme.of(context).colorScheme.surface,
      ),
      body: Container(
        decoration: BoxDecoration(
          color: Theme.of(context).colorScheme.background,
        ),
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Column(
                children: [
                  TextField(
                    controller: _titleController,
                    decoration: InputDecoration(
                      hintText: 'Title',
                      filled: true,
                      fillColor: Theme.of(context).colorScheme.surface,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(30),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    ),
                  ),
                  const SizedBox(height: 10),
                  TextField(
                    controller: _descriptionController,
                    decoration: InputDecoration(
                      hintText: 'Description',
                      filled: true,
                      fillColor: Theme.of(context).colorScheme.surface,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(30),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    ),
                  ),
                  const SizedBox(height: 10),
                  TextField(
                    controller: _filePathController,
                    decoration: InputDecoration(
                      hintText: 'File Path',
                      filled: true,
                      fillColor: Theme.of(context).colorScheme.surface,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(30),
                        borderSide: BorderSide.none,
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    ),
                  ),
                  const SizedBox(height: 10),
                  FloatingActionButton(
                    onPressed: _uploadAssignment,
                    child: const Icon(Icons.upload),
                    backgroundColor: Theme.of(context).colorScheme.primary,
                  ),
                ],
              ),
            ),
            Expanded(
              child: ListView.builder(
                itemCount: _assignments.length,
                itemBuilder: (context, index) {
                  final assignment = _assignments[index];
                  return ListTile(
                    title: Text(
                      assignment['title']!,
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.onPrimary,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    subtitle: Text(
                      assignment['description']!,
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.onPrimary,
                        fontSize: 14,
                      ),
                    ),
                    trailing: IconButton(
                      icon: const Icon(Icons.download),
                      onPressed: () {
                        // Handle file download
                      },
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}