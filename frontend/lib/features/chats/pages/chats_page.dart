import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../core/media/clipboard_files.dart';
import '../../../core/media/media_picker.dart';
import '../../../shared/responsive.dart';
import '../chats_controller.dart';

class ChatsPage extends StatefulWidget {
  const ChatsPage({super.key, this.initialSection = ChatSection.privateChats});

  final ChatSection initialSection;

  @override
  State<ChatsPage> createState() => _ChatsPageState();
}

class _ChatsPageState extends State<ChatsPage> {
  final controller = Get.find<ChatsController>();
  late final selected = widget.initialSection.obs;

  @override
  void initState() {
    super.initState();
    controller.load();
  }

  @override
  Widget build(BuildContext context) {
    return PanelScaffold(
      title: const Text('Ãšâ€ Ã˜ÂªÃ¢â‚¬Å’Ã™â€¡Ã˜Â§'),
      body: ResponsivePage(
        maxWidth: 920,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Obx(
              () => SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(
                  children: controller.visibleSections
                      .map(
                        (section) => Padding(
                          padding: const EdgeInsetsDirectional.only(end: 8),
                          child: ChoiceChip(
                            avatar: Icon(_sectionIcon(section), size: 18),
                            label: Text(_sectionTitle(section)),
                            selected: selected.value == section,
                            onSelected: (_) => selected.value = section,
                          ),
                        ),
                      )
                      .toList(),
                ),
              ),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: AppCard(
                padding: 0,
                child: Obx(() {
                  if (controller.loading.value) {
                    return const Center(child: CircularProgressIndicator());
                  }
                  final sections = controller.visibleSections;
                  final activeSection = sections.contains(selected.value)
                      ? selected.value
                      : sections.first;
                  return RefreshIndicator(
                    onRefresh: controller.load,
                    child: _SectionList(
                      section: activeSection,
                      controller: controller,
                    ),
                  );
                }),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SectionList extends StatelessWidget {
  const _SectionList({required this.section, required this.controller});

  final ChatSection section;
  final ChatsController controller;

  @override
  Widget build(BuildContext context) {
    final items = switch (section) {
      ChatSection.privateChats =>
        controller.privateChats
            .map(
              (item) => _ThreadItem(
                id: item.id,
                title: item.title.isEmpty
                    ? 'Ãšâ€ Ã˜Âª Ã˜Â®Ã˜ÂµÃ™Ë†Ã˜ÂµÃ›Å’'
                    : item.title,
                subtitle: 'Ã˜Â®Ã˜ÂµÃ™Ë†Ã˜ÂµÃ›Å’',
                icon: Icons.person_outline_rounded,
              ),
            )
            .toList(),
      ChatSection.fileChats =>
        controller.fileChats
            .map(
              (item) => _ThreadItem(
                id: item.id,
                title: item.title.isEmpty
                    ? 'Ãšâ€ Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž'
                    : item.title,
                subtitle: item.type == 'business_vault'
                    ? 'Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©'
                    : 'Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â´Ã˜Â®Ã˜ÂµÃ›Å’',
                icon: Icons.forum_outlined,
              ),
            )
            .toList(),
      ChatSection.businessFileChats =>
        controller.businessFileChats
            .map(
              (item) => _ThreadItem(
                id: item.id,
                title: item.title.isEmpty
                    ? 'Ãšâ€ Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©'
                    : item.title,
                subtitle: 'Ã™ÂÃ˜Â§Ã›Å’Ã™â€ž Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©',
                icon: Icons.real_estate_agent_outlined,
              ),
            )
            .toList(),
      ChatSection.personalVaults =>
        controller.personalVaults
            .map(
              (item) => _ThreadItem(
                id: item.channelId,
                title: item.title.isEmpty
                    ? 'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡ Ã˜Â´Ã˜Â®Ã˜ÂµÃ›Å’'
                    : item.title,
                subtitle: item.isMain
                    ? 'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡ Ã˜Â§Ã˜ÂµÃ™â€žÃ›Å’'
                    : 'Ã˜Â´Ã˜Â®Ã˜ÂµÃ›Å’',
                icon: Icons.inventory_2_outlined,
                opensVault: true,
              ),
            )
            .toList(),
      ChatSection.joinedVaults =>
        controller.joinedVaultChannels
            .map(
              (item) => _ThreadItem(
                id: item.id,
                title: item.title.isEmpty
                    ? 'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡ Ã˜Â¹Ã˜Â¶Ã™Ë† Ã˜Â´Ã˜Â¯Ã™â€¡'
                    : item.title,
                subtitle: item.isBusinessVault
                    ? 'Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©'
                    : 'Ã˜Â¹Ã˜Â¶Ã™Ë† Ã˜Â´Ã˜Â¯Ã™â€¡',
                icon: Icons.groups_2_outlined,
                opensVault: true,
              ),
            )
            .toList(),
      ChatSection.businessVaults =>
        controller.businessVaults
            .map(
              (item) => _ThreadItem(
                id: item.channelId,
                title: item.title.isEmpty
                    ? 'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡ Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©'
                    : item.title,
                subtitle: item.isMain
                    ? 'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡ Ã˜Â§Ã˜ÂµÃ™â€žÃ›Å’ Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©'
                    : 'Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©',
                icon: Icons.business_center_outlined,
                opensVault: true,
              ),
            )
            .toList(),
    };

    if (items.isEmpty) {
      return ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        children: [
          SizedBox(height: MediaQuery.sizeOf(context).height * 0.18),
          Icon(
            _sectionIcon(section),
            size: 54,
            color: Theme.of(context).hintColor,
          ),
          const SizedBox(height: 12),
          Center(
            child: Text(
              'Ã™â€¦Ã™Ë†Ã˜Â±Ã˜Â¯Ã›Å’ Ã˜Â¨Ã˜Â±Ã˜Â§Ã›Å’ Ã™â€ Ã™â€¦Ã˜Â§Ã›Å’Ã˜Â´ Ã™Ë†Ã˜Â¬Ã™Ë†Ã˜Â¯ Ã™â€ Ã˜Â¯Ã˜Â§Ã˜Â±Ã˜Â¯',
              style: Theme.of(context).textTheme.titleSmall,
            ),
          ),
        ],
      );
    }

    return ListView.separated(
      physics: const AlwaysScrollableScrollPhysics(),
      itemCount: items.length,
      separatorBuilder: (_, _) => const Divider(height: 1),
      itemBuilder: (context, index) => items[index],
    );
  }
}

class _ThreadItem extends StatelessWidget {
  const _ThreadItem({
    required this.id,
    required this.title,
    required this.subtitle,
    required this.icon,
    this.opensVault = false,
  });

  final String id;
  final String title;
  final String subtitle;
  final IconData icon;
  final bool opensVault;

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).colorScheme.secondary;
    return ListTile(
      minVerticalPadding: 14,
      leading: CircleAvatar(
        backgroundColor: color.withValues(alpha: 0.14),
        child: Icon(icon, color: color),
      ),
      title: Text(
        title,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: const TextStyle(fontWeight: FontWeight.w800),
      ),
      subtitle: Text(subtitle, maxLines: 1, overflow: TextOverflow.ellipsis),
      trailing: const Icon(Icons.chevron_left_rounded),
      onTap: () => showModalBottomSheet<void>(
        context: context,
        isScrollControlled: true,
        useSafeArea: true,
        builder: (_) => opensVault
            ? _VaultFilesSheet(channelId: id, title: title)
            : _ChatThreadSheet(channelId: id, title: title),
      ),
    );
  }
}

class _ChatThreadSheet extends StatefulWidget {
  const _ChatThreadSheet({required this.channelId, required this.title});

  final String channelId;
  final String title;

  @override
  State<_ChatThreadSheet> createState() => _ChatThreadSheetState();
}

class _ChatThreadSheetState extends State<_ChatThreadSheet> {
  final controller = Get.find<ChatsController>();
  final textController = TextEditingController();
  final messagesScrollController = ScrollController();
  late final ClipboardFilePasteDisposer disposePasteListener;

  @override
  void initState() {
    super.initState();
    messagesScrollController.addListener(_handleMessageScroll);
    controller.loadThread(widget.channelId);
    disposePasteListener = listenForClipboardFiles(
      (files) => controller.pasteFiles(widget.channelId, files),
    );
  }

  @override
  void dispose() {
    disposePasteListener();
    messagesScrollController.removeListener(_handleMessageScroll);
    messagesScrollController.dispose();
    textController.dispose();
    super.dispose();
  }

  void _handleMessageScroll() {
    if (!messagesScrollController.hasClients) {
      return;
    }
    final position = messagesScrollController.position;
    if (position.pixels >= position.maxScrollExtent - 120) {
      controller.loadOlderMessages(widget.channelId);
    }
    if (position.pixels <= 80) {
      controller.loadNewerMessages(widget.channelId);
    }
  }

  Future<void> _pickAndSend() async {
    final files = await pickMediaFiles();
    await controller.pasteFiles(widget.channelId, files);
  }

  Future<void> _sendText() async {
    final text = textController.text.trim();
    if (text.isEmpty) {
      return;
    }
    textController.clear();
    await controller.sendMessage(widget.channelId, text: text);
  }

  Future<void> _editMessage(ChannelMessageModel message) async {
    final textEditController = TextEditingController(text: message.text);
    final captionEditController = TextEditingController(text: message.caption);
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('ویرایش پیام'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: textEditController,
              minLines: 1,
              maxLines: 4,
              decoration: const InputDecoration(
                labelText: 'متن پیام',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 10),
            TextField(
              controller: captionEditController,
              minLines: 1,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: 'کپشن',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('انصراف'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('ذخیره'),
          ),
        ],
      ),
    );
    if (confirmed != true) {
      textEditController.dispose();
      captionEditController.dispose();
      return;
    }
    try {
      await controller.editMessage(
        widget.channelId,
        message.id,
        text: textEditController.text,
        caption: captionEditController.text,
      );
    } catch (e) {
      Get.snackbar('خطا', e.toString());
    } finally {
      textEditController.dispose();
      captionEditController.dispose();
    }
  }

  Future<void> _deleteMessage(ChannelMessageModel message) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('حذف پیام'),
        content: const Text('این پیام حذف شود؟'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('انصراف'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('حذف'),
          ),
        ],
      ),
    );
    if (confirmed != true) {
      return;
    }
    try {
      await controller.deleteMessage(widget.channelId, message.id);
    } catch (e) {
      Get.snackbar('خطا', e.toString());
    }
  }

  Future<void> _openPrivateChat(ChannelMemberModel member) async {
    try {
      final channel = await controller.startPrivateChat(member);
      if (!mounted) return;
      showModalBottomSheet<void>(
        context: context,
        isScrollControlled: true,
        useSafeArea: true,
        builder: (_) => _ChatThreadSheet(
          channelId: channel.id,
          title: channel.title.isEmpty
              ? controller.memberTitle(member)
              : channel.title,
        ),
      );
    } catch (e) {
      Get.snackbar('خطا', e.toString());
    }
  }

  void _showVaultPicker(ChannelMemberModel member) {
    final vaults = controller.accessibleVaultTargets;
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxHeight: MediaQuery.sizeOf(context).height * 0.62,
          ),
          child: vaults.isEmpty
              ? const Padding(
                  padding: EdgeInsets.all(24),
                  child: Center(
                    child: Text(
                      'صندوقچه قابل دسترس برای افزودن مخاطب وجود ندارد',
                    ),
                  ),
                )
              : ListView.separated(
                  shrinkWrap: true,
                  itemCount: vaults.length,
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, index) {
                    final vault = vaults[index];
                    return ListTile(
                      leading: const Icon(Icons.inventory_2_outlined),
                      title: Text(
                        vault.title,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      subtitle: Text(
                        vault.subtitle,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      onTap: () async {
                        Navigator.of(sheetContext).pop();
                        try {
                          await controller.addMemberToVault(vault, member);
                          Get.snackbar(
                            'افزوده شد',
                            '${controller.memberTitle(member)} به ${vault.title} اضافه شد.',
                          );
                        } catch (e) {
                          Get.snackbar('خطا', e.toString());
                        }
                      },
                    );
                  },
                ),
        ),
      ),
    );
  }

  Future<void> _showContactCategoryPicker(ChannelMemberModel member) async {
    try {
      await controller.loadContactCategories();
    } catch (e) {
      Get.snackbar('خطا', e.toString());
      return;
    }
    if (!mounted) return;
    final selected = <String>{};
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: StatefulBuilder(
          builder: (context, setSheetState) {
            final tags = controller.contactCategories;
            return ConstrainedBox(
              constraints: BoxConstraints(
                maxHeight: MediaQuery.sizeOf(context).height * 0.68,
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  ListTile(
                    title: const Text(
                      'افزودن به دسته‌بندی مخاطبین',
                      style: TextStyle(fontWeight: FontWeight.w800),
                    ),
                    subtitle: Text(
                      member.phone.isEmpty
                          ? controller.memberTitle(member)
                          : '${controller.memberTitle(member)} | ${member.phone}',
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  const Divider(height: 1),
                  if (tags.isEmpty)
                    const Padding(
                      padding: EdgeInsets.all(24),
                      child: Center(child: Text('دسته‌بندی مخاطبی وجود ندارد')),
                    )
                  else
                    Flexible(
                      child: ListView.builder(
                        shrinkWrap: true,
                        itemCount: tags.length,
                        itemBuilder: (context, index) {
                          final tag = tags[index];
                          return CheckboxListTile(
                            value: selected.contains(tag),
                            secondary: const Icon(Icons.sell_outlined),
                            title: Text(tag),
                            controlAffinity: ListTileControlAffinity.leading,
                            onChanged: (checked) {
                              setSheetState(() {
                                if (checked == true) {
                                  selected.add(tag);
                                } else {
                                  selected.remove(tag);
                                }
                              });
                            },
                          );
                        },
                      ),
                    ),
                  Padding(
                    padding: const EdgeInsets.all(12),
                    child: Row(
                      children: [
                        Expanded(
                          child: OutlinedButton(
                            onPressed: () => Navigator.of(sheetContext).pop(),
                            child: const Text('انصراف'),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: FilledButton.icon(
                            icon: const Icon(Icons.person_add_alt_1_outlined),
                            label: const Text('ثبت'),
                            onPressed: () async {
                              Navigator.of(sheetContext).pop();
                              try {
                                final result = await controller
                                    .addProfileToContactCategories(
                                      member,
                                      selected.toList(),
                                    );
                                final suffix = result.autoConsultant
                                    ? ' این فرد به صورت خودکار در دسته مشاور املاک هم ثبت شد.'
                                    : '';
                                Get.snackbar(
                                  result.existing ? 'به‌روزرسانی شد' : 'ثبت شد',
                                  '${controller.memberTitle(member)} در دسته‌بندی مخاطبین ثبت شد.$suffix',
                                );
                              } catch (e) {
                                Get.snackbar('خطا', e.toString());
                              }
                            },
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  void _showMemberProfile(ChannelMemberModel member) {
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: CircleAvatar(
                child: Text(
                  controller.memberTitle(member).isEmpty
                      ? '#'
                      : controller.memberTitle(member).substring(0, 1),
                ),
              ),
              title: Text(
                controller.memberTitle(member),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(fontWeight: FontWeight.w800),
              ),
              subtitle: Text(
                member.phone.isEmpty ? 'شماره ثبت نشده' : member.phone,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            ListTile(
              leading: const Icon(Icons.chat_bubble_outline_rounded),
              title: const Text('چت با فرد'),
              onTap: () {
                Navigator.of(sheetContext).pop();
                _openPrivateChat(member);
              },
            ),
            ListTile(
              leading: const Icon(Icons.person_add_alt_1_outlined),
              title: const Text('افزودن به صندوقچه'),
              onTap: () {
                Navigator.of(sheetContext).pop();
                _showVaultPicker(member);
              },
            ),
            ListTile(
              leading: const Icon(Icons.contacts_outlined),
              title: const Text('افزودن به دسته‌بندی مخاطبین'),
              onTap: () {
                Navigator.of(sheetContext).pop();
                _showContactCategoryPicker(member);
              },
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: MediaQuery.sizeOf(context).height * 0.86,
      child: Column(
        children: [
          ListTile(
            title: Text(
              widget.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w800),
            ),
            subtitle: Obx(() => Text(_presenceText(controller))),
            trailing: IconButton(
              tooltip: 'Ø¨Ø³ØªÙ†',
              icon: const Icon(Icons.close_rounded),
              onPressed: Get.back,
            ),
          ),
          const Divider(height: 1),
          Expanded(
            child: Obx(() {
              if (controller.threadLoading.value) {
                return const Center(child: CircularProgressIndicator());
              }
              final messages = controller.threadMessages;
              if (messages.isEmpty) {
                return const Center(
                  child: Text('Ù¾ÛŒØ§Ù…ÛŒ ÙˆØ¬ÙˆØ¯ Ù†Ø¯Ø§Ø±Ø¯'),
                );
              }
              return ListView.builder(
                controller: messagesScrollController,
                reverse: true,
                padding: const EdgeInsets.all(12),
                itemCount: messages.length,
                itemBuilder: (context, index) {
                  final message = messages[index];
                  final currentUserId = controller.currentUserId;
                  final isMine = message.isMine(currentUserId);
                  final isSeen = message.seenByOther(currentUserId);
                  final channel = controller.channelById(widget.channelId);
                  final isPrivateChat = channel?.type == 'private';
                  final canModify = isMine && (!isPrivateChat || !isSeen);
                  final member = controller.memberForUserId(message.authorId);
                  final authorName = member == null
                      ? message.authorName
                      : controller.memberTitle(member);
                  final showUnreadDivider =
                      controller.firstUnreadMessageId.value == message.id &&
                      controller.threadUnreadCount.value > 0;
                  return Column(
                    children: [
                      if (showUnreadDivider)
                        _UnreadMessagesDivider(
                          count: controller.threadUnreadCount.value,
                        ),
                      Align(
                        alignment: isMine
                            ? AlignmentDirectional.centerEnd
                            : AlignmentDirectional.centerStart,
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 8),
                          padding: const EdgeInsets.all(10),
                          constraints: const BoxConstraints(maxWidth: 520),
                          decoration: BoxDecoration(
                            color: Theme.of(
                              context,
                            ).colorScheme.surfaceContainerHighest,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              if (canModify)
                                Align(
                                  alignment: AlignmentDirectional.centerEnd,
                                  child: PopupMenuButton<String>(
                                    tooltip: 'گزینه‌های پیام',
                                    icon: const Icon(Icons.more_horiz_rounded),
                                    onSelected: (value) {
                                      if (value == 'edit') {
                                        _editMessage(message);
                                      } else if (value == 'delete') {
                                        _deleteMessage(message);
                                      }
                                    },
                                    itemBuilder: (context) => const [
                                      PopupMenuItem(
                                        value: 'edit',
                                        child: ListTile(
                                          leading: Icon(Icons.edit_outlined),
                                          title: Text('ویرایش'),
                                        ),
                                      ),
                                      PopupMenuItem(
                                        value: 'delete',
                                        child: ListTile(
                                          leading: Icon(Icons.delete_outline),
                                          title: Text('حذف'),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              if (!isMine && authorName.isNotEmpty)
                                InkWell(
                                  borderRadius: BorderRadius.circular(6),
                                  onTap: member == null
                                      ? null
                                      : () => _showMemberProfile(member),
                                  child: Padding(
                                    padding: const EdgeInsetsDirectional.only(
                                      bottom: 2,
                                      end: 6,
                                    ),
                                    child: Text(
                                      authorName,
                                      style: TextStyle(
                                        fontWeight: FontWeight.w800,
                                        color: member == null
                                            ? null
                                            : Theme.of(
                                                context,
                                              ).colorScheme.secondary,
                                      ),
                                    ),
                                  ),
                                ),
                              if (message.text.isNotEmpty) Text(message.text),
                              if (message.caption.isNotEmpty)
                                Text(message.caption),
                              if (message.vaultFileRef != null)
                                Container(
                                  margin: const EdgeInsets.only(top: 6),
                                  padding: const EdgeInsets.all(8),
                                  decoration: BoxDecoration(
                                    border: Border.all(
                                      color: Theme.of(context).dividerColor,
                                    ),
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Row(
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      Icon(
                                        _vaultFileIcon(
                                          message.vaultFileRef!.kind,
                                        ),
                                        size: 18,
                                      ),
                                      const SizedBox(width: 6),
                                      Flexible(
                                        child: Text(
                                          message.vaultFileRef!.title.isEmpty
                                              ? 'فایل صندوقچه'
                                              : message.vaultFileRef!.title,
                                          maxLines: 1,
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              for (final media in message.media)
                                Padding(
                                  padding: const EdgeInsets.only(top: 6),
                                  child: Text(
                                    '${media.kind} - ${(media.size / 1024).ceil()} KB',
                                    style: Theme.of(
                                      context,
                                    ).textTheme.bodySmall,
                                  ),
                                ),
                              if (isMine)
                                Align(
                                  alignment: AlignmentDirectional.centerEnd,
                                  child: Icon(
                                    isSeen
                                        ? Icons.done_all_rounded
                                        : Icons.done_rounded,
                                    size: 18,
                                    color: isSeen
                                        ? Theme.of(
                                            context,
                                          ).colorScheme.secondary
                                        : Theme.of(context).hintColor,
                                  ),
                                ),
                            ],
                          ),
                        ),
                      ),
                    ],
                  );
                },
              );
            }),
          ),
          const Divider(height: 1),
          Padding(
            padding: EdgeInsets.fromLTRB(
              12,
              10,
              12,
              10 + MediaQuery.viewInsetsOf(context).bottom,
            ),
            child: Row(
              children: [
                IconButton(
                  tooltip: 'Ø§Ø±Ø³Ø§Ù„ ÙØ§ÛŒÙ„',
                  onPressed: _pickAndSend,
                  icon: const Icon(Icons.attach_file_rounded),
                ),
                Expanded(
                  child: TextField(
                    controller: textController,
                    minLines: 1,
                    maxLines: 4,
                    decoration: const InputDecoration(
                      hintText: 'Ù¾ÛŒØ§Ù… ÛŒØ§ Ú©Ù¾Ø´Ù† Ø±Ø§ Ø¨Ù†ÙˆÛŒØ³ÛŒØ¯',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                Obx(
                  () => IconButton.filled(
                    tooltip: 'Ø§Ø±Ø³Ø§Ù„',
                    onPressed: controller.sending.value ? null : _sendText,
                    icon: controller.sending.value
                        ? const SizedBox(
                            width: 18,
                            height: 18,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : const Icon(Icons.send_rounded),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _UnreadMessagesDivider extends StatelessWidget {
  const _UnreadMessagesDivider({required this.count});

  final int count;

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).colorScheme.secondary;
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Row(
        children: [
          Expanded(child: Divider(color: color.withValues(alpha: 0.45))),
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 8),
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.12),
              borderRadius: BorderRadius.circular(999),
            ),
            child: Text(
              '$count پیام خوانده‌نشده',
              style: TextStyle(
                color: color,
                fontWeight: FontWeight.w800,
                fontSize: 12,
              ),
            ),
          ),
          Expanded(child: Divider(color: color.withValues(alpha: 0.45))),
        ],
      ),
    );
  }
}

class _VaultFilesSheet extends StatefulWidget {
  const _VaultFilesSheet({required this.channelId, required this.title});

  final String channelId;
  final String title;

  @override
  State<_VaultFilesSheet> createState() => _VaultFilesSheetState();
}

class _VaultFilesSheetState extends State<_VaultFilesSheet> {
  final controller = Get.find<ChatsController>();

  @override
  void initState() {
    super.initState();
    controller.loadVaultFiles(widget.channelId);
  }

  Future<void> _chatAboutFile(ChannelVaultFileModel file) async {
    await controller.startFileChat(widget.channelId, file);
    if (!mounted) return;
    Navigator.of(context).pop();
    showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      builder: (_) =>
          _ChatThreadSheet(channelId: widget.channelId, title: widget.title),
    );
  }

  void _showFileActions(ChannelVaultFileModel file) {
    showModalBottomSheet<void>(
      context: context,
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: Text(
                file.title.isEmpty ? 'فایل صندوقچه' : file.title,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
              subtitle: Text(
                '${file.kind} - ${(file.size / 1024).ceil()} KB',
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            ListTile(
              leading: const Icon(Icons.forum_outlined),
              title: const Text('چت درباره فایل'),
              onTap: () {
                Navigator.of(context).pop();
                _chatAboutFile(file);
              },
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: MediaQuery.sizeOf(context).height * 0.86,
      child: Column(
        children: [
          ListTile(
            title: Text(
              widget.title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w800),
            ),
            subtitle: const Text('برای هر فایل می‌توانید چت مرتبط بسازید.'),
            trailing: IconButton(
              tooltip: 'بستن',
              icon: const Icon(Icons.close_rounded),
              onPressed: Get.back,
            ),
          ),
          const Divider(height: 1),
          Expanded(
            child: Obx(() {
              final files = controller.vaultFiles;
              if (files.isEmpty) {
                return const Center(child: Text('فایلی در صندوقچه وجود ندارد'));
              }
              return ListView.separated(
                itemCount: files.length,
                separatorBuilder: (_, _) => const Divider(height: 1),
                itemBuilder: (context, index) {
                  final file = files[index];
                  return ListTile(
                    leading: CircleAvatar(
                      child: Icon(_vaultFileIcon(file.kind)),
                    ),
                    title: Text(
                      file.title.isEmpty ? 'فایل صندوقچه' : file.title,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    subtitle: Text(
                      '${file.kind} - ${(file.size / 1024).ceil()} KB',
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    trailing: const Icon(Icons.more_horiz_rounded),
                    onTap: () => _showFileActions(file),
                  );
                },
              );
            }),
          ),
        ],
      ),
    );
  }
}

String _presenceText(ChatsController controller) {
  final currentUserId = controller.currentUserId;
  final others = controller.threadMembers
      .where(
        (member) => member.userId.isEmpty || member.userId != currentUserId,
      )
      .toList();
  if (others.length == 1) {
    final member = others.first;
    if (member.isOnline) {
      return 'آنلاین';
    }
    return 'آفلاین';
  }
  final onlineCount = others.where((member) => member.isOnline).length;
  return '$onlineCount آنلاین از ${others.length} عضو';
}

IconData _vaultFileIcon(String kind) {
  switch (kind) {
    case 'image':
      return Icons.image_outlined;
    case 'video':
      return Icons.play_circle_outline_rounded;
    case 'audio':
      return Icons.graphic_eq_rounded;
    default:
      return Icons.description_outlined;
  }
}

String _sectionTitle(ChatSection section) => switch (section) {
  ChatSection.businessFileChats =>
    'Ãšâ€ Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©',
  ChatSection.businessVaults =>
    'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡Ã¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â§Ã™â€¦Ã™â€žÃ˜Â§ÃšÂ©',
  ChatSection.privateChats => 'Ãšâ€ Ã˜Âª Ã˜Â®Ã˜ÂµÃ™Ë†Ã˜ÂµÃ›Å’',
  ChatSection.fileChats => 'Ãšâ€ Ã˜Âª Ã™ÂÃ˜Â§Ã›Å’Ã™â€žÃ¢â‚¬Å’Ã™â€¡Ã˜Â§',
  ChatSection.personalVaults =>
    'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡Ã¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â´Ã˜Â®Ã˜ÂµÃ›Å’',
  ChatSection.joinedVaults =>
    'Ã˜ÂµÃ™â€ Ã˜Â¯Ã™Ë†Ã™â€šÃšâ€ Ã™â€¡Ã¢â‚¬Å’Ã™â€¡Ã˜Â§Ã›Å’ Ã˜Â¹Ã˜Â¶Ã™Ë† Ã˜Â´Ã˜Â¯Ã™â€¡',
};

IconData _sectionIcon(ChatSection section) => switch (section) {
  ChatSection.businessFileChats => Icons.real_estate_agent_outlined,
  ChatSection.businessVaults => Icons.business_center_outlined,
  ChatSection.privateChats => Icons.person_outline_rounded,
  ChatSection.fileChats => Icons.forum_outlined,
  ChatSection.personalVaults => Icons.inventory_2_outlined,
  ChatSection.joinedVaults => Icons.groups_2_outlined,
};
