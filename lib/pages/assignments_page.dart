import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:io';

class AssignmentsPage extends StatefulWidget {
  const AssignmentsPage({super.key});

  @override
  State<AssignmentsPage> createState() => _AssignmentsPageState();
}

class _AssignmentsPageState extends State<AssignmentsPage> {
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _filePathController = TextEditingController();
  final List<Map<String, dynamic>> _assignments = [];

  @override
  void initState() {
    super.initState();
    _fetchAssignments();
  }

  Future<void> _fetchAssignments() async {
    try {
      final response = await http.get(Uri.parse('localhost:8080/assignments'));
      if (response.statusCode == 200) {
        setState(() {
          _assignments.clear();
          _assignments.addAll(List<Map<String, dynamic>>.from(json.decode(response.body)));
        });
      }
    } catch (e) {
      print('Failed to fetch assignments: $e');
    }
  }


  Future<void> _uploadAssignment() async {
    final title = _titleController.text;
    final description = _descriptionController.text;
    final filePath = _filePathController.text;
    if (title.isEmpty || description.isEmpty || filePath.isEmpty) return;

    try {
      final response = await http.post(
        Uri.parse('http://localhost:8080/assignments'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'title': title,
          'description': description,
          'filePath': filePath,
        }),
      );

      if (response.statusCode == 200) {
        _fetchAssignments();
      } else {
        throw Exception('Failed to upload assignment');
      }
    } catch (e) {
      print('Failed to upload assignment: $e');
    }
  }

  Future<void> _editAssignment(int id, String title, String description, String filePath) async {
    try {
      final response = await http.post(
        Uri.parse('localhost:8080/edit-assignment'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'id': id,
          'title': title,
          'description': description,
          'filePath': filePath,
        }),
      );

      if (response.statusCode == 200) {
        _fetchAssignments();
      }
    } catch (e) {
      print('Failed to edit assignment: $e');
    }
  }


  Future<void> _downloadAssignment(String filePath) async {
    try {
      final response = await http.get(Uri.parse('localhost:8080/download-assignment?filePath=$filePath'));
      if (response.statusCode == 200) {
        final bytes = response.bodyBytes;
        final file = File(filePath);
        await file.writeAsBytes(bytes);
        print('File downloaded to $filePath');
      }
    } catch (e) {
      print('Failed to download assignment: $e');
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
          color: Theme.of(context).colorScheme.surface,
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
                  const SizedBox(height: 10),
                  FloatingActionButton(
                    onPressed: _uploadAssignment,
                    backgroundColor: Theme.of(context).colorScheme.primary,
                    child: const Icon(Icons.upload),
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