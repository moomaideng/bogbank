import { Stack } from 'expo-router';

export default function ReceiptsLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: 'Receipts' }} />
    </Stack>
  );
}
