import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/services/club_service.dart';
import 'package:clubs/services/news_service.dart';
import 'package:flutter/material.dart';

class CreateNewsScreen extends StatefulWidget {
  const CreateNewsScreen({super.key});

  @override
  State<CreateNewsScreen> createState() => _CreateNewsScreenState();
}

class _CreateNewsScreenState extends State<CreateNewsScreen> {
  final TextEditingController _newsTitleController = TextEditingController();
  final TextEditingController _newsContentController = TextEditingController();

  List<DropdownMenuItem> clubs = [];
  String? selectedClubId;

  @override
  void initState() {
    super.initState();
    loadClubs();
  }

  void createNews() async {
    if (_newsTitleController.text.isEmpty) return;
    await NewsService.createNews(selectedClubId!, _newsTitleController.text,
        _newsContentController.text);
  }

  void loadClubs() async {
    ClubService.getClubs().then((snapshot) {
      setState(() {
        for (var doc in snapshot.docs) {
          clubs.add(
            DropdownMenuItem(
              value: doc.id,
              child: Text(doc['name']),
            ),
          );
        }
      });
    });
  }

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
                'Create News',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              DropdownButton(
                items: clubs,
                onChanged: (value) {
                  setState(() {
                    selectedClubId = value as String;
                  });
                },
                value: selectedClubId,
              ),
              TextField(
                decoration: const InputDecoration(
                  labelText: 'News Title',
                ),
                controller: _newsTitleController,
              ),
              const SizedBox(height: 10),
              // multi-line text field
              TextField(
                decoration: const InputDecoration(
                  labelText: 'News Content',
                ),
                controller: _newsContentController,
                keyboardType: TextInputType.multiline,
                maxLines: null,
              ),
              const SizedBox(height: 10),
              MyButton(
                text: 'Post News',
                onTap: createNews,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
