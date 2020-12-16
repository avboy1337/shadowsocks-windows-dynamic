# shadowsocks-windows-dynamic
[![](https://img.shields.io/badge/Telegram-Group-blue)](https://t.me/aioCloud)
[![](https://img.shields.io/badge/Telegram-Channel-green)](https://t.me/aioCloud_channel) 

```csharp
namespace Shadowsocks
{
    public enum NameList : int
    {
        TYPE_LISN,
        TYPE_HOST,
        TYPE_PASS,
        TYPE_METH,
        TYPE_OBFS,
        TYPE_OBPA
    }

    public static class NativeMethods
    {
        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern bool ServerInfo(int name, byte[] value);

        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern bool Create();

        [DllImport("shadowsocks.bin", CallingConvention = CallingConvention.Cdecl)]
        public static extern void Delete();
    }

    public class Shadowsocks
    {
        public static void Create()
        {
            NativeMethods.ServerInfo((int)NameList.TYPE_LISN, Encoding.UTF8.GetBytes(":1080"));
            NativeMethods.ServerInfo((int)NameList.TYPE_HOST, Encoding.UTF8.GetBytes("1.1.1.1:80"));
            NativeMethods.ServerInfo((int)NameList.TYPE_PASS, Encoding.UTF8.GetBytes("114514"));
            NativeMethods.ServerInfo((int)NameList.TYPE_METH, Encoding.UTF8.GetBytes("chacha20-ietf"));
            NativeMethods.ServerInfo((int)NameList.TYPE_OBFS, Encoding.UTF8.GetBytes("HTTP"));
            NativeMethods.ServerInfo((int)NameList.TYPE_OBPA, Encoding.UTF8.GetBytes("dash.cloudflare.com"));

            if (!NativeMethods.Create())
            {
                Console.WriteLine("!NativeMethods.Create");
                Console.ReadLine();
                return;
            }
        }

        public static void Main(string[] args)
        {
            Create();
            Console.WriteLine("STARTED");
            Console.ReadLine();

            NativeMethods.Delete();
            Console.WriteLine("STOPPED");
            Console.ReadLine();
        }
    }
}
```
