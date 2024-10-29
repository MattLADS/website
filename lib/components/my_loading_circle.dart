import 'package:flutter/material.dart';

//showing loading circle
void showLoadingCircle(BuildContext context) {
  showDialog(
    context: context,
    builder: (context) => AlertDialog(
      backgroundColor: Colors.transparent,
      elevation: 0,
      content: Center(
        child: CircularProgressIndicator(
          valueColor: AlwaysStoppedAnimation<Color>(Theme.of(context).colorScheme.primary),
        ),
      ),
    )
  );
}

void hideLoadingCircle(BuildContext context) {
  Navigator.pop(context);
}