import { expect, test } from './fixtures';

test('el popup abre y muestra el buscador', async ({ page, extensionId }) => {
  await page.goto(`chrome-extension://${extensionId}/popup.html`);
  await expect(page.getByPlaceholder('Search tools...')).toBeVisible();
});
