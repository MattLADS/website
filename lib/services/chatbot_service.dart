import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class ChatbotService {
  static const String baseUrl = 'http://localhost:8080';

  Future<String> sendMessage(String message) async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    final url = Uri.parse('$baseUrl/chatbot');
    var username = prefs.getString('username');
    log("Sending message to chatbot: $message from $username");

    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'message': message, 'username': prefs.getString('username')}),
    );

    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 200) {
      return response.body;
    } else {
      throw Exception('Failed to get response from chatbot');
    }
  }
}