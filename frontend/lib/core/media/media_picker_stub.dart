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

Future<List<PickedMedia>> pickMediaFiles() async => const [];
