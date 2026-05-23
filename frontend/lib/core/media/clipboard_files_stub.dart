import 'dart:async';

import 'media_picker_stub.dart';

typedef ClipboardFilePasteDisposer = void Function();

ClipboardFilePasteDisposer listenForClipboardFiles(
  FutureOr<void> Function(List<PickedMedia> files) onFiles,
) {
  return () {};
}
