import 'package:flutter/material.dart';
import 'analysis_screen.dart';
import 'history_screen.dart';

/// Home screen with navigation
class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  static const List<Widget> _pages = <Widget>[
    AnalysisScreen(),
    HistoryScreen(),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            Icon(Icons.search, color: Colors.white),
            const SizedBox(width: 8),
            const Text(
              'SourceTracer',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ],
        ),
        backgroundColor: Colors.blue[700],
        elevation: 2,
        actions: [
          IconButton(
            icon: const Icon(Icons.info_outline, color: Colors.white),
            onPressed: () => _showAboutDialog(context),
            tooltip: 'About',
          ),
        ],
      ),
      body: IndexedStack(
        index: _selectedIndex,
        children: _pages,
      ),
      bottomNavigationBar: BottomNavigationBar(
        items: const <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: Icon(Icons.analytics),
            label: '分析',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.history),
            label: '履歴',
          ),
        ],
        currentIndex: _selectedIndex,
        selectedItemColor: Colors.blue[700],
        onTap: _onItemTapped,
      ),
    );
  }

  void _showAboutDialog(BuildContext context) {
    showAboutDialog(
      context: context,
      applicationName: 'SourceTracer',
      applicationVersion: '1.0.0',
      applicationIcon: const Icon(Icons.search, size: 48, color: Colors.blue),
      children: [
        const Text('情報源を、徹底的に追う。'),
        const SizedBox(height: 16),
        const Text(
          '自動文献収集・事実検証・情報源追跡システム\n\n'
          '意見と事実を峻別し、主張をサポートする文献を自動収集します。',
        ),
      ],
    );
  }
}
