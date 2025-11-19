
class Account{
  final String id;
  final String? name;
  final double balance;

  Account({required this.id, required this.name, required this.balance});

  // create Account instance from JSON data
  factory Account.fromJSON(Map<String, dynamic> json){
    return Account(
        id: json['accountId'].toString(),
        name: json['name'] as String? ?? 'Unknown',
        balance: (json['balance'] as num).toDouble(),
    );
  }

  Map<String, dynamic> toJson(){
    return {
      'id' : id,
      'name' : name,
      'balance' : balance,
    };
  }
}