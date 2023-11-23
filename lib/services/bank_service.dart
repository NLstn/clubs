import 'package:cloud_firestore/cloud_firestore.dart';

class BankService {
  static final FirebaseFirestore _firestore = FirebaseFirestore.instance;

  static Future<void> addFine(
      String clubId, String userId, String reason, int amount) async {
    await _firestore.collection('clubs').doc(clubId).collection('fines').add({
      'userId': userId,
      'reason': reason,
      'amount': amount,
      'createdAt': DateTime.now(),
    });
  }

  static Stream<QuerySnapshot> getFinesAsStream(String clubId) {
    return _firestore
        .collection('clubs')
        .doc(clubId)
        .collection('fines')
        .orderBy('createdAt', descending: true)
        .snapshots();
  }
}
