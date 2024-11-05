import 'dart:convert';
import 'package:http/http.dart' as http;

class ForumService {
  final String baseUrl = 'http://localhost:8080';

  // Fetch topics from the forum route
  Future<List<Map<String, dynamic>>> fetchTopics() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/forum/'));

      //this is an example of the HTTPs request code for the forum. Commented out for now.
    // try {
    //   final url = Uri.parse('$baseUrl/forum/'); 
    //   final response = await http.post(
    //     url,
    //     headers: {'Content-Type': 'application/json'},
    //     body: json.encode({
    //       'username': 'yourUsername',
    //       'password': 'yourPassword',
    //       'title': title,
    //       'comments': comments,
    //       'content': content,
    //     }),
    //   );

      if (response.statusCode == 200) {
        // Parse the JSON data into a list of maps
        List<dynamic> data = json.decode(response.body);
        return List<Map<String, dynamic>>.from(data);
      } else {
        throw Exception('Failed to load topics: ${response.reasonPhrase}');
      }
    } catch (e) {
      throw Exception('Failed to load topics: $e');
    }
  }
}
