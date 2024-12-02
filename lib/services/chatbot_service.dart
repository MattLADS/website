import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;

class ChatbotService {
  static const String baseUrl = 'http://localhost:8080';

  Future<String> sendMessage(String message) async {
    final url = Uri.parse('$baseUrl/chatbot');
    log("Sending message to chatbot: $message");

    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'message': message}),
    );

    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return data['response'];
    } else {
      throw Exception('Failed to get response from chatbot');
    }
  }
}