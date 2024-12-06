import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer';

class BackendAuthService {
  static const String baseUrl = 'http://localhost:8080';

  // Register with backend
  Future<bool> register(String username, String password, String email) async {
    final url = Uri.parse('$baseUrl/signup/');
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'username': username, 'password': password, 'email': email}),
    );
    
    print('Registration status code: ${response.statusCode}');
    print('Registration response body: ${response.body}');

    if (response.statusCode == 201) {
      //Navigator.pushNamed(context, '/login');
      return true;
    } else if (response.statusCode == 409) {
      throw Exception('Username already exists');
    } else {
      throw Exception('Failed to register user');
    }
  }

  // Login with backend
  Future<bool> login(String username, String password) async {
    final url = Uri.parse('$baseUrl/');
    final response = await http.post(
      url,
      headers: {
        'Content-Type': 'application/json'
      },
      body: json.encode({'username': username, 'password': password}),
    );

    log('Login request sent to: $url');
    log('Request body: ${json.encode({'username': username, 'password': password})}');
    log('Response status: ${response.statusCode}');
    log('Response body: ${response.body}');

    if (response.statusCode == 200) {
      log('Login successful');
      return true;
    } else if (response.statusCode == 401) {
      throw Exception('Invalid username or password');
    } else {
      throw Exception('Failed to login user');
    }
  }
  Future<bool> logout() async {
    final url = Uri.parse('$baseUrl/logout');
    final response = await  http.post(url);

    if (response.statusCode == 200) {
      // Logout was successful, navigate back to login page
      //Navigator.of(context).pushReplacementNamed('/login');
      return true;
    } else {
      throw Exception('Failed to log out');
    }
  }
}
