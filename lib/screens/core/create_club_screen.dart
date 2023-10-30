import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:flutter/material.dart';

class CreateClubScreen extends StatelessWidget {
  final TextEditingController _clubNameController = TextEditingController();

  CreateClubScreen({super.key});

  void createClub() async {}

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const MyAppBar(),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 25.0),
        child: Center(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text(
                'Create Club',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              TextField(
                decoration: const InputDecoration(
                  labelText: 'Club Name',
                ),
                controller: _clubNameController,
              ),
              const SizedBox(height: 15),
              MyButton(
                text: 'Create',
                onTap: createClub,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
