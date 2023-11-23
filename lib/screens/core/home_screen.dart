import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/screens/core/bank/create_fine_screen.dart';
import 'package:clubs/screens/core/clubs/club_list_screen.dart';
import 'package:clubs/screens/core/clubs/create_club_screen.dart';
import 'package:clubs/screens/core/news/create_news_screen.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const MyAppBar(),
      body: GridView.extent(
        maxCrossAxisExtent: 100,
        crossAxisSpacing: 5,
        mainAxisSpacing: 5,
        padding: const EdgeInsets.all(20),
        children: [
          HomeScreenTile(
            text: 'Create Club',
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => CreateClubScreen()),
              );
            },
          ),
          HomeScreenTile(
            text: 'List Clubs',
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const ClubListScreen()),
              );
            },
          ),
          HomeScreenTile(
            text: 'Create News',
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                    builder: (context) => const CreateNewsScreen()),
              );
            },
          ),
          HomeScreenTile(
            text: 'Create Fine',
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                    builder: (context) => const CreateFineScreen()),
              );
            },
          )
        ],
      ),
    );
  }
}

class HomeScreenTile extends StatelessWidget {
  final String text;
  final VoidCallback onTap;

  const HomeScreenTile({
    super.key,
    required this.text,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          border: Border.all(color: Colors.grey),
          borderRadius: BorderRadius.circular(15),
        ),
        child: Center(
          child: Text(text, textAlign: TextAlign.center),
        ),
      ),
    );
  }
}
