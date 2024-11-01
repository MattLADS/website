//post model
//what every post has

import "package:cloud_firestore/cloud_firestore.dart";

class Post {
  final String id;
  final String uid;
  final String name;
  final String username;
  final String message;
  final Timestamp timestamp;
  final String course;
  final int likeCount;
  final int commentCount;

  Post({
    required this.id,
    required this.uid,
    required this.name,
    required this.username,
    required this.message,
    required this.timestamp,
    required this.course,
    required this.likeCount,
    required this.commentCount,
  });

  //this is where backend merge should be. 
  //this is where the data from the database is converted to a post object

  try {
    Post.fromMap(Map<String, dynamic> map)
      : id = map['id'],
        uid = map['uid'],
        name = map['name'],
        username = map['username'],
        message = map['message'],
        timestamp = map[Timestamp.now()],
        course = map['course'],
        likeCount = map[0],
        commentCount = map[0];
  } catch (e) {
    print(e);
  }

  //posts
  List<Post> _allPosts = [];

  List<Post> get allPosts => _allPosts;

  //posting a message. this is where more backend has to be

}