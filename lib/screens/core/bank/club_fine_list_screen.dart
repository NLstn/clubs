import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/services/bank_service.dart';
import 'package:flutter/material.dart';

class ClubFineListScreen extends StatelessWidget {
  final String clubId;
  const ClubFineListScreen({
    super.key,
    required this.clubId,
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
              const Text(
                'Fines',
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              StreamBuilder(
                stream: BankService.getFinesAsStream(clubId),
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
                      final fine = snapshot.data!.docs[index];
                      return ListTile(
                        title: Text(fine['reason']),
                        subtitle: Text(fine['userId']),
                        trailing: Text(fine['amount'].toString()),
                      );
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
}
