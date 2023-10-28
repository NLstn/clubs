import 'package:clubs/screens/core/authentication/login_screen.dart';
import 'package:clubs/screens/core/authentication/signup_screen.dart';
import 'package:flutter/material.dart';

class LoginOrSignupScreen extends StatefulWidget {
  const LoginOrSignupScreen({super.key});

  @override
  State<LoginOrSignupScreen> createState() => _LoginOrSignupScreenState();
}

class _LoginOrSignupScreenState extends State<LoginOrSignupScreen> {
  bool showLoginScreen = true;

  void toggleScreens() {
    setState(() {
      showLoginScreen = !showLoginScreen;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (showLoginScreen) {
      return LoginScreen(toggleScreens: toggleScreens);
    } else {
      return SignupScreen(toggleScreens: toggleScreens);
    }
  }
}
