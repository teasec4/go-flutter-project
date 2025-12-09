import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend_flutter/core/colors.dart';
import 'package:frontend_flutter/di/service_locator.dart';
import 'package:frontend_flutter/domain/repository/account_repo.dart';
import 'package:frontend_flutter/presentation/account_cubit.dart';

class AccountPage extends StatefulWidget {
  final String accountNumber;
  const AccountPage({super.key, required this.accountNumber});

  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  late AccountCubit _accountCubit;

  @override
  void initState() {
    super.initState();
    _accountCubit = AccountCubit(getIt<AccountRepo>());
    final accountId = widget.accountNumber.trim();
    if (accountId.isNotEmpty) {
      _accountCubit.getAccountInfo(accountId);
    }
  }

  @override
  void dispose() {
    _accountCubit.close();
    super.dispose();
  }

  void _showAddTransactionModal(BuildContext context, bool isDeposit) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _TransactionListeningWrapper(
        isDeposit: isDeposit,
        accountNumber: widget.accountNumber,
        onConfirm: (amount) async {
          print('onConfirm called with amount=$amount');
          if (isDeposit) {
            print('Calling deposit...');
            await _accountCubit.deposit(widget.accountNumber, amount);
            print('Deposit completed');
          } else {
            print('Calling withdraw...');
            await _accountCubit.withdraw(widget.accountNumber, amount);
            print('Withdraw completed');
          }
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("My Account"),
        centerTitle: true,
        elevation: 0,
        leading: IconButton(
          onPressed: () => Navigator.pop(context),
          icon: const Icon(Icons.arrow_back),
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: BlocProvider.value(
            value: _accountCubit,
            child: BlocBuilder<AccountCubit, AccountState>(
              builder: (context, accountState) {
                return switch (accountState) {
                  AccountInitial() => const Center(
                      child: Text('Loading account info...'),
                    ),
                  AccountLoadingData() => const Center(
                      child: CircularProgressIndicator(),
                    ),
                  AccountLoadedData(:final account) => Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          // Account Info Section
                          Container(
                            padding: const EdgeInsets.all(20),
                            decoration: BoxDecoration(
                              color: AppColors.surfaceVariant,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: AppColors.border,
                              ),
                            ),
                            child: Column(
                              children: [
                                // Account Number
                                _InfoRow(
                                  label: 'Account Number',
                                  value: account.id,
                                ),
                                Divider(color: AppColors.divider, height: 24),
                                // Balance
                                _InfoRow(
                                  label: 'Current Balance',
                                  value: '\$${account.balance}',
                                  valueStyle: Theme.of(context)
                                      .textTheme
                                      .headlineSmall
                                      ?.copyWith(
                                        fontWeight: FontWeight.bold,
                                        color: AppColors.success,
                                      ),
                                ),
                                if (account.name != null &&
                                    account.name != 'Unknown') ...[
                                  Divider(
                                    color: AppColors.divider,
                                    height: 24,
                                  ),
                                  _InfoRow(
                                    label: 'Account Holder',
                                    value: account.name!,
                                  ),
                                ],
                              ],
                            ),
                          ),
                          const SizedBox(height: 32),
                          // Action Buttons
                          Row(
                            children: [
                              Expanded(
                                child: _TransactionButton(
                                  icon: Icons.arrow_downward,
                                  label: 'Withdraw',
                                  onPressed: () =>
                                      _showAddTransactionModal(context, false),
                                  color: AppColors.danger,
                                ),
                              ),
                              const SizedBox(width: 16),
                              Expanded(
                                child: _TransactionButton(
                                  icon: Icons.arrow_upward,
                                  label: 'Deposit',
                                  onPressed: () =>
                                      _showAddTransactionModal(context, true),
                                  color: AppColors.success,
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  AccountHadError(:final message) => Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(
                            Icons.error_outline,
                            color: AppColors.danger,
                            size: 48,
                          ),
                          const SizedBox(height: 16),
                          Text(
                            'Error',
                            style: Theme.of(context).textTheme.headlineSmall,
                          ),
                          const SizedBox(height: 8),
                          Text(
                            message,
                            textAlign: TextAlign.center,
                            style: Theme.of(context).textTheme.bodyMedium,
                          ),
                        ],
                      ),
                    ),
                };
              },
            ),
          ),
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final String label;
  final String value;
  final TextStyle? valueStyle;

  const _InfoRow({
    required this.label,
    required this.value,
    this.valueStyle,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: AppColors.textSecondary,
              ),
        ),
        Text(
          value,
          style: valueStyle ??
              Theme.of(context).textTheme.bodyLarge?.copyWith(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
        ),
      ],
    );
  }
}

class _TransactionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onPressed;
  final Color color;

