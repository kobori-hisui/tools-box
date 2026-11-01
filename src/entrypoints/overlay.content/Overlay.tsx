import { useEffect, useState } from 'react';

export default function Overlay() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const listener = (msg: unknown) => {
      if ((msg as { type?: string })?.type === 'palette:toggle') setOpen((v) => !v);
    };
    browser.runtime.onMessage.addListener(listener);
    return () => browser.runtime.onMessage.removeListener(listener);
  }, []);

  if (!open) return null;
  return <div className="tb-palette">Command palette</div>;
}
