import 'package:flutter/material.dart';
import 'package:lab02_chat/user_service.dart';

// UserProfile displays and updates user info
class UserProfile extends StatefulWidget {
  final UserService
      userService; // Accepts a user service for fetching user info
  const UserProfile({Key? key, required this.userService}) : super(key: key);

  @override
  State<UserProfile> createState() => _UserProfileState();
}

class _UserProfileState extends State<UserProfile> {
  // TODO: Add state for user data, loading, and error
  // TODO: Fetch user info from userService (simulate for tests)
  Map <String, String>? userData;
  bool isLoading = true;
  String? errorMessage;
  
  @override
  void initState() {
    super.initState();
    // TODO: Fetch user info and update state
    _fetchUserInfo();
  }

   // Fetch user info from userService
  Future<void> _fetchUserInfo() async {
    setState(() {
      isLoading = true;
      errorMessage = null;
    });

    try {
      final data = await widget.userService.fetchUser();
      setState(() {
        userData = data;
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        errorMessage = e.toString();
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    // TODO: Build user profile UI with loading, error, and user info
     return Scaffold(
      appBar: AppBar(title: const Text('User Profile')),
      body: Center(
        child: _buildProfileContent(),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _fetchUserInfo,
        tooltip: 'Refresh',
        child: const Icon(Icons.refresh),
      ),
    );
  }

  Widget _buildProfileContent() {
    if (isLoading) {
      return const CircularProgressIndicator();
    }

    if (errorMessage != null) {
      return Text(
        'Error: $errorMessage',
        style: const TextStyle(color: Colors.red),
      );
    }

    if (userData == null) {
      return const Text('No user data available');
    }

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          userData!['name'] ?? 'No name',
          style: Theme.of(context).textTheme.headline5,
        ),
        const SizedBox(height: 16),
        Text(
          userData!['email'] ?? 'No email',
          style: Theme.of(context).textTheme.subtitle1,
        ),
      ],
    );
  }
}
