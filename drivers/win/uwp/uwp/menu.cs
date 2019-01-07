using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Windows.Data.Json;
using Windows.System;
using Windows.UI.Core;
using Windows.UI.Xaml;
using Windows.UI.Xaml.Controls;
using Windows.UI.Xaml.Input;

namespace uwp
{
    public class Menu
    {
        public string ID { get; set; }
        Dictionary<string, object> Nodes { get; set; }
        public CompoNode Root { get; set; }


        public Menu(string ID)
        {
            this.ID = ID;
            this.Nodes = new Dictionary<string, object>();
        }

        public static void New(JsonObject input, string returnID)
        {
            var menu = new Menu(input.GetNamedString("ID"));
            Bridge.PutElem(menu.ID, menu);
            Bridge.Return(returnID, null, null);
        }

        public static void Load(JsonObject input, string returnID)
        {
            var menu = Bridge.GetElem<Menu>(input.GetNamedString("ID"));
            menu.Root = null;
            Bridge.Return(returnID, null, null);
        }

        public static async void Render(JsonObject input, string returnID)
        {
            await Window.Current.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                try
                {
                    var menu = Bridge.GetElem<Menu>(input.GetNamedString("ID"));
                    var changes = JsonArray.Parse(input.GetNamedString("Changes"));

                    foreach (var c in changes)
                    {
                        var change = c.GetObject();
                        var action = change.GetNamedNumber("Action");

                        switch (action)
                        {
                            case 0:
                                menu.setRoot(change);
                                break;

                            case 1:
                                menu.newNode(change);
                                break;

                            case 2:
                                menu.delNode(change);
                                break;

                            case 3:
                                menu.setAttr(change);
                                break;

                            case 4:
                                menu.delAttr(change);
                                break;

                            case 6:
                                menu.appendChild(change);
                                break;

                            default:
                                throw new Exception(string.Format("{0} change is not supported", action));
                        }
                    }

                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        public void setRoot(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var c = this.Nodes[nodeID] as CompoNode;
            c.isRootCompo = true;
            this.Root = c;
        }

        public void newNode(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var compoID = change.GetNamedString("CompoID", "");
            var type = change.GetNamedString("Type");
            var isCompo = change.GetNamedBoolean("IsCompo", false);

            if (isCompo)
            {
                var c = new CompoNode()
                {
                    ID = nodeID,
                    type = type,
                    isRootCompo = false,
                };

                this.Nodes[nodeID] = c;
                return;
            }

            if (type == "menu")
            {
                var container = new MenuContainer()
                {
                    ID = nodeID,
                    compoID = compoID,
                    elemID = this.ID,
                    item = new MenuFlyoutSubItem(),
                };


                container.item.FontSize = 12;
                this.Nodes[nodeID] = container;
                return;
            }

            if (type == "menuitem")
            {
                var item = new MenuItem()
                {
                    ID = nodeID,
                    compoID = compoID,
                    elemID = this.ID,
                    item = new MenuFlyoutItem(),
                };

                item.item.FontSize = 12;
                this.Nodes[nodeID] = item;
                return;
            }

            throw new Exception(string.Format("menu does not support {0} tag", type));
        }

        public void delNode(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            this.Nodes.Remove(nodeID);
        }

        public void setAttr(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var key = change.GetNamedString("Key");
            var value = change.GetNamedString("Value", "");

            var node = this.Nodes[nodeID] as IMenuWithAttr;
            node.setAttr(key, value);
        }

        public void delAttr(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var key = change.GetNamedString("Key");

            var node = this.Nodes[nodeID] as IMenuWithAttr;
            node.delAttr(key);
        }

        public void appendChild(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var childID = change.GetNamedString("ChildID");
            var node = this.Nodes[nodeID];

            if (node is CompoNode)
            {
                var cnode = node as CompoNode;
                cnode.rootID = childID;
                return;
            }

            var parent = node as MenuContainer;
            var child = this.Nodes[childID];
            var childRoot = this.CompoRoot(child);
            parent.appendChild(childRoot);
        }

        public object CompoRoot(object node)
        {
            if (node == null || !(node is CompoNode))
            {
                return node;
            }

            var c = node as CompoNode;
            return this.CompoRoot(this.Nodes[c.rootID]);
        }
    }

    public class CompoNode
    {
        public string ID;
        public string rootID;
        public string type;
        public bool isRootCompo;
    }

    public interface IMenuWithAttr
    {
        void setAttr(string key, string value);
        void delAttr(string key);
    }

    public class MenuContainer : IMenuWithAttr
    {
        public string ID { get; set; }
        public string compoID { get; set; }
        public string elemID { get; set; }
        public MenuFlyoutSubItem item { get; set; }

        public void setAttr(string key, string value)
        {
            switch (key)
            {
                case "label":
                    item.Text = value;
                    break;
            }
        }

        public void delAttr(string key)
        {
            switch (key)
            {
                case "label":
                    item.Text = "";
                    break;
            }
        }

        public void appendChild(object child)
        {
            if (child is MenuContainer)
            {
                var container = child as MenuContainer;
                this.item.Items.Add(container.item);
                return;
            }

            if (child is MenuItem)
            {
                var item = child as MenuItem;

                if (item.separator != null)
                {
                    this.item.Items.Add(item.separator);
                    return;
                }

                this.item.Items.Add(item.item);
                return;
            }

            throw new Exception("unknow child node type: " + child.GetType().ToString());
        }
    }

    public class MenuItem : IMenuWithAttr
    {
        public string ID { get; set; }
        public string compoID { get; set; }
        public string elemID { get; set; }
        public MenuFlyoutItem item { get; set; }
        public MenuFlyoutSeparator separator { get; set; }

        public void setAttr(string key, string value)
        {
            switch (key)
            {
                case "label":
                    item.Text = value;
                    return;

                case "role":
                    this.setRole(value);
                    return;

                case "separator":
                    this.separator = new MenuFlyoutSeparator();

                    if (this.item.Parent == null)
                    {
                        return;
                    }

                    var parent = this.item.Parent as MenuFlyoutSubItem;
                    var idx = parent.Items.IndexOf(this.item);
                    parent.Items.Insert(idx, this.separator);
                    parent.Items.Remove(this.item);
                    return;
            }
        }

        public void delAttr(string key)
        {
            switch (key)
            {
                case "label":
                    item.Text = "";
                    return;

                case "separator":
                    if (this.item.Parent == null)
                    {
                        return;
                    }

                    var parent = this.item.Parent as MenuFlyoutSubItem;
                    var idx = parent.Items.IndexOf(this.separator);
                    parent.Items.Insert(idx, this.item);
                    parent.Items.Remove(this.separator);

                    this.separator = null;
                    return;
            }
        }

        public void setRole(string role)
        {
            switch (role)
            {
                case "undo":
                    this.item.Icon = new SymbolIcon(Symbol.Undo);
                    break;

                case "redo":
                    this.item.Icon = new SymbolIcon(Symbol.Redo);
                    break;

                case "cut":
                    this.item.Icon = new SymbolIcon(Symbol.Cut);
                    this.setKeys("cmdorctrl+x");
                    break;

                case "copy":
                    this.item.Icon = new SymbolIcon(Symbol.Copy);
                    this.setKeys("cmdorctrl+c");
                    break;


                case "paste":
                    this.item.Icon = new SymbolIcon(Symbol.Paste);
                    this.setKeys("cmdorctrl+v");
                    break;

                case "pasteAndMatchStyle":
                    this.item.Icon = new SymbolIcon(Symbol.Paste);
                    break;

                case "selectAll":
                    this.item.Icon = new SymbolIcon(Symbol.SelectAll);
                    this.setKeys("cmdorctrl+a");
                    break;

                case "delete":
                    this.item.Icon = new SymbolIcon(Symbol.Delete);
                    break;

                case "minimize":
                    this.item.Visibility = Visibility.Collapsed;
                    break;

                case "close":
                    break;

                case "quit":
                    break;

                case "reload":
                    this.item.Icon = new SymbolIcon(Symbol.Refresh);
                    break;

                case "forceReload":
                    this.item.Icon = new SymbolIcon(Symbol.Refresh);
                    break;

                case "toggleFullScreen":
                    this.item.Icon = new SymbolIcon(Symbol.FullScreen);
                    break;

                default:
                    this.item.Visibility = Visibility.Visible;
                    break;
            }
        }

        public void setKeys(string keys)
        {
            keys = keys.ToLower();


            var acc = new KeyboardAccelerator();

            foreach (var k in keys.Split('+'))
            {
                switch (k)
                {
                    case "ctrl":
                    case "cmdorctrl":
                        acc.Modifiers |= VirtualKeyModifiers.Control;
                        break;

                    case "shift":
                        acc.Modifiers |= VirtualKeyModifiers.Shift;
                        break;

                    case "fn":
                        acc.Modifiers |= VirtualKeyModifiers.Menu;
                        break;

                    case "meta":
                        acc.Modifiers |= VirtualKeyModifiers.Windows;
                        break;

                    case "":
                    default:
                        var key = k;

                        if (key.Length == 0)
                        {
                            key = "+";
                        }

                        acc.Key = this.ParseVirtualKey(key);
                        break;
                }
            }

            this.item.KeyboardAccelerators.Clear();
            this.item.KeyboardAccelerators.Add(acc);
        }


        public VirtualKey ParseVirtualKey(string key)
        {
            switch (key)
            {
                case "a":
                    return VirtualKey.A;

                case "b":
                    return VirtualKey.B;

                case "c":
                    return VirtualKey.C;

                case "d":
                    return VirtualKey.D;

                case "e":
                    return VirtualKey.E;

                case "f":
                    return VirtualKey.F;

                case "g":
                    return VirtualKey.G;

                case "h":
                    return VirtualKey.H;

                case "i":
                    return VirtualKey.I;

                case "j":
                    return VirtualKey.J;

                case "k":
                    return VirtualKey.K;

                case "l":
                    return VirtualKey.L;

                case "m":
                    return VirtualKey.M;

                case "n":
                    return VirtualKey.N;

                case "o":
                    return VirtualKey.O;

                case "p":
                    return VirtualKey.P;

                case "q":
                    return VirtualKey.Q;

                case "r":
                    return VirtualKey.R;

                case "s":
                    return VirtualKey.S;

                case "t":
                    return VirtualKey.T;

                case "u":
                    return VirtualKey.U;

                case "v":
                    return VirtualKey.V;

                case "w":
                    return VirtualKey.W;

                case "x":
                    return VirtualKey.X;

                case "y":
                    return VirtualKey.Y;

                case "z":
                    return VirtualKey.Z;

                case "1":
                    return VirtualKey.Number1;

                case "2":
                    return VirtualKey.Number2;

                case "3":
                    return VirtualKey.Number3;

                case "4":
                    return VirtualKey.Number4;

                case "5":
                    return VirtualKey.Number5;

                case "6":
                    return VirtualKey.Number6;

                case "7":
                    return VirtualKey.Number7;

                case "8":
                    return VirtualKey.Number8;

                case "9":
                    return VirtualKey.Number9;

                case "0":
                    return VirtualKey.Number0;

                case "+":
                    return VirtualKey.Add;

                case "-":
                    return VirtualKey.Subtract;

                case "*":
                    return VirtualKey.Multiply;

                case "/":
                    return VirtualKey.Divide;

                case "f1":
                    return VirtualKey.F1;

                case "f2":
                    return VirtualKey.F2;

                case "f3":
                    return VirtualKey.F3;

                case "f4":
                    return VirtualKey.F4;

                case "f5":
                    return VirtualKey.F5;

                case "f6":
                    return VirtualKey.F6;

                case "f7":
                    return VirtualKey.F7;

                case "f8":
                    return VirtualKey.F8;

                case "f9":
                    return VirtualKey.F9;

                case "f10":
                    return VirtualKey.F10;

                case "f11":
                    return VirtualKey.F11;

                case "f12":
                    return VirtualKey.F12;

                default:
                    return new VirtualKey();
            }
        }
    }
}
