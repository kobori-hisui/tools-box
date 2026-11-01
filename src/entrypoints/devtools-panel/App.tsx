import { ToolShell } from '@/components/tool/ToolShell';

export default function App({ surface }: { surface: string }) {
  return <ToolShell surface="devtools-panel" />;
}
