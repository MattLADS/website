//user profile
//every user has this, to check their own profile?
//uid, name, email, classes in (like OS and senior teams in the listed)


class UserProfile {
  final String uid;
  final String name;
  final String email;
  final String username;
  final List<String> classes;

  UserProfile ({
    required this.uid, 
    required this.name, 
    required this.email, 
    required this.username,
    required this.classes
  });

  //takes backend data and formats it in JSON
  factory UserProfile.fromJSON(Map<String, dynamic> json) {
    return UserProfile(
       uid: json['uid'],
       name: json['name'],
       email: json['email'],
       username: json['username'],
       classes: json['classes'].cast<String>(),
    );
  }

  //convert user profile to JSON for storage
  Map<String, dynamic> toJson() {
    return {
      'uid': uid,
      'name': name,
      'email': email,
      'username': username,
      'classes': classes,
    };
  }
}