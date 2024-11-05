// go_server_macos.dart
import 'dart:ffi';
import 'dart:io';

typedef goServerType = Void Function();
typedef goServerFunc = void Function();

void startGoServer() async {
  final lib = DynamicLibrary.open('${Directory(Platform.resolvedExecutable).parent.path}/../Resources/goServer.so');
  final goServerFunc goServer = lib.lookup<NativeFunction<goServerType>>('goServer').asFunction();
  goServer();
}