import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';

class AuthService extends ChangeNotifier {
  Future<User> signUpWithEmailAndPassword(String email, String password) async {
    try {
      final UserCredential userCredential = await FirebaseAuth.instance
          .createUserWithEmailAndPassword(email: email, password: password);

      await FirebaseFirestore.instance
          .collection('users')
          .doc(userCredential.user!.uid)
          .set({
        'email': email,
      });

      return userCredential.user!;
    } on FirebaseAuthException {
      rethrow;
    }
  }

  Future<User> loginWithEmailAndPassword(String email, String password) async {
    try {
      final UserCredential userCredential = await FirebaseAuth.instance
          .signInWithEmailAndPassword(email: email, password: password);

      await FirebaseFirestore.instance
          .collection('users')
          .doc(userCredential.user!.uid)
          .set({
        'email': email,
      });
      return userCredential.user!;
    } on FirebaseAuthException {
      rethrow;
    }
  }

  void logout() async {
    await FirebaseAuth.instance.signOut();
  }
}
