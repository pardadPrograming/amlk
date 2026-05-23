import 'package:flutter/material.dart';

import '../../shared/responsive.dart';

class PlaceholderPage extends StatelessWidget {
  const PlaceholderPage({super.key, required this.title});
  final String title;

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: Text(title),
      body: ResponsivePage(
        child: Card(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Row(
              children: [
                Icon(
                  Icons.construction_outlined,
                  color: Theme.of(context).colorScheme.secondary,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    '$title در نسخه بعدی تکمیل می‌شود. route و ساختار ماژول آماده است.',
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
