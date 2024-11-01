//post tile widget 

import 'package:flutter/material.dart';
import 'package:matt_lads_app/models/post.dart';

class MyPostTile extends StatefulWidget {
  final Post post;
  const MyPostTile({
    super.key,
    required this.post,
  });

  @override
  State<MyPostTile> createState() => _MyPostTileState();
}
class _MyPostTileState extends State<MyPostTile> {
  // BUILD UI
  @override
  Widget build (BuildContext context) {
    // Container
    return Container(
      margin: EdgeInsets.symmetric (horizontal: 25, vertical: 4),

      padding: const EdgeInsets.all(20),

      decoration: BoxDecoration(
        //color of the post
        color: Theme.of(context).colorScheme.secondary,
        borderRadius: BorderRadius.circular(8),
      ),

      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                Icons.person,
                color: Theme.of(context).colorScheme.primary,

                ), 
              Text(widget.post.name,
                style: TextStyle(
                  color: Theme.of(context).colorScheme.primary,
                  fontWeight: FontWeight.bold,
                ),
              ),

              const SizedBox(width: 4),

              //username handle
              Text('@'+widget.post.username,
                style: TextStyle(
                  color: Theme.of(context).colorScheme.primary,
                ),
              ),
            ]
          ),

          const SizedBox(height: 10),

          //message
          Text(widget.post.message,
            style: TextStyle(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ],
      ),
    );
  }
}