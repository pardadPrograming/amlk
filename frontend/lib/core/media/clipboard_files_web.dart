// ignore_for_file: avoid_web_libraries_in_flutter, deprecated_member_use

import 'dart:async';
import 'dart:html' as html;
import 'dart:typed_data';

import 'media_picker_web.dart';

typedef ClipboardFilePasteDisposer = void Function();

ClipboardFilePasteDisposer listenForClipboardFiles(
  FutureOr<void> Function(List<PickedMedia> files) onFiles,
) {
  final subscription = html.document.onPaste.listen((event) async {
    final files = event.clipboardData?.files ?? const <html.File>[];
    if (files.isEmpty) {
      return;
    }
    event.preventDefault();
    final picked = <PickedMedia>[];
    for (final file in files) {
      picked.add(
        PickedMedia(
          name: file.name.isEmpty ? 'clipboard-file' : file.name,
          bytes: await _readFile(file),
          extension: file.name.contains('.')
              ? file.name.split('.').last.toLowerCase()
              : '',
        ),
      );
    }
    await onFiles(picked);
  });
  return subscription.cancel;
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
