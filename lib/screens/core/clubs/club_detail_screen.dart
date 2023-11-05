import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/screens/core/news/club_news_list_screen.dart';
import 'package:clubs/services/club_service.dart';
import 'package:flutter/material.dart';

class ClubDetailScreen extends StatefulWidget {
  final String clubId;
  final String clubName;
  const ClubDetailScreen({
    super.key,
    required this.clubId,
    required this.clubName,
  });

  @override
  State<ClubDetailScreen> createState() => _ClubDetailScreenState();
}

class _ClubDetailScreenState extends State<ClubDetailScreen> {
  final TextEditingController _emailController = TextEditingController();

  void openAddMemberPopup() {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Add Member'),
          content: TextField(
            controller: _emailController,
            decoration: const InputDecoration(
              labelText: 'Email',
            ),
          ),
          actions: [
            TextButton(
              onPressed: () async {
                await ClubService.addMemberByMail(
                    widget.clubId, _emailController.text);
                if (mounted) {
                  Navigator.pop(context);
                }
              },
              child: const Text('Add'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const MyAppBar(),
      body: Padding(
        padding: const EdgeInsets.all(25.0),
        child: Center(
          child: Column(
            children: [
              Text(
                widget.clubName,
                style: const TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              MyButton(
                text: 'News',
                onTap: () => Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) =>
                        ClubNewsListScreen(clubId: widget.clubId),
                  ),
                ),
              ),
              const SizedBox(height: 15),
              MyButton(text: 'Add Member', onTap: openAddMemberPopup),
              StreamBuilder(
                stream: ClubService.getMembersForClub(widget.clubId),
                builder: (context, snapshot) {
                  if (snapshot.hasError) {
                    return const Center(
                      child: Text('Something went wrong'),
                    );
                  }

                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(
                      child: CircularProgressIndicator(),
                    );
                  }

                  return ListView.builder(
                    shrinkWrap: true,
                    itemCount: snapshot.data!.docs.length,
                    itemBuilder: (context, index) {
                      final member = snapshot.data!.docs[index];
                      return _buildMemberTile(member);
                    },
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMemberTile(QueryDocumentSnapshot<Object?> member) {
    return ListTile(
      title: Text(member['email']),
      trailing: IconButton(
        icon: const Icon(Icons.delete),
        onPressed: () async {
          await ClubService.removeMember(widget.clubId, member['userId']);
        },
      ),
    );
  }
}
