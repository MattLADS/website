//authentication with backend. login, register, logout, delete account
import 'package:http/http.dart' as http;

class AuthService {
  static const String baseUrl = 'https://5ss3q5kd-8080.usw3.devtunnels.ms';

  // Login with backend
  Future<void> login(String username, String password) async {
    final url = Uri.parse('$baseUrl/login');
    final form = <String, dynamic>{};
    form['username'] = username;
    form['password'] = password;
    final response = await http.post(
      url,
      headers: {'Content-Type': 'multipart/form-data'},
      body: form,
    );

    if (response.statusCode == 200) {
      // Successful login
    } else if (response.statusCode == 401) {
      throw Exception('Invalid username or password');
    } else {
      throw Exception('Failed to login user');
    }
  }

  // Register with backend
  Future<void> register(String username, String password) async {
    final url = Uri.parse('$baseUrl/signup');
    final form = <String, dynamic>{};
    form['username'] = username;
    form['password'] = password;
    final response = await http.post(
      url,
      headers: {'Content-Type': 'multipart/form-data'},
      body: form,
    );

    if (response.statusCode == 200 || response.statusCode == 302) {
      // Successful registration
    } else if (response.statusCode == 409) {
      throw Exception('Username already exists');
    } else {
      throw Exception('Failed to register user');
    }
  }

  // Logout (if needed)
  Future<void> logout() async {
    // Implement logout functionality if needed
  }
}