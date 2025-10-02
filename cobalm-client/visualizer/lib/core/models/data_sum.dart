class DataSum {
  final String machine;
  final String key;
  final String value;
  final int count;
  final double sum;

  const DataSum({
    required this.machine,
    required this.key,
    required this.value,
    required this.count,
    required this.sum,
  });

  factory DataSum.fromJson(Map<String, dynamic> json) {
    return DataSum(
      machine: json['machine'].toString(),
      key: json['key'].toString(),
      value: json['value'].toString(),
      count: json['count'],
      sum: json['sum'].toDouble(),
    );
  }
}
