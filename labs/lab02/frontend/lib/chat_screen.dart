import 'package:flutter/material.dart';
import 'chat_service.dart';
import 'dart:async';

// ChatScreen displays the chat UI
class ChatScreen extends StatefulWidget {
  final ChatService chatService;
  const ChatScreen({super.key, required this.chatService});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  // TODO: Add TextEditingController for input
  // TODO: Add state for messages, loading, and error
  // TODO: Subscribe to chatService.messageStream
  // TODO: Implement UI for sending and displaying messages
  // TODO: Simulate chat logic for tests (current implementation is a simulation)
  final TextEditingController _textController = TextEditingController();
  final List<String> _messages = [];
  bool _isLoading = true;
  String? _error;
  late StreamSubscription<String> _messageSubscription;

  @override
  void initState() {
    super.initState();
    _connectToChat();
    // TODO: Connect to chat service and set up listeners
  }

  Future<void> _connectToChat() async {
    try {
      await widget.chatService.connect();
      _messageSubscription = widget.chatService.messageStream.listen((message) {
        setState(() {
          _messages.add(message);
        });
      }, onError: (e) {
        setState(() {
          _error = 'Connection error: $e';
        });
      });
      setState(() {
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Connection error: $e';
        _isLoading = false;
      });
    }
  }

  @override
  void dispose() {
    // TODO: Dispose controllers and subscriptions
    _textController.dispose();
    _messageSubscription.cancel();
    super.dispose();
  }

  void _sendMessage() async {
    // TODO: Send message using chatService
    final text = _textController.text.trim();
    if (text.isEmpty) return;

    setState(() {
      _messages.add(text);
      _textController.clear();
    });

    try {
      await widget.chatService.sendMessage(text);
    } catch (e) {
      setState(() {
        _error = 'Failed to send message: $e';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    // TODO: Build chat UI with loading, error, and message list
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(child: Text(_error!));
    }

    return Scaffold(
      appBar: AppBar(title: const Text('Chat')),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              itemCount: _messages.length,
              itemBuilder: (context, index) {
                return ListTile(title: Text(_messages[index]));
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _textController,
                    decoration: const InputDecoration(
                      hintText: 'Type a message',
                    ),
                    onSubmitted: (_) => _sendMessage(),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: _sendMessage,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
