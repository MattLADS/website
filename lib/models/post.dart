class Post {
final String id; // id of this post
final String uid; // uid of the poster
final String name; // name of the poster
final String username; // username of poster
final String message; // message of the post
final int likeCount; // like count of this post
final List<String> likedBy; // list of user IDs who liked this post
Post({
required this.id,
required this.uid,
required this.name,
required this.username,
required this.message,
required this. likeCount,
required this. likedBy,
});
}