export function registerContextMenus() {
  browser.runtime.onInstalled.addListener(() => {
    browser.contextMenus.create({ id: 'toolbox-root', title: 'Toolbox', contexts: ['all'] });
  });
  browser.contextMenus.onClicked.addListener((info) => {
    // TODO: delegar a tool.onContextMenu(info)
    console.debug('[ctx-menu]', info.menuItemId);
  });
}
