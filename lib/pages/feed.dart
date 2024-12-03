import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer.dart';
import '../services/database/forum_service.dart';
import 'dart:developer';


class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => HomePageState();
}

class HomePageState extends State<HomePage> {

  final ForumService forumService = ForumService();
  late Future<List<Map<String, dynamic>>> topicsFuture;

  @override
  void initState() {
    super.initState();
    log('Initializing HomePage - Fetching topics...');
    topicsFuture = forumService.fetchTopics();
    log('topicsFuture initialized.');
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
      body: FutureBuilder<List<Map<String, dynamic>>>(
        future: topicsFuture,
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
            log("FutureBuilder received data: ${snapshot.data}");
            List<Map<String, dynamic>> topics = snapshot.data!;
            //log("FutureBuilder received data: $topics");
            return ListView.builder(
              itemCount: topics.length,
              itemBuilder: (context, index) {
                final topic = topics[index];
                final title = topic['Title'] ?? 'Untitled Post'; // Provide a default if null
                final content = topic['Content'] ?? 'No Content'; // Provide a default if null
                final username = topic['Username'] ?? 'Unknown User'; // Provide a default if null

                return ListTile(
                  title: Text(title),
                  subtitle: Text("Author: @$username\n$content"),
                  onTap: () {
                    showDialog(
                      context: context,
                      builder: (context) {
                        return AlertDialog(
                          title: Text(title),
                          content: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text("Author: @$username"),
                              const SizedBox(height: 10),
                              Text(content),
                            ],
                          ),
                          actions: [
                            TextButton(
                              onPressed: () => Navigator.of(context).pop(),
                              child: const Text('Close'),
                            ),
                          ],
                        );
                      },
                    );
                  },
                );
              },
            );
          }
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          print("floating action button pressed");
          newPostButtonInfo();
        },
        child: Icon(Icons.add),
        tooltip: 'Create New Post',
      ),
    );
  }

  void newPostButtonInfo() {
    final postTopic = TextEditingController();
    final postContent = TextEditingController();

    log('newPostButtonInfo called');
    showDialog(
      context: context,
      builder: (context) {
        log('Dialog builder called');
        return AlertDialog(
          title: const Text('Create New Post'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: postTopic,
                decoration: const InputDecoration(labelText: 'Post Title'),
              ),
              const SizedBox(height: 10),
              TextField(
                controller: postContent,
                decoration: const InputDecoration(labelText: 'Post Content'),
                maxLines: 3,
              ),
            ],
          ),
          actions: [
            ElevatedButton(
              onPressed: () async {
                final title = postTopic.text;
                final content = postContent.text;
                
                if (title.isNotEmpty && content.isNotEmpty) {
                  try{
                    await ForumService().postTopic(title, content);
                    Navigator.of(context).pop();
                    setState(() {
                        // Refresh feed after post
                        topicsFuture = ForumService().fetchTopics();
                     });
                  } catch(e){
                    log('Error submitting post: $e');
                    showDialog(
                      context: context,
                      builder: (context) => AlertDialog(
                        title: const Text('Error'),
                        content: const Text('Failed to submit post'),
                        actions: [
                          TextButton(
                            onPressed: () => Navigator.of(context).pop(),
                            child: const Text('OK'),
                          ),
                        ],
                      ),
                    );
                  }
                  } else {
                    log('Title or content is empty');
                    showDialog(
                      context: context,
                      builder: (context) => AlertDialog(
                        title: const Text('Error'),
                        content: Text('Title and content cannot be empty.'),
                        actions: [
                          TextButton(
                            onPressed: () => Navigator.of(context).pop(),
                            child: const Text('OK'),
                          ),
                        ],
                      ),
                    );
                  }
                },
              child: const Text('Submit'),
            ),
          ],
        );
      },
    ).then((_) {
      log('Dialog dismissed'); // Log when dialog is closed
    });
  }
}