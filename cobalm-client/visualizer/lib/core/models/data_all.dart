class DataAll {
  final String machine;
  final String key;
  final double value;
  final String timestamp;
  final List<String>? params;

  const DataAll({
    required this.machine,
    required this.key,
    required this.value,
    required this.timestamp,
    required this.params,
  });

  factory DataAll.fromJson(Map<String, dynamic> json) {
    return DataAll(
      machine: json['machine'].toString(),
      key: json['key'].toString(),
      value: json['value'].toDouble(),
      timestamp: json['ts'],
      params: json['params'],
    );
  }
}
