import { Stack } from 'expo-router';

export default function LedgerLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: 'Ledger' }} />
    </Stack>
  );
}
