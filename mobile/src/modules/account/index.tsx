import { StyleSheet, Text, View } from 'react-native';

export function AccountHomeScreen() {
  return (
    <View style={styles.container}>
      <Text>Account</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
});
