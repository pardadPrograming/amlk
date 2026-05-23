import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../app/app.dart';
import '../app/app_theme.dart';
import '../app/theme_controller.dart';
import '../features/consultants/consultants_controller.dart';

class ResponsivePage extends StatelessWidget {
  const ResponsivePage({super.key, required this.child, this.maxWidth = 1180});
  final Widget child;
  final double maxWidth;

  @override
  Widget build(BuildContext context) {
    return AnimatedGradientBackground(
      child: Align(
        alignment: Alignment.topCenter,
        child: ConstrainedBox(
          constraints: BoxConstraints(maxWidth: maxWidth),
          child: Padding(
            padding: EdgeInsets.all(
              MediaQuery.sizeOf(context).width < 700 ? 16 : 24,
            ),
            child: AnimatedEntrance(child: child),
          ),
        ),
      ),
    );
  }
}

class AnimatedGradientBackground extends StatefulWidget {
  const AnimatedGradientBackground({super.key, required this.child});
  final Widget child;

  @override
  State<AnimatedGradientBackground> createState() =>
      _AnimatedGradientBackgroundState();
}

class _AnimatedGradientBackgroundState extends State<AnimatedGradientBackground>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 10),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, _) {
        final t = Curves.easeInOut.transform(_controller.value);
        return Container(
          width: double.infinity,
          height: double.infinity,
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.lerp(Alignment.topRight, Alignment.topLeft, t)!,
              end: Alignment.lerp(
                Alignment.bottomLeft,
                Alignment.bottomRight,
                t,
              )!,
              colors: isDark
                  ? const [
                      Color(0xFF050B14),
                      Color(0xFF0A1F36),
                      Color(0xFF102A45),
                      Color(0xFF07111F),
                    ]
                  : const [
                      Color(0xFFF7F9FC),
                      Color(0xFFEAF2FF),
                      Color(0xFFFFFAE8),
                      Color(0xFFF7F9FC),
                    ],
              stops: const [0, 0.42, 0.78, 1],
            ),
          ),
          child: widget.child,
        );
      },
    );
  }
}

class AnimatedEntrance extends StatelessWidget {
  const AnimatedEntrance({super.key, required this.child});
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0, end: 1),
      duration: const Duration(milliseconds: 520),
      curve: Curves.easeOutCubic,
      builder: (context, value, child) {
        return Opacity(
          opacity: value,
          child: Transform.translate(
            offset: Offset(0, 18 * (1 - value)),
            child: child,
          ),
        );
      },
      child: child,
    );
  }
}

class GradientHeader extends StatelessWidget {
  const GradientHeader({
    super.key,
    required this.title,
    required this.subtitle,
    this.icon = Icons.apartment_outlined,
  });

  final String title;
  final String subtitle;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final gradientColors = isDark
        ? const [Color(0xFF102A45), Color(0xFF2563EB), Color(0xFF4DE1FF)]
        : const [Color(0xFF172A4A), Color(0xFF2F80ED), Color(0xFF48D7FF)];

