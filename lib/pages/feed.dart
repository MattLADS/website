import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer.dart';
import '../services/database/forum_service.dart';

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
    _topicsFuture = _forumService.fetchTopics();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      drawer: MyDrawer(),
      appBar: AppBar(
        title: Text("H O M E"),
        foregroundColor: Theme.of(context).colorScheme.primary,
      ),
      body: FutureBuilder<List<Map<String, dynamic>>>(
        future: _topicsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error loading topics'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return Center(child: Text('No topics available'));
          } else {
            List<Map<String, dynamic>> topics = snapshot.data!;
            return ListView.builder(
              itemCount: topics.length,
              itemBuilder: (context, index) {
                final topic = topics[index];
                return ListTile(
                  title: Text(topic['Title']),
                  subtitle: Text(topic['Content']),
                );
              },
            );
          }
        },
      ),
    );
  }
}