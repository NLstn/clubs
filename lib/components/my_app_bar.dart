import 'package:clubs/screens/core/home_screen.dart';
import 'package:clubs/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class MyAppBar extends StatefulWidget implements PreferredSizeWidget {
  const MyAppBar({super.key});

  @override
  State<MyAppBar> createState() => _MyAppBarState();

  @override
  Size get preferredSize => const Size.fromHeight(50);
}

class _MyAppBarState extends State<MyAppBar> {
  void logOut() {
    AuthService authService = Provider.of<AuthService>(context, listen: false);
    authService.logout();
  }

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: InkWell(
        onTap: () {
          Navigator.pushAndRemoveUntil(
            context,
            MaterialPageRoute(builder: (context) => const HomeScreen()),
            (route) => false,
          );
        },
        child: const Text('Clubs'),
      ),
      actions: [
        IconButton(
          onPressed: logOut,
          icon: const Icon(Icons.logout),
        ),
      ],
    );
  }
}
