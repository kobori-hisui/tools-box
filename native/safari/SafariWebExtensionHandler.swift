// TEMPLATE - not compiled or tested in this repo (needs macOS + Xcode).
// Copy into the "Safari Web Extension" target that `pnpm safari:convert` creates.
//
// Safari does NOT launch arbitrary executables like Chrome/Firefox native
// messaging. Messages from the extension (browser.runtime.sendNativeMessage)
// arrive here, inside the containing app's extension process. To reuse the Go
// host or Rust sidecar, have the *containing app* start them (or talk to
// toolboxd on 127.0.0.1) and relay from this handler.

import SafariServices
import os.log

class SafariWebExtensionHandler: NSObject, NSExtensionRequestHandling {
    func beginRequest(with context: NSExtensionContext) {
        let item = context.inputItems.first as? NSExtensionItem
        let message = item?.userInfo?[SFExtensionMessageKey]
        os_log(.default, "toolbox: received native message: %@", String(describing: message))

        // TODO: dispatch on message["method"] and call into the app / toolboxd.
        let response = NSExtensionItem()
        response.userInfo = [SFExtensionMessageKey: ["ok": true, "echo": message ?? NSNull()]]
        context.completeRequest(returningItems: [response], completionHandler: nil)
    }
}
