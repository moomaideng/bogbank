import { StyleSheet, Text, View } from 'react-native';

export function ReceiptsHomeScreen() {
  return (
    <View style={styles.container}>
      <Text>Receipts</Text>
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
