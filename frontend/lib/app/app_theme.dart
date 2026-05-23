import 'package:flutter/material.dart';

class AppColors {
  static const primary = Color(0xFF1E3A5F);
  static const secondary = Color(0xFF2F80ED);
  static const accentGold = Color(0xFFF2C94C);
  static const success = Color(0xFF16A34A);
  static const warning = Color(0xFFF59E08);
  static const error = Color(0xFFDC2626);
  static const background = Color(0xFFF7F9FC);
  static const surface = Color(0xFFFFFFFF);
  static const border = Color(0xFFE5E7EB);
  static const textPrimary = Color(0xFF1F2937);
  static const textSecondary = Color(0xFF6B7280);
  static const textDisabled = Color(0xFF9CA3AF);
  static const darkBackground = Color(0xFF07111F);
  static const darkSurface = Color(0xFF0D1B2E);
  static const darkSurfaceAlt = Color(0xFF13263D);
  static const darkBorder = Color(0xFF21354F);
  static const electricCyan = Color(0xFF4DE1FF);
  static const neonBlue = Color(0xFF3B82F6);
}

class AppTheme {
  static ThemeData get light => _build(Brightness.light);
  static ThemeData get dark => _build(Brightness.dark);

  static ThemeData _build(Brightness brightness) {
    final isDark = brightness == Brightness.dark;
    final colorScheme = ColorScheme.fromSeed(
      seedColor: AppColors.secondary,
      brightness: brightness,
      primary: isDark ? AppColors.electricCyan : AppColors.primary,
      secondary: isDark ? AppColors.neonBlue : AppColors.secondary,
      error: AppColors.error,
      surface: isDark ? AppColors.darkSurface : AppColors.surface,
    );

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      scaffoldBackgroundColor: isDark
          ? AppColors.darkBackground
          : AppColors.background,
      colorScheme: colorScheme,
      pageTransitionsTheme: const PageTransitionsTheme(
        builders: {
          TargetPlatform.android: FadeUpwardsPageTransitionsBuilder(),
          TargetPlatform.iOS: CupertinoPageTransitionsBuilder(),
          TargetPlatform.windows: FadeUpwardsPageTransitionsBuilder(),
          TargetPlatform.macOS: CupertinoPageTransitionsBuilder(),
          TargetPlatform.linux: FadeUpwardsPageTransitionsBuilder(),
        },
      ),
      appBarTheme: AppBarTheme(
        backgroundColor: isDark
            ? AppColors.darkSurface.withValues(alpha: 0.86)
            : AppColors.primary.withValues(alpha: 0.90),
        foregroundColor: Colors.white,
        centerTitle: false,
        elevation: isDark ? 0 : 10,
        shadowColor: (isDark ? Colors.black : AppColors.primary).withValues(
          alpha: isDark ? 0.24 : 0.18,
        ),
        surfaceTintColor: Colors.transparent,
        toolbarHeight: 66,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(24)),
        ),
        iconTheme: const IconThemeData(color: Colors.white, size: 24),
        actionsIconTheme: const IconThemeData(color: Colors.white, size: 23),
        titleTextStyle: const TextStyle(
          color: Colors.white,
          fontSize: 22,
          fontWeight: FontWeight.w800,
        ),
      ),
      cardTheme: CardThemeData(
        color: (isDark ? AppColors.darkSurface : AppColors.surface).withValues(
          alpha: isDark ? 0.84 : 0.94,
        ),
        elevation: isDark ? 0 : 14,
        shadowColor: (isDark ? Colors.black : AppColors.primary).withValues(
          alpha: isDark ? 0.32 : 0.12,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(18),
          side: BorderSide(
            color: isDark
                ? AppColors.darkBorder.withValues(alpha: 0.74)
                : AppColors.border.withValues(alpha: 0.7),
          ),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: isDark
            ? AppColors.darkSurfaceAlt.withValues(alpha: 0.72)
            : Colors.white,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 16,
        ),
        labelStyle: TextStyle(
          color: isDark ? const Color(0xFFB8C7DA) : AppColors.textSecondary,
        ),
        prefixIconColor: isDark ? AppColors.electricCyan : AppColors.primary,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(
            color: isDark ? AppColors.darkBorder : AppColors.border,
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(
            color: isDark ? AppColors.electricCyan : AppColors.secondary,
            width: 1.5,
          ),
        ),
      ),
      textTheme:
          (isDark ? Typography.whiteMountainView : Typography.blackMountainView)
              .apply(
                bodyColor: isDark
                    ? const Color(0xFFE6EDF7)
                    : AppColors.textPrimary,
                displayColor: isDark
                    ? const Color(0xFFF8FBFF)
                    : AppColors.textPrimary,
              ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: isDark ? AppColors.neonBlue : AppColors.secondary,
          foregroundColor: Colors.white,
          elevation: 0,
          minimumSize: const Size(120, 52),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: isDark ? AppColors.neonBlue : AppColors.primary,
          foregroundColor: Colors.white,
          minimumSize: const Size(120, 52),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: isDark ? AppColors.electricCyan : AppColors.primary,
          minimumSize: const Size(112, 48),
          side: BorderSide(
            color: isDark ? AppColors.darkBorder : AppColors.border,
          ),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
        ),
      ),
    );
  }
}
