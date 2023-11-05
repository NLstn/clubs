import 'package:cloud_firestore/cloud_firestore.dart';

class NewsService {
  static Future<void> createNews(
      String clubId, String newsTitle, String newsContent) async {
    FirebaseFirestore.instance.collection('news').add({
      'clubId': clubId,
      'title': newsTitle,
      'content': newsContent,
    });
  }

  static Stream<QuerySnapshot> getNewsAsStream(String clubId) {
    return FirebaseFirestore.instance
        .collection('news')
        .where('clubId', isEqualTo: clubId)
        .snapshots();
  }
}
