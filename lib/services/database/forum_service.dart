import 'dart:convert';
import 'package:http/http.dart' as http;

class ForumService {
  final String baseUrl = 'http://localhost:8080';

  // Fetch topics from the forum route
  Future<List<Map<String, dynamic>>> fetchTopics() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/forum/'));

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
