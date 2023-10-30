import 'package:flutter/material.dart';

class MyClickableText extends StatelessWidget {
  final String text;
  final void Function() onTap;
  const MyClickableText({super.key, required this.text, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: onTap,
        child: Text(
          text,
          style: const TextStyle(
            color: Colors.blue,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
    );
  }
}
