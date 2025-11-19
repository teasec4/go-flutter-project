import 'package:flutter/material.dart';
import 'account_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  late TextEditingController _idNumberController;

  @override
  void initState(){
    super.initState();
    _idNumberController = TextEditingController();
  }

  @override
  void dispose() {
    _idNumberController.dispose();
    super.dispose();
  }

  void routeGoNextHandler(String accountNumber){
    Navigator.push(context, MaterialPageRoute(builder:
    (context) => AccountPage(accountNumber: accountNumber)
    ));
    _idNumberController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Go Flutter App"),
        centerTitle: true,
      ),
      body: SafeArea(
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              SizedBox(
                width: 300,
                child: TextField(
                  controller: _idNumberController,
                  keyboardType: TextInputType.numberWithOptions(
                    decimal: false
                  ),
                  decoration: InputDecoration(
                    labelText: "Account Number",
                    hintText: "Enter account number",
                    prefixIcon: const Icon(Icons.numbers),
                  ),
                ),
              ),
              SizedBox(height: 20,),
              SizedBox(
                width: 300,
                child: ElevatedButton(
                  onPressed: () {
                    routeGoNextHandler(_idNumberController.text);
                  },
                  child: Text("Next"),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
