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

    var user = await _firestore.collection('users').doc(userId).get();

    await _firestore.collection('clubs').doc(clubId).collection('members').add({
      'userId': userId,
      'email': user.data()!['email'],
    });
  }

  static Future<void> addMemberByMail(String clubId, String email) async {
    var user = await _firestore
        .collection('users')
        .where('email', isEqualTo: email)
        .get();

    if (user.docs.isEmpty) {
      throw Exception('User not found');
    }

    await addMember(clubId, user.docs.first.id);
  }

  static Future<void> removeMember(String clubId, String userId) async {
    await _firestore
        .collection('users')
        .doc(userId)
        .collection('memberOf')
        .where('clubId', isEqualTo: clubId)
        .get()
        .then((snapshot) {
      for (var doc in snapshot.docs) {
        doc.reference.delete();
      }
    });

    await _firestore
        .collection('clubs')
        .doc(clubId)
        .collection('members')
        .where('userId', isEqualTo: userId)
        .get()
        .then((snapshot) {
      for (var doc in snapshot.docs) {
        doc.reference.delete();
      }
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
