import 'package:clubs/screens/core/clubs/club_detail_screen.dart';
import 'package:clubs/services/club_service.dart';
import 'package:flutter/material.dart';

class ClubListScreen extends StatelessWidget {
  const ClubListScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Clubs'),
      ),
      body: StreamBuilder(
        stream: ClubService.getClubs(),
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
            scrollDirection: Axis.vertical,
            shrinkWrap: true,
            itemCount: snapshot.data!.docs.length,
            itemBuilder: (context, index) {
              final club = snapshot.data!.docs[index];
              return ListTile(
                title: Text(club['name']),
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => ClubDetailScreen(
                        clubId: club.id,
                        clubName: club['name'],
                      ),
                    ),
                  );
                },
              );
            },
          );
        },
      ),
    );
  }
}
