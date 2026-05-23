// ignore_for_file: avoid_web_libraries_in_flutter, deprecated_member_use

import 'dart:async';
import 'dart:html' as html;
import 'dart:typed_data';

class PickedMedia {
  const PickedMedia({
    required this.name,
    required this.bytes,
    this.extension = '',
  });

  final String name;
  final Uint8List bytes;
  final String extension;
}

Future<List<PickedMedia>> pickMediaFiles() async {
  final input = html.FileUploadInputElement()
    ..multiple = true
    ..accept = 'image/*,video/*';
  input.click();
  await input.onChange.first;
  final files = input.files ?? const <html.File>[];
  final picked = <PickedMedia>[];
  for (final file in files) {
    final bytes = await _readFile(file);
    picked.add(
      PickedMedia(
        name: file.name,
        bytes: bytes,
        extension: file.name.contains('.')
            ? file.name.split('.').last.toLowerCase()
            : '',
      ),
    );
  }
  return picked;
}

Future<Uint8List> _readFile(html.File file) {
  final completer = Completer<Uint8List>();
  final reader = html.FileReader();
  reader.onError.first.then((_) => completer.completeError('file_read_failed'));
  reader.onLoad.first.then((_) {
    final result = reader.result;
    if (result is ByteBuffer) {
      completer.complete(Uint8List.view(result));
      return;
    }
    completer.completeError('file_read_failed');
  });
  reader.readAsArrayBuffer(file);
  return completer.future;
}
