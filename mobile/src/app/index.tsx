import { useEffect, useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';

import { getLiveness } from '@/lib/api/health';

type CheckState =
  { kind: 'loading' } | { kind: 'ok'; status: string } | { kind: 'error'; message: string };

export default function Index() {
  const [check, setCheck] = useState<CheckState>({ kind: 'loading' });

  useEffect(() => {
    let cancelled = false;
    getLiveness()
      .then((health) => {
        if (!cancelled) {
          setCheck({ kind: 'ok', status: health.status });
        }
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setCheck({ kind: 'error', message: error instanceof Error ? error.message : 'unknown' });
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <View style={styles.container}>
      <Text>Bogbank</Text>
      <Text testID="health-status">Backend: {describe(check)}</Text>
    </View>
  );
}

function describe(check: CheckState): string {
  switch (check.kind) {
    case 'loading':
      return 'checking...';
    case 'ok':
      return check.status;
    case 'error':
      return `error (${check.message})`;
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
  },
});
