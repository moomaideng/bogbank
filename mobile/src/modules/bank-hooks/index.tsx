import { StyleSheet, Text, View } from 'react-native';

export function BankHooksHomeScreen() {
  return (
    <View style={styles.container}>
      <Text>Bank hooks</Text>
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
