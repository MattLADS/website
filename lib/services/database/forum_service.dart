import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/browser_client.dart';
import 'package:http/browser_client.dart' as http;

final client = BrowserClient()..withCredentials = true;


class ForumService {
  static const String baseUrl = 'http://localhost:8080';

  Future<List<Map<String, dynamic>>> fetchTopics() async {
    log("fetchTopics called");
    
    final url = Uri.parse('$baseUrl/forum/');
    log("Attempting to fetch topics from: $url");  // Log start of fetch
      
    final response = await http.get(
      url,
      headers: {'Content-Type': 'application/json',},
    );


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
    log('POST URL: $url');

    final client = http.BrowserClient(); // Use BrowserClient

    final response = await client.post(
      url,
      headers: {
        'Content-Type': 'application/json',
      },
      body: json.encode({'title': title, 'content': content}),
    );
    log('POST request to $url with body: ${json.encode({'title': title, 'content': content})}');
    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 201) {
      log('Post submitted successfully');
      return; // Success case
    } else {
      log('Failed to post topic. Status: ${response.statusCode}, Body: ${response.body}');
      throw Exception('Failed to post topic');
    }
  }

  Future<void> postComment(String title, String comment) async {
    final url = Uri.parse('$baseUrl/new-comment/');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'title': title, 'comment': comment}),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to post comment');
    }
  }
}
