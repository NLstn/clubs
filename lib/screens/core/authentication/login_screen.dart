import 'package:clubs/components/my_button.dart';
import 'package:clubs/components/my_clickable_text.dart';
import 'package:clubs/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class LoginScreen extends StatefulWidget {
  final void Function() toggleScreens;

  const LoginScreen({
    super.key,
    required this.toggleScreens,
  });

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  void login() {
    if (_emailController.text.isEmpty || _passwordController.text.isEmpty) {
      return;
    }
    final AuthService authService =
        Provider.of<AuthService>(context, listen: false);

    try {
      authService.loginWithEmailAndPassword(
        _emailController.text,
        _passwordController.text,
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Invalid email or password'),
        ),
      );
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 25.0),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const Text(
                'Login',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              TextField(
                decoration: const InputDecoration(
                  hintText: 'Email',
                ),
                controller: _emailController,
                onSubmitted: (value) => login(),
              ),
              TextField(
                decoration: const InputDecoration(
                  hintText: 'Password',
                ),
                obscureText: true,
                controller: _passwordController,
                onSubmitted: (value) => login(),
              ),
              const SizedBox(height: 15),
              MyButton(
                text: 'Login',
                onTap: login,
              ),
              const SizedBox(height: 25),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text('Don\'t have an account?'),
                  const SizedBox(width: 4),
                  MyClickableText(
                    text: 'Signup',
                    onTap: widget.toggleScreens,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
