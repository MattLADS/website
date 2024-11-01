//user profile
//every user has this, to check their own profile?
//uid, name, email, classes in (like OS and senior teams in the listed)

import 'package:cloud_firestore/cloud_firestore.dart';

class UserProfile {
  final String uid;
  final String name;
  final String email;
  final String username;
  final List<String> classes;

  UserProfile ({
    required this.uid, 
    required this.name, 
    required this.email, 
    required this.username,
    required this.classes
  });

  //app to firebase
  //converting firestore document to user profile

  factory UserProfile.fromDocument(DocumentSnapshot doc) {
    return UserProfile(
       uid: doc['uid'],
       name: doc['name'],
       email: doc['email'],
       username: doc['username'],
       classes: doc['classes'].cast<String>(),
    );
  }
  //firebase to app
  //convert user profile to a map for storage

  Map<String, dynamic> toMap() {
    return {
      'uid': uid,
      'name': name,
      'email': email,
      'username': username,
      'classes': classes,
    };
  }
}