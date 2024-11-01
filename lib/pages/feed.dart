import 'package:flutter/material.dart';
import 'package:matt_lads_app/components/my_drawer.dart';
import 'package:matt_lads_app/models/post.dart';

//HOMEPAGE, the first page the user sees and main page with posts

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  //based on the database provider. firebase/backend merge line here
  //listents for changes in database
  late final listeningProvider = Provider.of<DatabaseProvider>(context);
  //executes functions
  late final databaseProvider = Provider.of<DatabaseProvider>(context, listen : false);

  @override
  void initState() {
    super.initState();
    //get all posts
    loadAllposts();
  }

  Future<void> loadAllPosts() async {
    try {
      final allPosts = await _db.getAllPosts();
    } catch (e) {
      print(e);
    }
  }

  //text controllets 
  final _messageController = TextEditingController();

  //open post message box
  void _openPostMessageBox() {
    showDialog(context: context, builder: (context) => AlertDialog(
      title: Text("New Question"),
      )
    );
  }

  //If a user wants to post a messsage...
  Future<void> postMessage(String message) async {
    try {
      await databaseProvider.postMessage(message);
    } catch (e) {
      print(e);
    }
  }
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.surface,
      drawer: MyDrawer(),

      //App Bar
      appBar: AppBar(
        title: Text("H O M E"),
        foregroundColor: Theme.of(context).colorScheme.primary,
      ),

      floatingActionButton: FloatingActionButton(
        onPressed: _openPostMessageBox,
        child:const Icon(Icons.add),

        ),

        body: _buildPostList(listeningProvider.allPosts),
    );

    Widget _buildPostList(List<Post> posts) {
      return posts.isEmpty ? 
      
      const Center (
        child: Text("Nothing here..."),
      )

      //post list NOT EMPTY
      : ListView.builder(
        itemCount: posts.length,
        itemBuilder: (context, index) {
          final post = posts[index];

          return MyPostTile(post: post);

          

          return Container(
            child: Text(post.message),
            );
          },
        );
    }
  }
}