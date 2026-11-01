export function registerCommands() {
  browser.commands.onCommand.addListener(async (command) => {
    if (command !== 'open-command-palette') return;
    const [tab] = await browser.tabs.query({ active: true, currentWindow: true });
    if (tab?.id != null) await browser.tabs.sendMessage(tab.id, { type: 'palette:toggle' });
  });
}
