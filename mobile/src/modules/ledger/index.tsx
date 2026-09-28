import { StyleSheet, Text, View } from 'react-native';

export function LedgerHomeScreen() {
  return (
    <View style={styles.container}>
      <Text>Ledger</Text>
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
