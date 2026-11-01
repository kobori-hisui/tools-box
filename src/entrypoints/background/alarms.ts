export function registerAlarms() {
  browser.alarms.onAlarm.addListener((alarm) => {
    // TODO: delegar a services/scheduler segun alarm.name
    console.debug('[alarm]', alarm.name);
  });
}
