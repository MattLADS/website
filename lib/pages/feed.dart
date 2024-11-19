import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer.dart';
import '../services/database/forum_service.dart';
import 'dart:developer';


class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {

  final ForumService _forumService = ForumService();
  late Future<List<Map<String, dynamic>>> _topicsFuture;

  @override
  void initState() {
    super.initState();
    log('Initializing HomePage - Fetching topics...');
    _topicsFuture = _forumService.fetchTopics();
    log('_topicsFuture initialized.');
  }

  @override
  Widget build(BuildContext context) {
    print("HomePage build method called");
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      drawer: MyDrawer(),
      appBar: AppBar(
        title: const Text("H O M E"),
        foregroundColor: Theme.of(context).colorScheme.primary,
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          log("floatingActionButton pressed");
          newPostButtonInfo();
        },
        child: Icon(Icons.add),
        tooltip: 'Create New Post',
      ),

      
      body: FutureBuilder<List<Map<String, dynamic>>>(
        future: _topicsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            log('Error in FutureBuilder: ${snapshot.error}');
            return const Center(child: Text('Error loading topics'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            log('Topics retrieved: ${snapshot.data}');
            return const Center(child: Text('No topics available'));
          } else {
            List<Map<String, dynamic>> topics = snapshot.data!;
            return ListView.builder(
              itemCount: topics.length,
              itemBuilder: (context, index) {
                final topic = topics[index];
                final title = topic['title'] ?? 'No Title'; // Provide a default if null
                final content = topic['content'] ?? 'No Content'; // Provide a default if null
    
                return ListTile(
                  title: Text(title),
                  subtitle: Text(content),
                  onTap: () {
                    //eventually Navigate to a topic detail page
                    },
                );
              },
            );
          }
        },
      ),
    );
  }
  
  void newPostButtonInfo() {
    log('newPostButtonInfo called');
    showDialog(
    context: context,
    builder: (context) {
      final postTopic = TextEditingController();
      final postContent = TextEditingController();
      final ForumService forumService = ForumService();

      return AlertDialog(
        title: Text('Create New Post'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: postTopic,
              decoration: InputDecoration(labelText: 'Topic'),
            ),
            SizedBox(height: 10),
            TextField(
              controller: postContent,
              decoration: InputDecoration(labelText: 'Content'),
              maxLines: 3,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
            },
            child: Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              log('Submit button pressed');
              final title = postTopic.text;
              final content = postContent.text;

              if (title.isNotEmpty && content.isNotEmpty) {
                log('Calling postTopic with Title: $title, Content: $content');
                try {
                  await forumService.postTopic(title, content);
                  Navigator.of(context).pop(); // Close the pop up
                  setState(() {
                    _topicsFuture = forumService.fetchTopics(); // Refresh topics
                  });
                  log('Post submitted successfully and feed updated');
                } catch (e) {
                  log('Error submitting post: $e');
                  showDialog(
                    context: context,
                    builder: (context) => AlertDialog(
                      title: Text('Error'),
                      content: Text('Failed to submit post. Error: $e.'),
                      actions: [
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Text('OK'),
                        ),
                      ],
                    ),
                  );
                }
              }
            },
            child: Text('Submit'),
          ),
        ],
      );
    },
  );
  }
}