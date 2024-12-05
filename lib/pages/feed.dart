import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer.dart';
import '../services/database/forum_service.dart';
import 'dart:developer';

class HomePage extends StatefulWidget {
  @override
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
            return ListView.builder(
              itemCount: topics.length,
              itemBuilder: (context, index) {
                final topic = topics[index];
                final title = topic['Title'] ?? 'Untitled Post';
                final content = topic['Content'] ?? 'No Content';
                final username = topic['Username'] ?? 'Unknown User';

                return Card(
                  margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  elevation: 4,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: InkWell(
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
                                child: const Text('C L O S E'),
                              ),
                            ],
                          );
                        },
                      );
                    },
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            title,
                            style: Theme.of(context).textTheme.titleLarge,
                          ),
                          const SizedBox(height: 8),
                          Text(
                            "Author: @$username",
                            style: Theme.of(context).textTheme.titleSmall?.copyWith(color: Colors.grey),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            content,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: Theme.of(context).textTheme.bodyMedium,
                          ),
                        ],
                      ),
                    ),
                  ),
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
        tooltip: 'Create New Post',
        child: const Icon(Icons.add),
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
                  try {
                    await ForumService().postTopic(title, content);
                    Navigator.of(context).pop();
                    setState(() {
                      // Refresh feed after post
                      topicsFuture = ForumService().fetchTopics();
                    });
                  } catch (e) {
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
                      content: const Text('Title and content cannot be empty.'),
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
      log('Dialog dismissed');
    });
  }
}
