import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/screens/core/news/create_news_screen.dart';
import 'package:clubs/services/news_service.dart';
import 'package:flutter/material.dart';

class ClubNewsListScreen extends StatelessWidget {
  final String clubId;
  const ClubNewsListScreen({
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
                'News',
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              MyButton(
                text: 'Create News',
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => CreateNewsScreen(
                        clubId: clubId,
                      ),
                    ),
                  );
                },
              ),
              const SizedBox(height: 15),
              StreamBuilder(
                stream: NewsService.getNewsAsStream(clubId),
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
                    physics: const NeverScrollableScrollPhysics(),
                    shrinkWrap: true,
                    itemCount: snapshot.data!.docs.length,
                    itemBuilder: (context, index) {
                      final news = snapshot.data!.docs[index];
                      return ListTile(
                        title: Text(news['title']),
                        subtitle: Text(news['content']),
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
