import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/screens/core/create_club_screen.dart';
import 'package:flutter/material.dart';
import 'members_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

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
            child: const Text('Members'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const MembersPage()),
              );
            },
          ),
          HomeScreenTile(
            child: const Text('Create Club'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => CreateClubScreen()),
              );
            },
          ),
          HomeScreenTile(
            child: const Text('Dummy2'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const MembersPage()),
              );
            },
          ),
          HomeScreenTile(
            child: const Text('Dummy3'),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const MembersPage()),
              );
            },
          ),
        ],
      ),
    );
  }
}

class HomeScreenTile extends StatelessWidget {
  const HomeScreenTile({
    Key? key,
    required this.child,
    required this.onTap,
  }) : super(key: key);

  final Widget child;
  final VoidCallback onTap;

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
          child: child,
        ),
      ),
    );
  }
}
