//handles fata to and from firebase
//methods for user profile, posting messages, likes + comments
//account specifics
//search users
//search posts and classes

//this is where I handled a lot of backend stuff with firebase. This is essentuially
//irrelevant if you're not using firebase, so this is just a reference 
//for now 

import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:matt_lads_app/models/user.dart';

class DatabaseService {
  //get instance of firestore database & auth
  final _db = FirebaseFirestore.instance;
  final _auth = FirebaseAuth.instance;

  //user profile, when a new user registers create an account and store details
  //n the database to display on their profile

  //save user info
  Future<void> saveUserInfoInFirebase({required String name, email}) async {
    //get user id
    String uid = _auth.currentUser!.uid;

    //extract username from email
    String username = email.split('@')[0];
    //if user signs up with test@gmail.com, username = test

    //create user profile
    UserProfile user = UserProfile(
      uid: uid,
      name: name,
      email: email,
      username: username,
      classes: [],
    );

    //convert user into map for firebase db storage 
    final userMap = user.toMap();

    //save user info in firebase
    await _db.collection("Users").doc(uid).set(userMap);
  }

  //get user info
  Future<UserProfile?> getUserInfoFromFirebase(String uid) async {

    try {
      DocumentSnapshot userDoc = await _db.collection("Users").doc(uid).get();
      return UserProfile.fromDocument(userDoc);
    } catch (e) {
      print(e);
      return null;
    }
  }
 
  //post message

  //likes

  //comments

  //account specifics 
}