import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/services/club_service.dart';
import 'package:flutter/material.dart';

class ClubDetailScreen extends StatelessWidget {
  final String clubId;
  final String clubName;
  const ClubDetailScreen({
    super.key,
    required this.clubId,
    required this.clubName,
  });

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
                  clubName,
                  style: const TextStyle(
                    fontSize: 25,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                StreamBuilder(
                  stream: ClubService.getMembersForClub(clubId),
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
                        final member = snapshot.data!.docs[index];
                        return ListTile(
                          title: Text(member['userId']),
                        );
                      },
                    );
                  },
                ),
              ],
            ),
          ),
        ));
  }
}
