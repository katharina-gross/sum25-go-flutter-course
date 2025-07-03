class UserService {
  Future<Map<String, String>> fetchUser() async {
    // TODO: Simulate fetching user data for tests
    await Future.delayed(const Duration(milliseconds: 10));

    // Return simulated user data
    return {
      'name': 'Default User',
      'email': 'user@example.com'
      };
  }
}