  const _TransactionButton({
    required this.icon,
    required this.label,
    required this.onPressed,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: onPressed,
      icon: Icon(icon),
      label: Text(label),
      style: ElevatedButton.styleFrom(
        backgroundColor: color,
        foregroundColor: AppColors.surface,
        padding: const EdgeInsets.symmetric(vertical: 16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
      ),
    );
  }
}

class _AddTransactionSheet extends StatefulWidget {
  final bool isDeposit;
  final String accountNumber;
  final Function(int) onConfirm;

  const _AddTransactionSheet({
    required this.isDeposit,
    required this.accountNumber,
    required this.onConfirm,
  });

  @override
  State<_AddTransactionSheet> createState() => _AddTransactionSheetState();
}

class _TransactionListeningWrapper extends StatelessWidget {
  final bool isDeposit;
  final String accountNumber;
  final Function(int) onConfirm;

  const _TransactionListeningWrapper({
    required this.isDeposit,
    required this.accountNumber,
    required this.onConfirm,
  });

  @override
  Widget build(BuildContext context) {
    return BlocListener<AccountCubit, AccountState>(
      listener: (context, state) {
        if (state is AccountLoadedData) {
          print('Transaction successful, showing toast');
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                isDeposit ? "Successfully deposited" : "Successfully withdrew",
              ),
              backgroundColor: AppColors.success,
              behavior: SnackBarBehavior.floating,
              duration: const Duration(seconds: 2),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            ),
          );
          // Close modal after showing toast
          Navigator.pop(context);
        } else if (state is AccountHadError) {
          print('Transaction error: ${state.message}');
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Error: ${state.message}'),
              backgroundColor: AppColors.danger,
              behavior: SnackBarBehavior.floating,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            ),
          );
        }
      },
      child: _AddTransactionSheet(
        isDeposit: isDeposit,
        accountNumber: accountNumber,
        onConfirm: onConfirm,
      ),
    );
  }
}

class _AddTransactionSheetState extends State<_AddTransactionSheet> {
  late final TextEditingController _amountController;

  @override
  void initState() {
    super.initState();
    _amountController = TextEditingController();
  }

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  void _handleConfirmation() async {
    print('_handleConfirmation called');
    if (_amountController.text.isEmpty) {
      print('Amount is empty');
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter an amount')),
      );
      return;
    }

    try {
      int amount = int.parse(_amountController.text);
      if (amount <= 0) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Amount must be greater than 0')),
        );
        return;
      }

      print('Starting ${widget.isDeposit ? 'deposit' : 'withdraw'} for amount: $amount');
      await widget.onConfirm(amount);
      print('${widget.isDeposit ? 'Deposit' : 'Withdraw'} completed successfully');
    } on FormatException {
      print('FormatException occurred');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Please enter a valid number')),
        );
      }
    } catch (e) {
      print('Exception occurred: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: ${e.toString()}')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return DraggableScrollableSheet(
      expand: false,
      minChildSize: 0.4,
      initialChildSize: 0.5,
      builder: (context, scrollController) {
        return Container(
          decoration: const BoxDecoration(
            color: AppColors.surface,
            borderRadius: BorderRadius.vertical(
              top: Radius.circular(20),
            ),
          ),
          child: SingleChildScrollView(
            controller: scrollController,
            child: Padding(
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom + 16,
                left: 24,
                right: 24,
                top: 24,
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Handle indicator
                  Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: AppColors.divider,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                  const SizedBox(height: 24),
                  // Title
                  Text(
                    widget.isDeposit ? 'Deposit Money' : 'Withdraw Money',
                    style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  const SizedBox(height: 24),
                  // Amount TextField
                  TextField(
                    controller: _amountController,
                    keyboardType:
                        const TextInputType.numberWithOptions(decimal: true),
                    decoration: InputDecoration(
                      labelText: 'Amount',
                      hintText: 'Enter amount',
                      prefixIcon: const Icon(Icons.attach_money),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 16,
                        vertical: 16,
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),
                  // Confirm Button
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () {
                        _handleConfirmation();
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: widget.isDeposit
                            ? AppColors.success
                            : AppColors.danger,
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: Text(
                        widget.isDeposit ? 'Deposit' : 'Withdraw',
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                          color: AppColors.surface,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),
                  // Cancel Button
                  SizedBox(
                    width: double.infinity,
                    child: OutlinedButton(
                      onPressed: () => Navigator.pop(context),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                      child: const Text('Cancel'),
                    ),
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
