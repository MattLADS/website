import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/browser_client.dart';
import 'package:http/browser_client.dart' as http;


class ForumService {
  static const String baseUrl = 'http://localhost:8080';

  Future<List<Map<String, dynamic>>> fetchTopics() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    final cookies = prefs.getString('username');
    log("fetchTopics called");
    
    final url = Uri.parse('$baseUrl/forum/');
    log("Attempting to fetch topics from: $url");  // Log start of fetch

    final client = http.BrowserClient();
      
    final response = await client.get(url);

    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 200) {
      final List<dynamic> data = json.decode(response.body);
      return data.map((topic) => topic as Map<String, dynamic>).toList();
    } else {
      throw Exception('Failed to load topics');
    }
  
  }

  Future<void> postTopic(String title, String content) async {
    log('postTopic called with Title: $title, Content: $content');
    final url = Uri.parse('$baseUrl/new-topic/');

    // Initialize SharedPreferences
    SharedPreferences prefs = await SharedPreferences.getInstance();
    final username = prefs.getString('username') ?? '';

    final client = http.BrowserClient();
    final response = await client.post(
      url,
      headers: {
        'Content-Type': 'application/json',
      },
      body: json.encode({
        'title': title, 
        'content': content,
        'username': username,
      }),
    );
    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 201) {
      log('Post submitted successfully');
    } else {
      log('Failed to post topic. Status: ${response.statusCode}, Body: ${response.body}');
      throw Exception('Failed to post topic');
    }
  }

  Future<void> postComment(String title, String comment) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    final cookies = prefs.getString('username');
    final url = Uri.parse('$baseUrl/new-comment/');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json', 'Cookie': 'username=$cookies'},
      body: json.encode({'title': title, 'comment': comment}),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to post comment');
    }
  }
}