    return ClipRRect(
      borderRadius: BorderRadius.circular(24),
      child: Stack(
        children: [
          Container(
            padding: const EdgeInsets.all(22),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: gradientColors,
                begin: Alignment.topRight,
                end: Alignment.bottomLeft,
              ),
              boxShadow: [
                BoxShadow(
                  color: (isDark ? AppColors.electricCyan : AppColors.secondary)
                      .withValues(alpha: isDark ? 0.16 : 0.26),
                  blurRadius: 34,
                  offset: const Offset(0, 18),
                ),
              ],
            ),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 8, sigmaY: 8),
              child: Row(
                children: [
                  Container(
                    width: 56,
                    height: 56,
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.13),
                      borderRadius: BorderRadius.circular(18),
                      border: Border.all(
                        color: Colors.white.withValues(alpha: 0.24),
                      ),
                    ),
                    child: Icon(
                      icon,
                      color: isDark ? Colors.white : AppColors.accentGold,
                      size: 30,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          title,
                          style: Theme.of(context).textTheme.headlineSmall
                              ?.copyWith(
                                color: Colors.white,
                                fontWeight: FontWeight.w900,
                              ),
                        ),
                        const SizedBox(height: 7),
                        Text(
                          subtitle,
                          style: Theme.of(context).textTheme.bodyMedium
                              ?.copyWith(
                                color: Colors.white.withValues(alpha: 0.82),
                                height: 1.45,
                              ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          Positioned(
            left: -18,
            top: 18,
            child: Transform.rotate(
              angle: -0.38,
              child: Container(
                width: 180,
                height: 18,
                color: Colors.white.withValues(alpha: isDark ? 0.08 : 0.12),
              ),
            ),
          ),
          Positioned(
            right: 26,
            bottom: 18,
            child: Transform.rotate(
              angle: -0.38,
              child: Container(
                width: 220,
                height: 2,
                color: AppColors.accentGold.withValues(
                  alpha: isDark ? 0.24 : 0.38,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class ThemeModeAction extends StatelessWidget {
  const ThemeModeAction({super.key});

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<ThemeController>();
    return Obx(
      () => IconButton(
        tooltip: 'تم: ${controller.label}',
        onPressed: controller.cycle,
        icon: Icon(controller.icon),
      ),
    );
  }
}

class PanelScaffold extends StatelessWidget {
  const PanelScaffold({
    super.key,
    required this.title,
    required this.body,
    this.floatingActionButton,
  });

  final Widget title;
  final Widget body;
  final Widget? floatingActionButton;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final desktop = constraints.maxWidth >= 1024;
        if (!desktop) {
          return Scaffold(
            drawer: const AmlakDrawer(),
            appBar: GlassAppBar(title: title),
            floatingActionButton: floatingActionButton,
            body: body,
          );
        }
        final navController = Get.put(
          PanelNavigationController(),
          permanent: true,
        );
        return Obx(
          () => Scaffold(
            floatingActionButton: floatingActionButton,
            body: Row(
              children: [
                DesktopSideNav(expanded: navController.expanded.value),
                Expanded(
                  child: Scaffold(
                    backgroundColor: Colors.transparent,
                    appBar: GlassAppBar(
                      title: title,
                      leading: IconButton(
                        tooltip: navController.expanded.value
                            ? 'جمع کردن منو'
                            : 'باز کردن منو',
                        onPressed: navController.toggle,
                        icon: Icon(
                          navController.expanded.value
                              ? Icons.menu_open_rounded
                              : Icons.menu_rounded,
                        ),
                      ),
                    ),
                    body: body,
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}

class PanelNavigationController extends GetxController {
  final expanded = true.obs;

  void toggle() => expanded.toggle();
}

class DesktopSideNav extends StatelessWidget {
  const DesktopSideNav({super.key, required this.expanded});

  final bool expanded;

  static const _baseItems = [
    _NavItem(
      title: 'داشبورد',
      icon: Icons.dashboard_customize_outlined,
      route: AppRoutes.dashboard,
    ),
    _NavItem(
      title: 'فایل‌ها',
      icon: Icons.apartment_outlined,
      route: AppRoutes.properties,
    ),
    _NavItem(
      title: 'جدیدترین فایل‌ها',
      icon: Icons.dynamic_feed_outlined,
      route: AppRoutes.latestFiles,
    ),
    _NavItem(
      title: 'چت‌ها',
      icon: Icons.chat_bubble_outline_rounded,
      route: AppRoutes.chats,
    ),
    _NavItem(
      title: 'صندوقچه‌ها',
      icon: Icons.inventory_2_outlined,
      route: AppRoutes.vaults,
    ),
    _NavItem(
      title: 'مشاورین',
      icon: Icons.manage_accounts_outlined,
      route: AppRoutes.consultants,
    ),
    _NavItem(
      title: 'تنظیمات',
      icon: Icons.settings_outlined,
      route: AppRoutes.settings,
    ),
    _NavItem(
      title: 'مدیریت کل',
      icon: Icons.admin_panel_settings_outlined,
      route: AppRoutes.admin,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final consultants = Get.find<ConsultantsController>();
    return AnimatedContainer(
      duration: const Duration(milliseconds: 260),
      curve: Curves.easeOutCubic,
      width: expanded ? 264 : 84,
      height: double.infinity,
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topRight,
          end: Alignment.bottomLeft,
          colors: isDark
              ? const [Color(0xFF07111F), Color(0xFF0D2138)]
              : const [Color(0xFFFFFFFF), Color(0xFFEAF2FF)],
        ),
        border: Border(
          left: BorderSide(
            color: (isDark ? Colors.white : AppColors.border).withValues(
              alpha: isDark ? 0.08 : 0.8,
            ),
          ),
        ),
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(12, 14, 12, 12),
          child: Column(
            children: [
              _SideBrand(expanded: expanded),
              const SizedBox(height: 20),
              Expanded(
                child: Obx(() {
                  final items = [
                    ..._baseItems.take(5),
                    if (consultants.pendingInboxCount > 0)
                      const _NavItem(
                        title: 'دعوت‌نامه‌ها',
                        icon: Icons.inbox_outlined,
                        route: AppRoutes.inbox,
                      ),
                    ..._baseItems.skip(5),
                  ];
                  return ListView.separated(
                    itemCount: items.length,
                    separatorBuilder: (_, _) => const SizedBox(height: 8),
                    itemBuilder: (context, index) {
                      final item = items[index];
                      return _SideNavButton(
                        item: item,
                        expanded: expanded,
                        badgeCount: item.route == AppRoutes.inbox
                            ? consultants.pendingInboxCount
                            : 0,
                      );
                    },
                  );
                }),
              ),
              const SizedBox(height: 10),
              _SideThemeButton(expanded: expanded),
            ],
          ),
        ),
      ),
    );
  }
}

class _SideBrand extends StatelessWidget {
  const _SideBrand({required this.expanded});

  final bool expanded;

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: const Duration(milliseconds: 240),
      height: 58,
      padding: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [AppColors.primary, AppColors.secondary],
          begin: Alignment.topRight,
          end: Alignment.bottomLeft,
        ),
        borderRadius: BorderRadius.circular(18),
        boxShadow: [
          BoxShadow(
            color: AppColors.secondary.withValues(alpha: 0.20),
            blurRadius: 18,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: expanded
            ? MainAxisAlignment.start
            : MainAxisAlignment.center,
        children: [
          const Icon(Icons.real_estate_agent_outlined, color: Colors.white),
          if (expanded) ...[
            const SizedBox(width: 10),
            const Expanded(
              child: Text(
                'Amlak CRM',
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w900,
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _SideNavButton extends StatelessWidget {
  const _SideNavButton({
    required this.item,
    required this.expanded,
    this.badgeCount = 0,
  });

  final _NavItem item;
  final bool expanded;
  final int badgeCount;

  @override
  Widget build(BuildContext context) {
    final currentRoute = Get.currentRoute;
    final selected =
        currentRoute == item.route ||
        (item.route == AppRoutes.settings &&
            currentRoute.startsWith(AppRoutes.settings)) ||
        (item.route == AppRoutes.dashboard && currentRoute.isEmpty);
    return _SideNavAction(
      expanded: expanded,
      icon: item.icon,
      title: item.title,
      selected: selected,
      badgeCount: badgeCount,
      onTap: () {
        if (Get.currentRoute != item.route) {
          _replaceRoute(item.route);
        }
      },
    );
  }
}

class _SideNavAction extends StatelessWidget {
  const _SideNavAction({
    required this.expanded,
    required this.icon,
    required this.title,
    required this.onTap,
    this.selected = false,
    this.badgeCount = 0,
  });

  final bool expanded;
  final IconData icon;
  final String title;
  final bool selected;
  final int badgeCount;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final activeColor = isDark ? AppColors.electricCyan : AppColors.secondary;
    final textColor = selected
        ? activeColor
        : Theme.of(context).textTheme.bodyLarge?.color;
    return Tooltip(
      message: expanded ? '' : title,
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(16),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 220),
            height: 54,
            padding: EdgeInsets.symmetric(horizontal: expanded ? 14 : 0),
            decoration: BoxDecoration(
              color: selected
                  ? activeColor.withValues(alpha: isDark ? 0.16 : 0.12)
                  : Colors.transparent,
              borderRadius: BorderRadius.circular(16),
              border: Border.all(
                color: selected
                    ? activeColor.withValues(alpha: 0.34)
                    : Colors.transparent,
              ),
            ),
            child: Row(
              mainAxisAlignment: expanded
                  ? MainAxisAlignment.start
                  : MainAxisAlignment.center,
              children: [
                Stack(
                  clipBehavior: Clip.none,
                  children: [
                    Icon(icon, color: selected ? activeColor : null),
                    if (badgeCount > 0)
                      Positioned(
                        top: -7,
                        left: -8,
                        child: _PendingBadge(count: badgeCount),
                      ),
                  ],
                ),
                if (expanded) ...[
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          title,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                          style: TextStyle(
                            color: textColor,
                            fontWeight: selected
                                ? FontWeight.w800
                                : FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _PendingBadge extends StatelessWidget {
  const _PendingBadge({required this.count});

  final int count;

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: const BoxConstraints(minWidth: 18, minHeight: 18),
      padding: const EdgeInsets.symmetric(horizontal: 5),
      decoration: BoxDecoration(
        color: AppColors.error,
        borderRadius: BorderRadius.circular(9),
        border: Border.all(color: Colors.white.withValues(alpha: 0.9)),
      ),
      alignment: Alignment.center,
      child: Text(
        count > 9 ? '9+' : '$count',
        style: const TextStyle(
          color: Colors.white,
          fontSize: 10,
          fontWeight: FontWeight.w900,
        ),
      ),
    );
  }
}

class _SideThemeButton extends StatelessWidget {
  const _SideThemeButton({required this.expanded});

  final bool expanded;

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<ThemeController>();
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final accent = isDark ? AppColors.electricCyan : AppColors.secondary;
    return Obx(
      () => Tooltip(
        message: expanded ? '' : 'تم: ${controller.label}',
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: controller.cycle,
            borderRadius: BorderRadius.circular(20),
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 240),
              curve: Curves.easeOutCubic,
              height: 62,
              padding: EdgeInsets.symmetric(horizontal: expanded ? 12 : 0),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topRight,
                  end: Alignment.bottomLeft,
                  colors: isDark
                      ? [
                          AppColors.darkSurfaceAlt.withValues(alpha: 0.92),
                          const Color(0xFF0B1728).withValues(alpha: 0.92),
                        ]
                      : [
                          Colors.white.withValues(alpha: 0.94),
                          const Color(0xFFEAF2FF).withValues(alpha: 0.88),
                        ],
                ),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: accent.withValues(alpha: 0.28)),
                boxShadow: [
                  BoxShadow(
                    color: accent.withValues(alpha: isDark ? 0.14 : 0.18),
                    blurRadius: 18,
                    offset: const Offset(0, 10),
                  ),
                ],
              ),
              child: Row(
                mainAxisAlignment: expanded
                    ? MainAxisAlignment.start
                    : MainAxisAlignment.center,
                children: [
                  Container(
                    width: 38,
                    height: 38,
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [
                          accent.withValues(alpha: 0.24),
                          AppColors.accentGold.withValues(alpha: 0.18),
                        ],
                        begin: Alignment.topRight,
                        end: Alignment.bottomLeft,
                      ),
                      borderRadius: BorderRadius.circular(14),
                    ),
                    child: Icon(controller.icon, color: accent, size: 22),
                  ),
                  if (expanded) ...[
                    const SizedBox(width: 10),
                    Expanded(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'تم برنامه',
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: Theme.of(context).textTheme.labelMedium
                                ?.copyWith(
                                  color: Theme.of(context).hintColor,
                                  fontWeight: FontWeight.w700,
                                ),
                          ),
                          const SizedBox(height: 3),
                          Text(
                            controller.label,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: Theme.of(context).textTheme.bodyMedium
                                ?.copyWith(
                                  color: accent,
                                  fontWeight: FontWeight.w900,
                                ),
                          ),
                        ],
                      ),
                    ),
                    Icon(
                      Icons.sync_rounded,
                      color: accent.withValues(alpha: 0.78),
                      size: 18,
                    ),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _NavItem {
  const _NavItem({
    required this.title,
    required this.icon,
    required this.route,
  });

  final String title;
  final IconData icon;
  final String route;
}

class _DrawerThemeTile extends StatelessWidget {
  const _DrawerThemeTile();

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<ThemeController>();
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final accent = isDark ? AppColors.electricCyan : AppColors.secondary;
    return Obx(
      () => Padding(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 16),
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: controller.cycle,
            borderRadius: BorderRadius.circular(20),
            child: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topRight,
                  end: Alignment.bottomLeft,
                  colors: isDark
                      ? [
                          AppColors.darkSurfaceAlt.withValues(alpha: 0.92),
                          const Color(0xFF0B1728).withValues(alpha: 0.92),
                        ]
                      : [
                          Colors.white.withValues(alpha: 0.96),
                          const Color(0xFFEAF2FF).withValues(alpha: 0.90),
                        ],
                ),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: accent.withValues(alpha: 0.28)),
                boxShadow: [
                  BoxShadow(
                    color: accent.withValues(alpha: isDark ? 0.12 : 0.16),
                    blurRadius: 16,
                    offset: const Offset(0, 8),
                  ),
                ],
              ),
              child: Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: accent.withValues(alpha: 0.14),
                      borderRadius: BorderRadius.circular(14),
                    ),
                    child: Icon(controller.icon, color: accent),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          'تم برنامه',
                          style: Theme.of(context).textTheme.labelMedium
                              ?.copyWith(
                                color: Theme.of(context).hintColor,
                                fontWeight: FontWeight.w700,
                              ),
                        ),
                        const SizedBox(height: 3),
                        Text(
                          controller.label,
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: accent,
                                fontWeight: FontWeight.w900,
                              ),
                        ),
                      ],
                    ),
                  ),
                  Icon(Icons.sync_rounded, color: accent, size: 20),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

void _replaceRoute(String route) {
  if (Get.currentRoute == route) return;
  Get.offNamed(route);
}

class AmlakDrawer extends StatelessWidget {
  const AmlakDrawer({super.key});

  @override
  Widget build(BuildContext context) {
    final consultants = Get.find<ConsultantsController>();
    return Drawer(
      child: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: GradientHeader(
                title: 'Amlak CRM',
                subtitle: 'منوی اصلی مدیریت املاک',
                icon: Icons.real_estate_agent_outlined,
              ),
            ),
            ListTile(
              leading: const Icon(Icons.dashboard_customize_outlined),
              title: const Text('داشبورد'),
              onTap: () => _replaceRoute(AppRoutes.dashboard),
            ),
            ListTile(
              leading: const Icon(Icons.apartment_outlined),
              title: const Text('فایل‌ها'),
              onTap: () => _replaceRoute(AppRoutes.properties),
            ),
            ListTile(
              leading: const Icon(Icons.dynamic_feed_outlined),
              title: const Text('جدیدترین فایل‌ها'),
              onTap: () => _replaceRoute(AppRoutes.latestFiles),
            ),
            ListTile(
              leading: const Icon(Icons.chat_bubble_outline_rounded),
              title: const Text('چت‌ها'),
              onTap: () => _replaceRoute(AppRoutes.chats),
            ),
            ListTile(
              leading: const Icon(Icons.inventory_2_outlined),
              title: const Text('صندوقچه‌ها'),
              onTap: () => _replaceRoute(AppRoutes.vaults),
            ),
            ListTile(
              leading: const Icon(Icons.manage_accounts_outlined),
              title: const Text('مشاورین'),
              onTap: () => _replaceRoute(AppRoutes.consultants),
            ),
            Obx(
              () => consultants.pendingInboxCount > 0
                  ? ListTile(
                      leading: Badge(
                        label: Text('${consultants.pendingInboxCount}'),
                        child: const Icon(Icons.inbox_outlined),
                      ),
                      title: const Text('دعوت‌نامه‌ها'),
                      onTap: () => _replaceRoute(AppRoutes.inbox),
                    )
                  : const SizedBox.shrink(),
            ),
            ListTile(
              leading: const Icon(Icons.settings_outlined),
              title: const Text('تنظیمات'),
              onTap: () => _replaceRoute(AppRoutes.settings),
            ),
            ListTile(
              leading: const Icon(Icons.admin_panel_settings_outlined),
              title: const Text('مدیریت کل'),
              onTap: () => _replaceRoute(AppRoutes.admin),
            ),
            const Spacer(),
            const _DrawerThemeTile(),
          ],
        ),
      ),
    );
  }
}

class GlassAppBar extends StatelessWidget implements PreferredSizeWidget {
  const GlassAppBar({
    super.key,
    required this.title,
    this.actions,
    this.leading,
  });

  final Widget title;
  final List<Widget>? actions;
  final Widget? leading;

  @override
  Size get preferredSize => const Size.fromHeight(82);

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topRight,
          end: Alignment.bottomLeft,
          colors: isDark
              ? const [Color(0xFF050B14), Color(0xFF0A1F36), Color(0xFF102A45)]
              : const [Color(0xFFF7F9FC), Color(0xFFEAF2FF), Color(0xFFFFFAE8)],
        ),
      ),
      child: SafeArea(
        bottom: false,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(12, 8, 12, 0),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(24),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 14, sigmaY: 14),
              child: AppBar(
                title: title,
                leading: leading,
                actions: actions,
                automaticallyImplyLeading: true,
                backgroundColor: isDark
                    ? AppColors.darkSurface.withValues(alpha: 0.86)
                    : AppColors.primary.withValues(alpha: 0.90),
                elevation: 0,
                shape: const RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(Radius.circular(24)),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class StatCard extends StatelessWidget {
  const StatCard({
    super.key,
    required this.title,
    required this.value,
    required this.icon,
  });
  final String title;
  final String value;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0.96, end: 1),
      duration: const Duration(milliseconds: 420),
      curve: Curves.easeOutBack,
      builder: (context, scale, child) {
        return Transform.scale(scale: scale, child: child);
      },
      child: Card(
        child: Container(
          padding: const EdgeInsets.all(18),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(18),
            gradient: LinearGradient(
              colors: isDark
                  ? const [Color(0xFF102238), Color(0xFF0D1B2E)]
                  : const [Colors.white, Color(0xFFF8FBFF)],
              begin: Alignment.topRight,
              end: Alignment.bottomLeft,
            ),
          ),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    colors: [
                      (isDark ? AppColors.electricCyan : AppColors.secondary)
                          .withValues(alpha: isDark ? 0.20 : 0.16),
                      AppColors.accentGold.withValues(alpha: 0.18),
                    ],
                    begin: Alignment.topRight,
                    end: Alignment.bottomLeft,
                  ),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Icon(
                  icon,
                  size: 26,
                  color: isDark ? AppColors.electricCyan : AppColors.primary,
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(title, style: Theme.of(context).textTheme.bodyMedium),
                    const SizedBox(height: 6),
                    Text(
                      value,
                      style: Theme.of(context).textTheme.headlineSmall
                          ?.copyWith(
                            fontWeight: FontWeight.w900,
                            color: isDark
                                ? AppColors.electricCyan
                                : AppColors.primary,
                          ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class PremiumActionButton extends StatelessWidget {
  const PremiumActionButton({
    super.key,
    required this.onPressed,
    required this.child,
    this.icon,
    this.matchTextDirection = false,
  });

  final VoidCallback? onPressed;
  final Widget child;
  final IconData? icon;
  final bool matchTextDirection;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: isDark
              ? const [Color(0xFF4DE1FF), Color(0xFF2563EB), Color(0xFF14213D)]
              : const [AppColors.secondary, AppColors.primary],
          begin: Alignment.centerRight,
          end: Alignment.centerLeft,
        ),
        borderRadius: BorderRadius.circular(14),
        boxShadow: [
          BoxShadow(
            color: (isDark ? AppColors.electricCyan : AppColors.secondary)
                .withValues(alpha: isDark ? 0.18 : 0.24),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: ElevatedButton.icon(
        onPressed: onPressed,
        icon: icon == null
            ? const SizedBox.shrink()
            : Transform.scale(
                scaleX:
                    matchTextDirection &&
                        Directionality.of(context) == TextDirection.rtl
                    ? -1
                    : 1,
                child: Icon(icon),
              ),
        label: child,
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.transparent,
          shadowColor: Colors.transparent,
          disabledBackgroundColor: Colors.transparent,
          minimumSize: const Size(double.infinity, 54),
        ),
      ),
    );
  }
}

class AppCard extends StatelessWidget {
  const AppCard({super.key, required this.child, this.padding = 24});
  final Widget child;
  final double padding;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return Card(
      child: ClipRRect(
        borderRadius: BorderRadius.circular(18),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 12, sigmaY: 12),
          child: DecoratedBox(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: isDark
                    ? [
                        AppColors.darkSurface.withValues(alpha: 0.92),
                        AppColors.darkSurfaceAlt.withValues(alpha: 0.74),
                      ]
                    : [
                        Colors.white.withValues(alpha: 0.94),
                        Colors.white.withValues(alpha: 0.78),
                      ],
                begin: Alignment.topRight,
                end: Alignment.bottomLeft,
              ),
            ),
            child: Padding(padding: EdgeInsets.all(padding), child: child),
          ),
        ),
      ),
    );
  }
}
