import 'package:clubs/components/my_app_bar.dart';
import 'package:clubs/components/my_button.dart';
import 'package:clubs/services/bank_service.dart';
import 'package:clubs/services/club_service.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class CreateFineScreen extends StatefulWidget {
  const CreateFineScreen({super.key});

  @override
  State<CreateFineScreen> createState() => _CreateFineScreenState();
}

class _CreateFineScreenState extends State<CreateFineScreen> {
  List<DropdownMenuItem> clubs = [];
  String? selectedClubId;

  List<DropdownMenuItem> members = [];
  String? selectedMemberId;

  final TextEditingController _reasonController = TextEditingController();
  final TextEditingController _amountController = TextEditingController();

  @override
  void initState() {
    super.initState();
    loadClubs();
  }

  void loadClubs() async {
    ClubService.getClubs().then((snapshot) {
      setState(() {
        for (var doc in snapshot.docs) {
          clubs.add(
            DropdownMenuItem(
              value: doc.id,
              child: Text(doc['name']),
            ),
          );
        }
      });
    });
  }

  void setSelectedClub(String clubId) {
    setState(() {
      selectedClubId = clubId;
    });

    ClubService.getMembersForClub(clubId).then((snapshot) {
      setState(() {
        for (var doc in snapshot.docs) {
          members.add(
            DropdownMenuItem(
              value: doc.id,
              child: Text(doc['email']),
            ),
          );
        }
      });
    });
  }

  void addFine() async {
    BankService.addFine(
      selectedClubId!,
      selectedMemberId!,
      _reasonController.text,
      int.parse(_amountController.text),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const MyAppBar(),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 25.0),
        child: Center(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text(
                'Create News',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 10),
              DropdownButton(
                items: clubs,
                onChanged: (value) => setSelectedClub(value.toString()),
                value: selectedClubId,
              ),
              const SizedBox(height: 10),
              DropdownButton(
                items: members,
                onChanged: (value) => setState(() {
                  selectedMemberId = value.toString();
                }),
                value: selectedMemberId,
              ),
              const SizedBox(height: 10),
              TextField(
                controller: _reasonController,
                decoration: const InputDecoration(
                  labelText: 'Reason',
                ),
              ),
              const SizedBox(height: 10),
              TextField(
                controller: _amountController,
                decoration: const InputDecoration(
                  labelText: 'Amount',
                ),
                keyboardType: TextInputType.number,
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                ],
              ),
              const SizedBox(height: 10),
              MyButton(
                onTap: addFine,
                text: 'Create',
              ),
            ],
          ),
        ),
      ),
    );
  }
}
