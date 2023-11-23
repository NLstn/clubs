import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/screens/core/bank/club_fine_list_screen.dart';
import 'package:clubs/screens/core/clubs/club_member_list_screen.dart';
import 'package:clubs/screens/core/news/club_news_list_screen.dart';
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
                text: 'Members',
                onTap: () => Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) =>
                        ClubMemberListScreen(clubId: widget.clubId),
                  ),
                ),
              ),
              const SizedBox(height: 15),
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
              MyButton(
                text: 'Fines',
                onTap: () => Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) =>
                        ClubFineListScreen(clubId: widget.clubId),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
