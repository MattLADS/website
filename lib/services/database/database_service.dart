//handles fata to and from firebase
//methods for user profile, posting messages, likes + comments
//account specifics
//search users
//search posts and classes

import 'dart:developer';
import 'dart:convert';
import 'package:matt_lads_app/models/user.dart';
import 'package:http/http.dart' as http;

class DatabaseService {

  final String baseUrl = 'http://localhost:8080'; //backend URL

  //user profile, when a new user registers create an account and store details
  //n the database to display on their profile

  //save user info
  Future<void> saveUserInfoInFirebase({required String name, email}) async {

    //extract username from email
    String username = email.split('@')[0];
    //if user signs up with test@gmail.com, username = test

    //create user profile
    UserProfile user = UserProfile(
      uid: '', //this will be assigned in the backend
      name: name,
      email: email,
      username: username,
      classes: [],
    );
    //convert data to JSON and send to backend
    final response = await http.post(
      Uri.parse('$baseUrl/signup'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(user.toJson()),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to save user info');
    }
  }

  //get user info from backend
  Future<UserProfile?> getUserInfo(String uid) async {

    try {
      final response = await http.get(Uri.parse('$baseUrl/users/$uid'));
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        return UserProfile.fromJSON(data);
      } else {
        log('Failed to load user info: ${response.reasonPhrase}');
        return null;
      }

    } catch (e) {
      log('Error getting user info: $e');
      return null;
    }
  }
 
  //post message
  

  //likes

  //comments

  //account specifics 
}