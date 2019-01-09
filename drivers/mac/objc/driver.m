#include "driver.h"
#include "controller.h"
#include "dock.h"
#include "json.h"
#include "menu.h"
#include "notification.h"
#include "panel.h"
#include "sandbox.h"
#include "status.h"
#include "window.h"

@implementation Driver
+ (instancetype)current {
  static Driver *driver = nil;

  @synchronized(self) {
    if (driver == nil) {
      driver = [[Driver alloc] init];
      NSApplication *app = [NSApplication sharedApplication];
      app.delegate = driver;
    }
  }

  return driver;
}

- (instancetype)init {
  self = [super init];

  self.elements = [[NSMutableDictionary alloc] init];
  self.macRPC = [[MacRPC alloc] init];
  self.goRPC = [[GoRPC alloc] init];

  self.roles = @{
    @"undo" : @"undo:",
    @"redo" : @"redo:",
    @"cut" : @"cut:",
    @"copy" : @"copy:",
    @"paste" : @"paste:",
    @"pasteAndMatchStyle" : @"pasteAsPlainText:",
    @"selectAll" : @"selectAll:",
    @"delete" : @"delete:",
    @"minimize" : @"performMiniaturize:",
    @"close" : @"performClose:",
    @"quit" : @"terminate:",
    @"reload" : @"reload:",
    @"forceReload" : @"reloadFromOrigin:",
    @"toggleFullScreen" : @"toggleFullScreen:",

    @"about" : @"orderFrontStandardAboutPanel:",
    @"hide" : @"hide:",
    @"hideOthers" : @"hideOtherApplications:",
    @"unhide" : @"unhideAllApplications:",
    @"startSpeaking" : @"startSpeaking:",
    @"stopSpeaking" : @"stopSpeaking:",
    @"front" : @"arrangeInFront:",
    @"zoom" : @"performZoom:",
  };

  // Driver handlers.
  [self.macRPC handle:@"driver.Run"
          withHandler:^(id in, NSString *returnID) {
            return [self run:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Bundle"
          withHandler:^(id in, NSString *returnID) {
            return [self bundle:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetContextMenu"
          withHandler:^(id in, NSString *returnID) {
            return [self setContextMenu:in return:returnID];
          }];
  [self.macRPC handle:@"driver.SetMenubar"
          withHandler:^(id in, NSString *returnID) {
            return [self setMenubar:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Share"
          withHandler:^(id in, NSString *returnID) {
            return [self share:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Close"
          withHandler:^(id in, NSString *returnID) {
            return [self close:in return:returnID];
          }];
  [self.macRPC handle:@"driver.Terminate"
          withHandler:^(id in, NSString *returnID) {
            return [self terminate:in return:returnID];
          }];

  // Window handlers.
  [self.macRPC handle:@"windows.New"
          withHandler:^(id in, NSString *returnID) {
            return [Window new:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Load"
          withHandler:^(id in, NSString *returnID) {
            return [Window load:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Render"
          withHandler:^(id in, NSString *returnID) {
            return [Window render:in return:returnID];
          }];
  [self.macRPC handle:@"windows.EvalJS"
          withHandler:^(id in, NSString *returnID) {
            return [Window evalJS:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Position"
          withHandler:^(id in, NSString *returnID) {
            return [Window position:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Move"
          withHandler:^(id in, NSString *returnID) {
            return [Window move:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Center"
          withHandler:^(id in, NSString *returnID) {
            return [Window center:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Size"
          withHandler:^(id in, NSString *returnID) {
            return [Window size:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Resize"
          withHandler:^(id in, NSString *returnID) {
            return [Window resize:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Focus"
          withHandler:^(id in, NSString *returnID) {
            return [Window focus:in return:returnID];
          }];
  [self.macRPC handle:@"windows.ToggleFullScreen"
          withHandler:^(id in, NSString *returnID) {
            return [Window toggleFullScreen:in return:returnID];
          }];
  [self.macRPC handle:@"windows.ToggleMinimize"
          withHandler:^(id in, NSString *returnID) {
            return [Window toggleMinimize:in return:returnID];
          }];
  [self.macRPC handle:@"windows.Close"
          withHandler:^(id in, NSString *returnID) {
            return [Window close:in return:returnID];
          }];

  // Menu handlers.
  [self.macRPC handle:@"menus.New"
          withHandler:^(id in, NSString *returnID) {
            return [Menu new:in return:returnID];
          }];
  [self.macRPC handle:@"menus.Load"
          withHandler:^(id in, NSString *returnID) {
            return [Menu load:in return:returnID];
          }];
  [self.macRPC handle:@"menus.Render"
          withHandler:^(id in, NSString *returnID) {
            return [Menu render:in return:returnID];
          }];
  [self.macRPC handle:@"menus.Delete"
          withHandler:^(id in, NSString *returnID) {
            return [Menu delete:in return:returnID];
          }];

  // Status menu handlers.
  [self.macRPC handle:@"statusMenus.New"
          withHandler:^(id in, NSString *returnID) {
            return [StatusMenu new:in return:returnID];
          }];
  [self.macRPC handle:@"statusMenus.SetMenu"
          withHandler:^(id in, NSString *returnID) {
            return [StatusMenu setMenu:in return:returnID];
          }];
  [self.macRPC handle:@"statusMenus.SetText"
          withHandler:^(id in, NSString *returnID) {
            return [StatusMenu setText:in return:returnID];
          }];
  [self.macRPC handle:@"statusMenus.SetIcon"
          withHandler:^(id in, NSString *returnID) {
            return [StatusMenu setIcon:in return:returnID];
          }];
  [self.macRPC handle:@"statusMenus.Close"
          withHandler:^(id in, NSString *returnID) {
            return [StatusMenu close:in return:returnID];
          }];

  // Dock handlers.
  [self.macRPC handle:@"docks.SetMenu"
          withHandler:^(id in, NSString *returnID) {
            return [Dock setMenu:in return:returnID];
          }];
  [self.macRPC handle:@"docks.SetBadge"
          withHandler:^(id in, NSString *returnID) {
            return [Dock setBadge:in return:returnID];
          }];
  [self.macRPC handle:@"docks.SetIcon"
          withHandler:^(id in, NSString *returnID) {
            return [Dock setIcon:in return:returnID];
          }];

  // Controller handlers.
  [self.macRPC handle:@"controller.New"
          withHandler:^(id in, NSString *returnID) {
            return [Controller new:in return:returnID];
          }];
  [self.macRPC handle:@"controller.Listen"
          withHandler:^(id in, NSString *returnID) {
            return [Controller listen:in return:returnID];
          }];
  [self.macRPC handle:@"controller.Close"
          withHandler:^(id in, NSString *returnID) {
            return [Controller close:in return:returnID];
          }];

  // File panel handlers.
  [self.macRPC handle:@"files.NewPanel"
          withHandler:^(id in, NSString *returnID) {
            return [FilePanel newFilePanel:in return:returnID];
          }];
  [self.macRPC handle:@"files.NewSavePanel"
          withHandler:^(id in, NSString *returnID) {
            return [FilePanel newSaveFilePanel:in return:returnID];
          }];

  // Notification handlers.
  [self.macRPC handle:@"notifications.New"
          withHandler:^(id in, NSString *returnID) {
            return [Notification new:in return:returnID];
          }];

  // Notifications.
  NSUserNotificationCenter *userNotificationCenter =
      [NSUserNotificationCenter defaultUserNotificationCenter];
  userNotificationCenter.delegate = self;

  return self;
}

- (void)run:(NSDictionary *)in return:(NSString *)returnID {
  [NSApp run];
  [self.macRPC return:returnID withOutput:nil andError:nil];
}

- (void)bundle:(NSDictionary *)in return:(NSString *)returnID {
  NSBundle *mainBundle = [NSBundle mainBundle];

  NSMutableDictionary *out = [[NSMutableDictionary alloc] init];
  out[@"AppName"] = mainBundle.infoDictionary[@"CFBundleName"];
  out[@"Resources"] = mainBundle.resourcePath;
  out[@"Support"] = [self support];

  [self.macRPC return:returnID withOutput:out andError:nil];
}

- (NSString *)support {
  NSBundle *mainBundle = [NSBundle mainBundle];

  if ([mainBundle isSandboxed]) {
    return NSHomeDirectory();
  }

  NSArray *paths = NSSearchPathForDirectoriesInDomains(
      NSApplicationSupportDirectory, NSUserDomainMask, YES);
  NSString *applicationSupportDirectory = [paths firstObject];

  if (mainBundle.bundleIdentifier.length == 0) {
    return [NSString
        stringWithFormat:@"%@/goapp/{appname}", applicationSupportDirectory];
  }
  return [NSString stringWithFormat:@"%@/%@", applicationSupportDirectory,
                                    mainBundle.bundleIdentifier];
}

- (SEL)selectorFromRole:(NSString *)role {
  NSString *selector = self.roles[role];

  if (selector == nil) {
    return nil;
  }

  return NSSelectorFromString(selector);
}

- (void)setContextMenu:(NSString *)menuID return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[menuID];
    NSWindow *win = NSApp.keyWindow;

    if (win == nil) {
      [self.macRPC return:returnID withOutput:nil
          andError:@"no window to host the context menu"];
      return;
    }

    [menu.root popUpMenuPositioningItem:menu.root.itemArray[0]
                             atLocation:[win mouseLocationOutsideOfEventStream]
                                 inView:win.contentView];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setMenubar:(NSString *)menuID return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[menuID];
    NSApp.mainMenu = menu.root;
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDock:(NSString *)menuID return:(NSString *)returnID {
  defer(returnID, ^{
    Menu *menu = self.elements[menuID];
    self.dock = menu.root;
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDockIcon:(NSString *)icon return:(NSString *)returnID {
  defer(returnID, ^{
    if (icon.length != 0) {
      NSApp.applicationIconImage = [[NSImage alloc] initByReferencingFile:icon];
    } else {
      NSApp.applicationIconImage = nil;
    }

    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)setDockBadge:(NSString *)badge return:(NSString *)returnID {
  defer(returnID, ^{
    [NSApp.dockTile setBadgeLabel:badge];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)share:(NSDictionary *)in return:(NSString *)returnID {
  defer(returnID, ^{
    NSWindow *rawWindow = NSApp.keyWindow;
    if (rawWindow == nil) {
      [NSException raise:@"NoKeyWindowExeption"
                  format:@"no window to host the share menu"];
    }
    Window *win = (Window *)rawWindow.windowController;

    id share = in[@"Share"];
    if ([in[@"Type"] isEqual:@"url"]) {
      [NSURL URLWithString:share];
    }

    NSPoint pos = [win.window mouseLocationOutsideOfEventStream];
    pos = [win.webview convertPoint:pos fromView:rawWindow.contentView];
    NSRect rect = NSMakeRect(pos.x, pos.y, 1, 1);

    NSSharingServicePicker *picker =
        [[NSSharingServicePicker alloc] initWithItems:@[ share ]];
    [picker showRelativeToRect:rect
                        ofView:win.webview
                 preferredEdge:NSMinYEdge];

    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)close:(id)in return:(NSString *)returnID {
  [NSApp terminate:self];
  [self.macRPC return:returnID withOutput:nil andError:nil];
}

- (void)terminate:(id)in return:(NSString *)returnID {
  defer(returnID, ^{
    [NSApp replyToApplicationShouldTerminate:YES];
    [self.macRPC return:returnID withOutput:nil andError:nil];
  });
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  [self.goRPC call:@"driver.OnRun" withInput:nil];
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  [self.goRPC call:@"driver.OnFocus" withInput:nil];
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  [self.goRPC call:@"driver.OnBlur" withInput:nil];
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  NSDictionary *in = @{
    @"HasVisibleWindows" : [NSNumber numberWithBool:flag],
  };

  [self.goRPC call:@"driver.OnReopen" withInput:in];
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  NSDictionary *in = @{
    @"Filenames" : filenames,
  };

  [NSApp activateIgnoringOtherApps:YES];
  [self.goRPC call:@"driver.OnFilesOpen" withInput:in];
}

- (void)applicationWillFinishLaunching:(NSNotification *)aNotification {
  NSAppleEventManager *appleEventManager =
      [NSAppleEventManager sharedAppleEventManager];
  [appleEventManager
      setEventHandler:self
          andSelector:@selector(handleGetURLEvent:withReplyEvent:)
        forEventClass:kInternetEventClass
           andEventID:kAEGetURL];
}

- (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
           withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  NSDictionary *in = @{
    @"URL" : [event paramDescriptorForKeyword:keyDirectObject].stringValue,
  };

  [self.goRPC call:@"driver.OnURLOpen" withInput:in];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  [self.goRPC call:@"driver.OnClose" withInput:nil];
  return NSTerminateLater;
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}

- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center
     shouldPresentNotification:(NSUserNotification *)notification {
  return YES;
}

- (void)userNotificationCenter:(NSUserNotificationCenter *)center
       didActivateNotification:(NSUserNotification *)notification {
  NSMutableDictionary *in = [[NSMutableDictionary alloc] init];
  in[@"ID"] = notification.identifier;

  if (notification.activationType == NSUserNotificationActivationTypeReplied) {
    in[@"Reply"] = notification.response.string;
    [self.goRPC call:@"notifications.OnReply" withInput:in];
  }
}
@end
