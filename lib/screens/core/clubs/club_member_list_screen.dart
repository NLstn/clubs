import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/services/club_service.dart';
import 'package:flutter/material.dart';

class ClubMemberListScreen extends StatefulWidget {
  final String clubId;
  const ClubMemberListScreen({
    super.key,
    required this.clubId,
  });

  @override
  State<ClubMemberListScreen> createState() => _ClubMemberListScreenState();
}

class _ClubMemberListScreenState extends State<ClubMemberListScreen> {
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
              const Text(
                'Members',
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              MyButton(text: 'Add Member', onTap: openAddMemberPopup),
              const SizedBox(height: 15),
              StreamBuilder(
                stream: ClubService.getMembersForClubAsStream(widget.clubId),
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
