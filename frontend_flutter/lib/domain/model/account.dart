
class Account{
  final String id;
  final String name;
  final double balance;

  Account({required this.id, required this.name, required this.balance});

  // create Account instance from JSON data
  factory Account.fromJSON(Map<String, dynamic> json){
    return Account(
        id: json['id'] as String,
        name: json['name'] as String,
        balance: json['balance'] as double,
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