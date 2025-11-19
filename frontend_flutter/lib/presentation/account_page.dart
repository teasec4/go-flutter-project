import 'package:flutter/material.dart';
import 'package:frontend_flutter/data/datasource/account_remote_datasource_impl.dart';
import 'package:frontend_flutter/domain/model/account.dart';

class AccountPage extends StatefulWidget {
  final String accountNumber ;
  const AccountPage({super.key, required this.accountNumber});

  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  late Future<Account> account;

  @override
  void initState(){
    super.initState();
    account = AccountRemoteDatasourceImpl().getAccountInfo(widget.accountNumber);
  }
  // show a modal bottom
  void _showAddTransactionModal(BuildContext context, bool isDeposit){
    showModalBottomSheet(
        context: context,
        isScrollControlled: true,
        builder: (context) => _AddTransactionSheet(
            isDeposit: isDeposit
        )
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Account number ${widget.accountNumber}"),
        centerTitle: true,
        leading: IconButton(
          onPressed: () {
            Navigator.pop(context);
          },
          icon: Icon(Icons.arrow_back_sharp),
        ),
      ),
      body: SafeArea(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(15.0),
            child: FutureBuilder<Account>(
              future: account,
              builder: (context, snapshot){
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return Center(child: CircularProgressIndicator());
                }
                if (snapshot.hasError) {
                  return Center(child: Text('Error: ${snapshot.error}'));
                }
                final accountData = snapshot.data!;
                return Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Text("Account Number - ${accountData.id}")
                      ],
                    ),
                    SizedBox(height: 20,),
                    Text("Balance - ${accountData.balance}"),
                    SizedBox(height: 20,),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        ElevatedButton(
                            onPressed: () => _showAddTransactionModal(context, false),
                            child: Text("Withdraw")
                        ),
                        ElevatedButton(
                            onPressed: () =>  _showAddTransactionModal(context, true),
                            child: Text("Deposit")
                        ),
                      ],
                    )
                  ],
                );
              }
            ),
          ),
        ),
      ),
    );
  }
}


// modal bottom sheet
class _AddTransactionSheet extends StatefulWidget{
  final bool isDeposit;
  const _AddTransactionSheet({required this.isDeposit});

  @override
  State<_AddTransactionSheet> createState() => _AddTransactionSheetState();

}

class _AddTransactionSheetState extends State<_AddTransactionSheet>{
  late final TextEditingController _amountController;

  @override
  void initState(){
    super.initState();
    _amountController = TextEditingController();
  }

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  // handle confirmation
  void _handleConfirmation(bool isDeposit){
    ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
            content: isDeposit ? Text("Success Deposit") : Text("Success Withdraw"),
          backgroundColor: Colors.green.shade300,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8)
          ),

        )
    );
    Navigator.pop(context);
    _amountController.clear();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
        left: 16,
        right: 16,
        top: 16,
      ),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            SizedBox(height:20),
            TextField(
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: const InputDecoration(
                labelText: "Amount",
                hintText: "Enter Amount",
                prefixIcon: Icon(Icons.attach_money),
              ),
            ),
            SizedBox(height:20),
            ElevatedButton(
                onPressed: () => _handleConfirmation(widget.isDeposit),
                child: widget.isDeposit ? Text("Deposit") : Text("Withdraw"),
            ),
            SizedBox(height:20),
        
          ],
        ),
      ),
    );
  }
}