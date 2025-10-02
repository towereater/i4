import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:visualizer/core/models/data_all.dart';

class SalesLineChart extends StatelessWidget {
  final List<DataAll> machineData;

  const SalesLineChart({super.key, required this.machineData});

  @override
  Widget build(BuildContext context) {
    return LineChart(
      LineChartData(
        gridData: FlGridData(
          show: true,
          drawVerticalLine: true,
          horizontalInterval: 50,
          verticalInterval: 1,
          getDrawingHorizontalLine: (value) => FlLine(color: Colors.white.withValues(alpha:0.1), strokeWidth: 1),
          getDrawingVerticalLine: (value) => FlLine(color: Colors.white.withValues(alpha:0.1), strokeWidth: 1),
        ),
        titlesData: FlTitlesData(
          show: true,
          rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              reservedSize: 30,
              interval: 1,
              getTitlesWidget: _bottomTitleWidgets,
            ),
          ),
          leftTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              interval: 50,
              getTitlesWidget: _leftTitleWidgets,
              reservedSize: 42,
            ),
          ),
        ),
        borderData: FlBorderData(show: true, border: Border.all(color: const Color(0xff37434d))),
        minX: 0,
        maxX: (machineData.length - 1).toDouble(),
        minY: 0,
        maxY: _calculateMaxY(),
        lineBarsData: [_mainLineBarData()],
      ),
      duration: const Duration(milliseconds: 250),
    );
  }

  double _calculateMaxY() {
    double maxValue = 0;
    for (var data in machineData) {
      if (data.value > maxValue) {
        maxValue = data.value;
      }
    }
    return (maxValue / 50).ceil() * 50;
  }
  
  Widget _bottomTitleWidgets(double value, TitleMeta meta) {
    const style = TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.grey);
    int index = value.toInt();
    if (index >= 0 && index < machineData.length) {
      return SideTitleWidget(meta: meta, child: Text(machineData[index].timestamp, style: style));
    }
    return Container();
  }

  Widget _leftTitleWidgets(double value, TitleMeta meta) {
    const style = TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Colors.grey);
    return Text('${value.toInt()} s', style: style, textAlign: TextAlign.left);
  }

  LineChartBarData _mainLineBarData() {
    return LineChartBarData(
      spots: List.generate(machineData.length, (i) => FlSpot(i.toDouble(), machineData[i].value.toDouble())),
      isCurved: true,
      gradient: const LinearGradient(colors: [Colors.cyan, Colors.blue]),
      barWidth: 5,
      isStrokeCapRound: true,
      dotData: const FlDotData(show: false),
      belowBarData: BarAreaData(
        show: true,
        gradient: LinearGradient(
          colors: [Colors.cyan.withValues(alpha:0.3), Colors.blue.withValues(alpha:0.3)],
        ),
      ),
    );
  }
}