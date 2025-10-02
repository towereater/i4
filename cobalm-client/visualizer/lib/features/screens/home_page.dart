import 'package:flutter/material.dart';
import 'package:visualizer/core/models/data_all.dart';
import 'package:visualizer/core/models/data_sum.dart';
import 'package:visualizer/core/services/api_services.dart';
import 'package:visualizer/features/widgets/sales_bar_chart.dart';
import 'package:visualizer/features/widgets/sales_line_chart.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  late Future<List<DataAll>> futureDataAll;
  late Future<List<DataSum>> futureDataSum;
  final ApiService _apiService = ApiService(
    clientCode: '00001',
    apiCode: '45mfio4fm54jkf4wionfneil',
  );

  String machine = '00001';
  String key = 'SpindleTC';

  @override
  void initState() {
    super.initState();
    _reloadData();
  }

  void _reloadData() {
    setState(() {
      futureDataAll = _apiService.getAllData(machine, key);
      futureDataSum = _apiService.getTotalPerJob(machine);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Dashboard principale'),
        backgroundColor: const Color(0xFF2C2C44),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            tooltip: 'Ricarica dati',
            onPressed: _reloadData,
          ),
        ],
      ),
      body: Center(
        child: ListView(
          children: [
            FutureBuilder<List<DataSum>>(
              future: futureDataSum,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const CircularProgressIndicator();
                } else if (snapshot.hasError) {
                  return Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Errore: ${snapshot.error}',
                        style: const TextStyle(color: Colors.redAccent),
                      ),
                      const SizedBox(height: 20),
                      ElevatedButton(
                        onPressed: _reloadData,
                        child: const Text('Riprova'),
                      ),
                    ],
                  );
                } else if (snapshot.hasData) {
                  final data = snapshot.data!;
                  return Padding(
                    padding: const EdgeInsets.all(24.0),
                    child: AnimatedSwitcher(
                      duration: const Duration(milliseconds: 500),
                      transitionBuilder:
                          (Widget child, Animation<double> animation) {
                            return FadeTransition(
                              opacity: animation,
                              child: child,
                            );
                          },
                      child: AspectRatio(
                        aspectRatio: 1.7,
                        child: SalesBarChart(
                          machineData: data,
                          key: const ValueKey('bar_chart'),
                        ),
                      ),
                    ),
                  );
                } else {
                  return const Text('Nessun dato disponibile.');
                }
              },
            ),
            FutureBuilder<List<DataAll>>(
              future: futureDataAll,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const CircularProgressIndicator();
                } else if (snapshot.hasError) {
                  return Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Errore: ${snapshot.error}',
                        style: const TextStyle(color: Colors.redAccent),
                      ),
                      const SizedBox(height: 20),
                      ElevatedButton(
                        onPressed: _reloadData,
                        child: const Text('Riprova'),
                      ),
                    ],
                  );
                } else if (snapshot.hasData) {
                  final data = snapshot.data!;
                  return Padding(
                    padding: const EdgeInsets.all(24.0),
                    child: AnimatedSwitcher(
                      duration: const Duration(milliseconds: 500),
                      transitionBuilder:
                          (Widget child, Animation<double> animation) {
                            return FadeTransition(
                              opacity: animation,
                              child: child,
                            );
                          },
                      child: AspectRatio(
                        aspectRatio: 1.7,
                        child: SalesLineChart(
                          machineData: data,
                          key: const ValueKey('line_chart'),
                        ),
                      ),
                    ),
                  );
                } else {
                  return const Text('Nessun dato disponibile.');
                }
              },
            ),
          ],
        ),
      ),
    );
  }
}
