import 'package:dio/dio.dart';
import '../models/data_all.dart';
import '../models/data_sum.dart';

class ApiService {
  final Dio _dio = Dio();
  
  static const String _url = 'http://localhost:12001/clients/{code}/data';
  final String clientCode;
  final String apiCode;

  ApiService({
    required this.clientCode,
    required this.apiCode,
  });

  Future<List<DataAll>> getAllData(String machine, String key) async {
    String url = '${_url.replaceAll('{code}', clientCode)}/gauge?machine=$machine&key=$key';

    try {
      final response = await _dio.get(url, options: Options(headers: {'Authentication': apiCode}));

      if (response.statusCode == 200) {
        final List<dynamic> jsonResponse = response.data;
        return jsonResponse.map((data) => DataAll.fromJson(data)).toList();
      } else {
        throw Exception('Errore nel caricamento dei dati dall\'API (URL: $url, Status code: ${response.statusCode})');
      }
    } on DioException catch (e) {
      throw Exception('Errore di connessione: $e');
    } catch (e) {
      throw Exception('Si è verificato un errore imprevisto: $e');
    }
  }

  Future<List<DataSum>> getTotalPerJob(String machine) async {
    String url = '${_url.replaceAll('{code}', clientCode)}/interval/sum?machine=$machine';

    try {
      final response = await _dio.get(url, options: Options(headers: {'Authentication': apiCode}));

      if (response.statusCode == 200) {
        final List<dynamic> jsonResponse = response.data;
        return jsonResponse.map((data) => DataSum.fromJson(data)).toList();
      } else {
        throw Exception('Errore nel caricamento dei dati dall\'API (URL: $url, Status code: ${response.statusCode})');
      }
    } on DioException catch (e) {
      throw Exception('Errore di connessione: $e');
    } catch (e) {
      throw Exception('Si è verificato un errore imprevisto: $e');
    }
  }
}
