import { Stack } from 'expo-router';

export default function BankHooksLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: 'Bank hooks' }} />
    </Stack>
  );
}
