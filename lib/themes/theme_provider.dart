import 'package:flutter/material.dart';
import 'package:matt_lads_app/themes/light_mode.dart';
import 'package:matt_lads_app/themes/dark_mode.dart';

/*
Theme Provider
Helper class to provide the theme to the app and change themes.
*/

class ThemeProvider with ChangeNotifier{
  ThemeData _themeData = lightMode;
  ThemeData get themeData => _themeData;

  bool get isDarkMode => _themeData == darkMode;

  set themeData(ThemeData themeData){
    _themeData = themeData;
    notifyListeners();
  }

  void toggleTheme(){
    _themeData = _themeData == lightMode ? darkMode : lightMode;
    notifyListeners();
  }
}