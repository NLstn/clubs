import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';

class ClubService {
  static final FirebaseAuth _auth = FirebaseAuth.instance;
  static final FirebaseFirestore _firestore = FirebaseFirestore.instance;

  static Future<void> createClub(String name) async {
    final User? user = _auth.currentUser;
    final club = _firestore.collection('clubs').doc();
    await club.set({
      'name': name,
      'owner': user!.uid,
    });

    addMember(club.id, user.uid);
  }

  static Future<void> addMember(String clubId, String userId) async {
    await _firestore
        .collection('users')
        .doc(userId)
        .collection('memberOf')
        .add({
      'clubId': clubId,
    });

    await _firestore.collection('clubs').doc(clubId).collection('members').add({
      'userId': userId,
    });
  }

  static Stream<QuerySnapshot> getClubs() {
    return _firestore
        .collection('clubs')
        .where('owner', isEqualTo: _auth.currentUser!.uid)
        .snapshots();
  }

  static Stream<QuerySnapshot> getMembersForClub(String clubId) {
    return _firestore
        .collection('clubs')
        .doc(clubId)
        .collection('members')
        .snapshots();
  }
}
