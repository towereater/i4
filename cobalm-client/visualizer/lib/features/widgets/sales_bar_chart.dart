import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:visualizer/core/models/data_sum.dart';

class SalesBarChart extends StatelessWidget {
  final List<DataSum> machineData;

  const SalesBarChart({super.key, required this.machineData});

  @override
  Widget build(BuildContext context) {
    return BarChart(
      BarChartData(
        alignment: BarChartAlignment.spaceAround,
        maxY: _calculateMaxY(),
        barTouchData: BarTouchData(
          touchTooltipData: BarTouchTooltipData(
            getTooltipColor: (group) => Colors.indigo.withValues(alpha:0.8),
            getTooltipItem: (group, groupIndex, rod, rodIndex) {
              String value = machineData[groupIndex].value;
              String job = value.substring(value.lastIndexOf('\\') + 1);
              return BarTooltipItem(
                '$job\n',
                const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 16,
                ),
                children: <TextSpan>[
                  TextSpan(
                    text: '${rod.toY.round().toString()} s',
                    style: const TextStyle(
                      color: Colors.cyanAccent,
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              );
            },
          ),
        ),
        titlesData: FlTitlesData(
          show: true,
          topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              getTitlesWidget: _getBottomTitles,
              reservedSize: 38,
            ),
          ),
          leftTitles: AxisTitles(
            sideTitles: SideTitles(
              showTitles: true,
              reservedSize: 40,
              getTitlesWidget: _getLeftTitles,
            ),
          ),
        ),
        borderData: FlBorderData(show: false),
        barGroups: _generateBarGroups(),
        gridData: FlGridData(
          show: true,
          drawVerticalLine: false,
          getDrawingHorizontalLine: (value) => FlLine(
            color: Colors.white.withValues(alpha:0.1),
            strokeWidth: 1,
          ),
        ),
      ),
      duration: const Duration(milliseconds: 250),
    );
  }

  double _calculateMaxY() {
    double maxSum = 0;
    for (var data in machineData) {
      if (data.sum > maxSum) {
        maxSum = data.sum;
      }
    }
    return maxSum * 1.2;
  }

  List<BarChartGroupData> _generateBarGroups() {
    return List.generate(machineData.length, (index) {
      final data = machineData[index];
      return BarChartGroupData(
        x: index,
        barRods: [
          BarChartRodData(
            toY: data.sum,
            color: Colors.cyanAccent,
            width: 16,
            borderRadius: const BorderRadius.only(
              topLeft: Radius.circular(4),
              topRight: Radius.circular(4),
            ),
          ),
        ],
      );
    });
  }
  
  Widget _getBottomTitles(double value, TitleMeta meta) {
    final style = TextStyle(
      color: Colors.grey.shade400,
      fontWeight: FontWeight.bold,
      fontSize: 14,
    );
    int index = value.toInt();
    String v = machineData[index].value;
    String job = v.substring(v.lastIndexOf('\\') + 1);
    if (index >= 0 && index < machineData.length) {
      return SideTitleWidget(
        meta: meta,
        space: 16,
        child: Text(job, style: style),
      );
    }
    return Container();
  }
  
  Widget _getLeftTitles(double value, TitleMeta meta) {
    if (value % 50 != 0 && value != 0) {
      return Container();
    }
    return Text(
      '${value.toInt()} s',
      style: TextStyle(color: Colors.grey.shade400, fontSize: 12),
      textAlign: TextAlign.left,
    );
  }
}