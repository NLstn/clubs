import 'package:clubs/models/core/member.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class MembersPage extends StatefulWidget {
  const MembersPage({Key? key}) : super(key: key);

  @override
  State<MembersPage> createState() => _MembersPageState();
}

class _MembersPageState extends State<MembersPage> {
  final List<Member> members = [
    Member(
      name: 'John Doe',
      email: 'johndoe@example.com',
      joinDate: DateTime.parse('2021-01-01'),
    ),
    Member(
      name: 'Jane Doe',
      email: 'janedoe@example.com',
      joinDate: DateTime.parse('2021-02-01'),
    ),
    Member(
      name: 'Bob Smith',
      email: 'bobsmith@example.com',
      joinDate: DateTime.parse('2021-03-01'),
    ),
    Member(
      name: 'Alice Smith',
      email: 'alicesmith@example.com',
      joinDate: DateTime.parse('2021-04-01'),
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final dateFormat =
        DateFormat.yMd(Localizations.localeOf(context).languageCode);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Members'),
      ),
      body: ListView.builder(
        itemCount: members.length,
        itemBuilder: (BuildContext context, int index) {
          return ListTile(
            title: Text(members[index].name),
            subtitle: Text(members[index].email),
            trailing: Text(dateFormat.format(members[index].joinDate)),
          );
        },
      ),
    );
  }
}
