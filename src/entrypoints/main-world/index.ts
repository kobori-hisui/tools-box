/** Se inyecta al mundo MAIN con injectScript('/main-world.js'). */
export default defineUnlistedScript(() => {
  document.currentScript?.dispatchEvent(new CustomEvent('toolbox:main-world-ready'));
});
