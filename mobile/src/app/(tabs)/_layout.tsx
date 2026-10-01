import { Tabs } from 'expo-router';

export default function TabsLayout() {
  return (
    <Tabs screenOptions={{ headerShown: false }}>
      <Tabs.Screen name="dashboard" options={{ title: 'Dashboard' }} />
      <Tabs.Screen name="ledger" options={{ title: 'Ledger' }} />
      <Tabs.Screen name="receipts" options={{ title: 'Receipts' }} />
      <Tabs.Screen name="bank-hooks" options={{ title: 'Bank hooks' }} />
      <Tabs.Screen name="account" options={{ title: 'Account' }} />
    </Tabs>
  );
}
