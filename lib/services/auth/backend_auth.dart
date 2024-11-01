import 'package:http/http.dart' as http;
import 'dart:convert';

class BackendAuthService {
  static const String baseUrl = 'http://localhost:8080';

  // Register with backend
  Future<void> register(String username, String password) async {
    final url = Uri.parse('$baseUrl/signup/');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'username': username, 'password': password}),
    );

    if (response.statusCode == 200 || response.statusCode == 302) {
      //Navigator.pushNamed(context, '/login');
    } else if (response.statusCode == 409) {
      throw Exception('Username already exists');
    } else {
      throw Exception('Failed to register user');
    }
  }

  // Login with backend
  Future<void> login(String username, String password) async {
    final url = Uri.parse('$baseUrl/');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'username': username, 'password': password}),
    );

    if (response.statusCode == 200) {
      // Successful login
    } else if (response.statusCode == 401) {
      throw Exception('Invalid username or password');
    } else {
      throw Exception('Failed to login user');
    }
  }
}
